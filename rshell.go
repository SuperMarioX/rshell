package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/luckywinds/lwssh"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/prompt"
	"github.com/luckywinds/rshell/pkg/utils"
	. "github.com/luckywinds/rshell/types"
	"io"
	"log"
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

var cfg = options.GetCfg()
var auths = options.GetAuths()
var hostgroups = options.GetHostgroups()
var tasks, isScriptMode = options.GetTasks()

func main() {
	if !isScriptMode {
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

func interactiveRun() {
	l, err := prompt.NewReadline(cfg, hostgroups)
	if err != nil {
		log.Fatal(err)
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
			d, h, c, err := utils.GetDo(hostgroups, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			for _, cstr := range strings.Split(c, cfg.CmdSeparator) {
				for _, bstr := range cfg.BlackCmdList {
					if strings.HasPrefix(strings.TrimSpace(cstr), bstr) {
						fmt.Printf("DANGEROUS: Get a blacklist command [%s], DO NOT RUN.\n", strings.TrimSpace(cstr))
						goto retry
					}
				}
			}
			t := Task{
				Taskname:   d,
				Hostgroups: h,
				Sshtasks:   strings.Split(c, cfg.CmdSeparator),
			}
			tasks.Ts = append(tasks.Ts, t)
			for _, value := range t.Sshtasks {
				if strings.TrimSpace(value) != "" {
					prompt.AddCmd(strings.TrimSpace(value))
				}
			}
		case strings.HasPrefix(line, "sudo "):
			s, h, c, err := utils.GetSudo(hostgroups, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			for _, cstr := range strings.Split(c, cfg.CmdSeparator) {
				for _, bstr := range cfg.BlackCmdList {
					if strings.HasPrefix(strings.TrimSpace(cstr), bstr) {
						fmt.Printf("DANGEROUS: Get a blacklist command [%s], DO NOT RUN.\n", strings.TrimSpace(cstr))
						goto retry
					}
				}
			}
			t := Task{
				Taskname:   s,
				Hostgroups: h,
				Sudoroot:   true,
				Sshtasks:   strings.Split(c, cfg.CmdSeparator),
			}
			tasks.Ts = append(tasks.Ts, t)
			for _, value := range t.Sshtasks {
				if strings.TrimSpace(value) != "" {
					prompt.AddCmd(strings.TrimSpace(value))
				}
			}
		case strings.HasPrefix(line, "download "):
			d, h, sf, dd, err := utils.GetDownload(hostgroups, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			t := Task{
				Taskname:   d,
				Hostgroups: h,
				Sftptasks: []Sftptask{
					{
						Type:    d,
						SrcFile: sf,
						DesDir:  dd,
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
			prompt.AddSrcFile(strings.TrimSpace(sf))
			prompt.AddDesDir(strings.TrimSpace(dd))
		case strings.HasPrefix(line, "upload "):
			u, h, sf, dd, err := utils.GetUpload(hostgroups, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			t := Task{
				Taskname:   u,
				Hostgroups: h,
				Sftptasks: []Sftptask{
					{
						Type:    u,
						SrcFile: sf,
						DesDir:  dd,
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
			prompt.AddSrcFile(strings.TrimSpace(sf))
			prompt.AddDesDir(strings.TrimSpace(dd))
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
		defer close(taskchs)

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
					}
					hostresult.Stdout = stdout
					hostresult.Stderr = stderr
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

		var newtaskresult Taskresult
		newtaskresult.Name = taskresult.Name
		for _, h := range hg.Hosts {
			for _, r := range taskresult.Results {
				if r.Hostaddr == h {
					newtaskresult.Results = append(newtaskresult.Results, r)
					break
				}
			}
		}

		utils.Output(newtaskresult)
	}

	return nil
}
