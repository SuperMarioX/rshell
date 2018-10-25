package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/deckarep/golang-set"
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
	SUDO        string = "sudo"
)

var (
	cfgFile   = flag.String("cfg", path.Join(".rshell", "cfg.yaml"), "System Config file to read, Default: " + path.Join(".rshell", "cfg.yaml"))
	hostsFile = flag.String("hosts", path.Join(".rshell", "hosts.yaml"), "Hosts Config file to read, Default: " + path.Join(".rshell", "hosts.yaml"))
	authFile  = flag.String("auth", path.Join(".rshell", "auth.yaml"), "Auth Config file to read, Default: " + path.Join(".rshell", "hosts.yaml"))
	script    = flag.String("f", "", "The script yaml.")
)

func initFlag() {
	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`Usage:`, cmd, `[<options>]

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
	return prompt.FilterHasPrefix(pssl, d.GetWordBeforeCursor(), true)
}

var ps = []prompt.Suggest{
	{
		Text: DOWNLOAD,
		Description: "KEYWORDS: Download file from remote host",
	},
	{
		Text: UPLOAD,
		Description: "KEYWORDS: Upload file to remote host",
	},
	{
		Text: SUDO,
		Description: "KEYWORDS: auto change root",
	},
}

var pss = mapset.NewSet()
var pssl = []prompt.Suggest{}

func main() {
	os.Mkdir(".rshell", os.ModeDir)
	initFlag()
	initCfg()
	initHosts()
	initAuths()

	if len(hostgroups.Hgs) == 0 {
		fmt.Println("The hostgroup empty.")
	}

	if *script == "" {
		interactiveRun()
	} else {
		scriptRun()
	}
}

func showInteractiveRunUsage() {
	fmt.Println(`Usage:
    SSH Task: hostgroup cmd1; cmd2; cmd3
    SSH Task(sudo root): hostgroup sudo cmd1; cmd2; cmd3
    SFTP Task: hostgroup download/upload srcFile desDir
    Exit: Ctrl c
    Help: ?`)
}

func interactiveRun() {
	showInteractiveRunUsage()
	for _, value := range ps {
		pss.Add(value)
	}
	for _, value := range hostgroups.Hgs {
		p := prompt.Suggest{
			Text:        value.Groupname,
			Description: "HOSTGROUP",
		}
		pss.Add(p)
	}

	var his = []string{}
	var sug = prompt.Suggest{}
	for {
		pssl = []prompt.Suggest{}
		for _, value := range pss.ToSlice() {
			v := value.(prompt.Suggest)
			pssl = append(pssl, v)
		}
		v := prompt.Input("rshell: ", hgCompleter, prompt.OptionHistory(his), prompt.OptionMaxSuggestion(6))
		if strings.Trim(string(v), " ") == "" {
			continue
		}
		if strings.Trim(string(v), " ") == "?" {
			showInteractiveRunUsage()
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
		if strings.HasPrefix(vs[1], DOWNLOAD + " ") || strings.HasPrefix(vs[1], UPLOAD + " ") {
			vss := strings.Split(vs[1], " ")
			if len(vss) != 3 {
				fmt.Println("The sftp commands needs.")
				showInteractiveRunUsage()
				continue
			}
			t := Task{
				Taskname:   "DEFAULT",
				Hostgroups: vs[0],
				Sshtasks:   []string{},
				Sftptasks:  []Sftptask{
					{
						Type:    vss[0],
						SrcFile: vss[1],
						DesDir:  vss[2],
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
			sug = prompt.Suggest{
				Text:        vss[1],
				Description: "FILE",
			}
			pss.Add(sug)
			sug = prompt.Suggest{
				Text:        vss[2],
				Description: "DIRECTORY",
			}
			pss.Add(sug)
		} else if strings.HasPrefix(vs[1], "sudo") {
			vss := strings.SplitN(vs[1], " ", 2)
			if len(vss) <= 1 || strings.Trim(vss[1], " ") == "" {
				fmt.Println("The ssh sudo commands needs.")
				showInteractiveRunUsage()
				continue
			}
			t := Task{
				Taskname:   "DEFAULT",
				Hostgroups: vs[0],
				Sudoroot: true,
				Sshtasks:   strings.Split(vss[1], ";"),
			}
			tasks.Ts = append(tasks.Ts, t)
			for _, value := range t.Sshtasks {
				if strings.Trim(value, " ") == "" {
					continue
				}
				sug = prompt.Suggest{
					Text:        strings.Trim(value, " "),
					Description: "COMMAND",
				}
				pss.Add(sug)
			}
		} else {
			t := Task{
				Taskname:   "DEFAULT",
				Hostgroups: vs[0],
				Sshtasks:   strings.Split(vs[1], ";"),
			}
			tasks.Ts = append(tasks.Ts, t)
			for _, value := range t.Sshtasks {
				if strings.Trim(value, " ") == "" {
					continue
				}
				sug = prompt.Suggest{
					Text:        strings.Trim(value, " "),
					Description: "COMMAND",
				}
				pss.Add(sug)
			}
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
		if len(task.Sshtasks) != 0 && task.Sshtasks[len(task.Sshtasks)-1] != "exit" && !strings.HasPrefix(task.Sshtasks[len(task.Sshtasks)-1], "exit ") {
			task.Sshtasks = append(task.Sshtasks, "exit")
		}
		if len(task.Sshtasks) != 0 && task.Sudoroot {
			if sudotype == "" {
				sudotype = "su"
			}
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
					hostresult.Stdout += "\n"
					var err error
					for _, value := range task.Sftptasks {
						if value.Type == DOWNLOAD {
							err = lwssh.ScpDownload(host, sshport, username, password, privatekey, passphrase, ciphers, value.SrcFile, path.Join(value.DesDir, hg.Groupname))
							if err == nil {
								hostresult.Stdout += "DOWNLOAD Success [" + value.SrcFile + " -> " + path.Join(value.DesDir, hg.Groupname) + "]\n"
							} else {
								hostresult.Error += "DOWNLOAD Failed [" + value.SrcFile + " -> " + path.Join(value.DesDir, hg.Groupname) + "] " + err.Error() + "\n"
							}
						} else if value.Type == UPLOAD {
							err = lwssh.ScpUpload(host, sshport, username, password, privatekey, passphrase, ciphers, value.SrcFile, value.DesDir)
							if err == nil {
								hostresult.Stdout += "UPLOAD Success [" + value.SrcFile + " -> " + value.DesDir + "]\n"
							} else {
								hostresult.Error += "UPLOAD Failed [" + value.SrcFile + " -> " + value.DesDir + "] " + err.Error() + "\n"
							}
						} else {
							err = fmt.Errorf("%s", "")
							hostresult.Error += "SFTP Failed. Not support type[" + value.Type + "], not in [download/upload].\n"
						}
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
