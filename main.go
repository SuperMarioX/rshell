package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/luckywinds/lwssh"
	. "github.com/luckywinds/rshell/types"
	"github.com/luckywinds/rshell/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

const (
	Interactive string = "interactive"
	SCRIPT      string = "script"
	DOWNLOAD    string = "download"
	UPLOAD      string = "upload"
)

var (
	cfgFile   = flag.String("cfg", "cfg.yaml", "System Config file to read")
	hostsFile = flag.String("hosts", "hosts.yaml", "Hosts Config file to read")
	authFile  = flag.String("auth", "auth.yaml", "Auth Config file to read")
	mode      = flag.String("mode", "interactive", "The Run mode, options[interactive/script], DEFAULT is interactive.")
	script    = flag.String("f", "", "The script yaml.")
	hg        = flag.String("g", "", "The target hostgroup.")
	cmd       = flag.String("c", "", "The commands split by ;")
)

func initFlag() {
	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`Usage:`, cmd, `[<options>] <arguments>

Options:`)
		flag.PrintDefaults()
	}
	flag.Parse()
}

var cfg Cfg

func initCfg() {
	c, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatalf("Can not find cfg file[%s].", *cfgFile)
	}
	err = yaml.Unmarshal(c, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

var hostgroups Hostgroups

func initHosts() {
	h, err := ioutil.ReadFile(*hostsFile)
	if err != nil {
		log.Fatalf("Can not find hosts file[%s].", *hostsFile)
	}
	err = yaml.Unmarshal(h, &hostgroups)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, hg := range hostgroups.Hgs {
		for _, h := range hg.Hosts {
			if net.ParseIP(h) == nil {
				log.Fatalf("IP illegal [%s/%s].", hg.Groupname, h)
			}
		}
	}
}

var auths Auths

func initAuths() {
	a, err := ioutil.ReadFile(*authFile)
	if err != nil {
		log.Fatalf("Can not find auth file[%s].", *authFile)
	}
	err = yaml.Unmarshal(a, &auths)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

var tasks Tasks

func initArgs(playbook string) {
	p, err := ioutil.ReadFile(playbook)
	if err != nil {
		log.Fatalf("Can not find script file[%s].", playbook)
	}

	err = yaml.Unmarshal(p, &tasks)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func hgCompleter(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(ps, d.GetWordBeforeCursor(), true)
}

var ps = []prompt.Suggest{
	{
		Text: DOWNLOAD,
		Description: "COMMAND: Download file from remote host",
	},
	{
		Text: UPLOAD,
		Description: "COMMAND: Upload file to remote host",
	},
}

func main() {
	initFlag()
	initCfg()
	initHosts()
	initAuths()

	if len(hostgroups.Hgs) == 0 {
		fmt.Println("The hostgroup empty.")
	}

	if *mode == Interactive {
		interactiveRun()
	} else if *mode == SCRIPT {
		scriptRun()
	} else {
		log.Fatal("The mode not supported. not in [interactive/script]")
	}
}

func interactiveRun() {
	for _, value := range hostgroups.Hgs {
		p := prompt.Suggest{
			Text:        value.Groupname,
			Description: "HOSTGROUP",
		}
		ps = append(ps, p)
	}
	fmt.Println("Usage: [hostgroup cmd1; cmd2; cmd3] or [hostgroup download/upload srcFile desDir] or [Ctrl c] exit.")
	var his = []string{}
	for {
		v := prompt.Input(">>> ", hgCompleter, prompt.OptionHistory(his))
		if strings.Trim(string(v), " ") == "" {
			continue
		}
		vs := strings.SplitN(strings.Trim(string(v), " "), " ", 2)
		if len(vs) != 2 {
			fmt.Println("The hostgroup and commands needed.")
			continue
		}
		if vs[0] == "" || vs[1] == "" {
			fmt.Println("The hostgroup and commands needed.")
			continue
		}
		tasks = Tasks{}
		if strings.HasPrefix(vs[1], DOWNLOAD) || strings.HasPrefix(vs[1], UPLOAD) {
			vss := strings.Split(vs[1], " ")
			if len(vss) != 3 {
				fmt.Println("The sftp commands needs. Usage: hostgroup download/upload srcFile desDir")
				continue
			}
			t := Task{
				Taskname:   "DEFAULT",
				Hostgroups: vs[0],
				Sshtasks:   []string{},
				Sftptasks:  []Sftptask{
					Sftptask{
						Type:    vss[0],
						SrcFile: vss[1],
						DesDir:  vss[2],
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
		} else {
			t := Task{
				Taskname:   "DEFAULT",
				Hostgroups: vs[0],
				Sshtasks:   strings.Split(vs[1], ";"),
			}
			tasks.Ts = append(tasks.Ts, t)
		}
		his = append(his, strings.Trim(string(v), " "))
		err := run()
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}

func scriptRun() {
	if *script == "" {
		fmt.Println("The script yaml needed.")
	}
	initArgs(*script)
	err := run()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func run() error {
	if len(tasks.Ts) == 0 {
		return fmt.Errorf("%s", "Tasks empty.")
	}

	limit := make(chan bool, cfg.Concurrency)
	defer close(limit)

	for _, task := range tasks.Ts {
		hg := utils.ChooseHostgroups(hostgroups, task.Hostgroups)
		if hg.Groupname == "" {
			return fmt.Errorf("%s", "The hostgroup not found.")
		}

		auth := utils.ChooseAuthmethod(auths, hg.Authmethod)
		if auth.Name == "" {
			return fmt.Errorf("%s", "The auth method not found.")
		}

		username := auth.Username
		password := auth.Password
		privatekey := auth.Privatekey
		passphrase := auth.Passphrase
		ciphers := auth.Ciphers
		sshport := hg.Sshport
		sudotype := auth.Sudotype
		sudopass := auth.Sudopass

		if len(hg.Hosts) == 0 {
			return fmt.Errorf("%s", "Hosts empty.")
		}
		if sshport <= 0 {
			return fmt.Errorf("%s", "SSH Port not right.")
		}
		if username == "" {
			return fmt.Errorf("%s", "Username empty.")
		}
		if password == "" && privatekey == "" {
			return fmt.Errorf("%s", "Password and PrivateKey empty.")
		}
		if len(task.Sshtasks) == 0 && len(task.Sftptasks) == 0 {
			return fmt.Errorf("%s", "SSH and SFTP Tasks empty.")
		}
		if len(task.Sshtasks) != 0 && !strings.HasPrefix(task.Sshtasks[len(task.Sshtasks)-1], "exit") {
			task.Sshtasks = append(task.Sshtasks, "exit")
		}
		if len(task.Sshtasks) != 0 && sudotype != "" && sudopass != "" {
			task.Sshtasks = append([]string{sudotype, sudopass}, task.Sshtasks...)
			task.Sshtasks = append(task.Sshtasks, "exit")
		}

		taskchs := make(chan Hostresult, len(hg.Hosts))
		var taskresult Taskresult

		for _, host := range hg.Hosts {
			limit <- true
			go func(host string, sshport int, username, password, privatekey, passphrase, sudotype, sudopass string, cipers []string, task Task) {
				var hostresult Hostresult
				hostresult.Hostaddr = host

				if len(task.Sshtasks) != 0 {
					var stdout, stderr string
					var err error
					stdout, stderr, err = lwssh.SSHShell(host, sshport, username, password, privatekey, passphrase, ciphers, task.Sshtasks)
					if err != nil {
						hostresult.Error = err.Error()
					} else {
						hostresult.Stdout = stdout
						hostresult.Stderr = stderr
					}
				}

				if len(task.Sftptasks) != 0 {
					var err error
					for _, value := range task.Sftptasks {
						if value.Type == DOWNLOAD {
							err = lwssh.ScpDownload(host, sshport, username, password, privatekey, passphrase, ciphers, value.SrcFile, path.Join(value.DesDir, hg.Groupname))
							if err == nil {
								hostresult.Stdout = "DOWNLOAD Success."
							}
						} else if value.Type == UPLOAD {
							err = lwssh.ScpUpload(host, sshport, username, password, privatekey, passphrase, ciphers, value.SrcFile, value.DesDir)
							if err == nil {
								hostresult.Stdout = "UPLOAD Success."
							}
						} else {
							err = fmt.Errorf("%s", "Not support scp type, not in [download/upload].")
						}
					}
					if err != nil {
						hostresult.Error = err.Error()
					}
				}

				taskchs <- hostresult
				<-limit
			}(host, sshport, username, password, privatekey, passphrase, sudotype, sudopass, ciphers, task)
		}

		for i := 0; i < len(hg.Hosts); i++ {
			taskresult.Name = task.Taskname
			select {
			case res := <-taskchs:
				taskresult.Results = append(taskresult.Results, res)
			case <-time.After(time.Duration(cfg.Tasktimeout) * time.Second):
				break
			}
		}
		if len(taskresult.Results) != len(hg.Hosts) {
			m := make(map[string]bool)
			for _, v := range taskresult.Results {
				m[v.Hostaddr] = true
			}
			for _, v := range hg.Hosts {
				if _, ok := m[v]; !ok {
					var h Hostresult
					h.Hostaddr = v
					h.Error = "TIMEOUT"
					taskresult.Results = append(taskresult.Results, h)
				}
			}
		}
		utils.Output(taskresult)
		close(taskchs)
	}

	return nil
}
