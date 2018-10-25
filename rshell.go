package main

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/luckywinds/lwssh"
	. "github.com/luckywinds/rshell/types"
	"github.com/luckywinds/rshell/utils"
	"gopkg.in/yaml.v2"
	"io"
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
	cfgFile   = flag.String("cfg", path.Join(".rshell", "cfg.yaml"), "System Config file to read, Default: "+path.Join(".rshell", "cfg.yaml"))
	hostsFile = flag.String("hosts", path.Join(".rshell", "hosts.yaml"), "Hosts Config file to read, Default: "+path.Join(".rshell", "hosts.yaml"))
	authFile  = flag.String("auth", path.Join(".rshell", "auth.yaml"), "Auth Config file to read, Default: "+path.Join(".rshell", "hosts.yaml"))
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
	fmt.Println(`Usage: <keywords> <hostgroup> <agruments>

do hostgroup cmd1; cmd2; cmd3
    --- Run cmds on hostgroup use normal user
sudo hostgroup sudo cmd1; cmd2; cmd3
    --- Run cmds on hostgroup use root which auto change from normal user
download hostgroup srcFile desDir
    --- Download srcFile from hostgroup to local desDir
upload hostgroup srcFile desDir
    --- Upload srcFile from local to hostgroup desDir
ctrl c
    --- Exit
?
    --- Help`)
}

func listHgs() func(string) []string {
	return func(line string) []string {
		hgs := []string{}
		for _, value := range hostgroups.Hgs {
			hgs = append(hgs, value.Groupname)
		}
		return hgs
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("do",
		readline.PcItemDynamic(listHgs()),
	),
	readline.PcItem("sudo",
		readline.PcItemDynamic(listHgs()),
	),
	readline.PcItem("download",
		readline.PcItemDynamic(listHgs()),
	),
	readline.PcItem("upload",
		readline.PcItemDynamic(listHgs()),
	),
)

func interactiveRun() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31mrshell:\033[0m ",
		HistoryFile:         ".rshell/rshell.history",
		AutoComplete:        completer,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		retry:
		tasks = Tasks{}
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "do "):
			if err := utils.CheckDo(hostgroups, line); err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				break
			}
			lines := strings.SplitN(line, " ", 3)
			t := Task{
				Taskname:   "do",
				Hostgroups: lines[1],
				Sshtasks:   strings.Split(lines[2], ";"),
			}
			tasks.Ts = append(tasks.Ts, t)
		case strings.HasPrefix(line, "sudo "):
			if err := utils.CheckSudo(hostgroups, line); err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				break
			}
			lines := strings.SplitN(line, " ", 3)
			t := Task{
				Taskname:   "sudo",
				Hostgroups: lines[1],
				Sudoroot:   true,
				Sshtasks:   strings.Split(lines[2], ";"),
			}
			tasks.Ts = append(tasks.Ts, t)
		case strings.HasPrefix(line, "download "):
			if err := utils.CheckDownload(hostgroups, line); err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				break
			}
			lines := strings.SplitN(line, " ", 4)
			t := Task{
				Taskname:   "download",
				Hostgroups: lines[1],
				Sftptasks: []Sftptask{
					{
						Type:    "download",
						SrcFile: lines[2],
						DesDir:  lines[3],
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
		case strings.HasPrefix(line, "upload "):
			if err := utils.CheckUpload(hostgroups, line); err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				break
			}
			lines := strings.SplitN(line, " ", 4)
			t := Task{
				Taskname:   "upload",
				Hostgroups: lines[1],
				Sftptasks: []Sftptask{
					{
						Type:    "upload",
						SrcFile: lines[2],
						DesDir:  lines[3],
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
		case line == "?":
			showInteractiveRunUsage()
			goto retry
		case line == "":
			goto retry
		default:
			showInteractiveRunUsage()
			goto retry
		}

		if err := run(); err != nil {
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
		if sshport < 0 {
			return fmt.Errorf("%s", "SSH Port not right.")
		}
		if sshport == 0 {
			sshport = 22
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
