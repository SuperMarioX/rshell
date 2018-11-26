package main

import (
	"fmt"
	"github.com/luckywinds/lwssh"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/crypt"
	"github.com/luckywinds/rshell/pkg/filters"
	"github.com/luckywinds/rshell/pkg/prompt"
	"github.com/luckywinds/rshell/pkg/update"
	"github.com/luckywinds/rshell/pkg/utils"
	. "github.com/luckywinds/rshell/types"
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
	SSH         string = "ssh"
	SFTP        string = "sftp"
	AES         string = "aes"
)

var cfg = options.GetCfg()
var auths, authsMap = options.GetAuths()
var hostgroups, hostgroupsMap = options.GetHostgroups()
var tasks, isScriptMode = options.GetTasks()
var version = "3.5"

func main() {
	go update.Update(cfg, version)

	if !isScriptMode {
		interactiveRun()
	} else {
		scriptRun()
	}
}

func showIntro() {
	fmt.Println(`
 ______     ______     __  __     ______     __         __
/\  == \   /\  ___\   /\ \_\ \   /\  ___\   /\ \       /\ \
\ \  __<   \ \___  \  \ \  __ \  \ \  __\   \ \ \____  \ \ \____
 \ \_\ \_\  \/\_____\  \ \_\ \_\  \ \_____\  \ \_____\  \ \_____\
  \/_/ /_/   \/_____/   \/_/\/_/   \/_____/   \/_____/   \/_____/
------ Rshell @`+version+` Type "?" or "help" for more information. -----`)
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

encrypt_aes cleartext_password
    --- Encrypt cleartext_password with aes 256 cfb
decrypt_aes ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb

exit
    --- Exit rshell
?
    --- Help`)
}

func interactiveRun() {
	showIntro()
	l, err := prompt.New(cfg, hostgroups)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
	retry:
		tasks = Tasks{}
		line, err := prompt.Prompt(l, cfg)
		if err == prompt.ErrPromptAborted {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == prompt.ErrPromptError {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "do "):
			d, h, c, err := utils.GetDo(hostgroupsMap, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			t := Task{
				Name:      d,
				Hostgroup: h,
				Subtasks: []Subtask{
					{
						Name: "DEFAULT",
						Mode: SSH,
						Sudo: false,
						Cmds: strings.Split(c, cfg.CmdSeparator),
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
			for _, value := range strings.Split(c, cfg.CmdSeparator) {
				if strings.TrimSpace(value) != "" {
					prompt.AddCmd(strings.TrimSpace(value))
				}
			}
		case strings.HasPrefix(line, "sudo "):
			s, h, c, err := utils.GetSudo(hostgroupsMap, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			t := Task{
				Name:      s,
				Hostgroup: h,
				Subtasks: []Subtask{
					{
						Name: "DEFAULT",
						Mode: SSH,
						Sudo: true,
						Cmds: strings.Split(c, cfg.CmdSeparator),
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
			for _, value := range strings.Split(c, cfg.CmdSeparator) {
				if strings.TrimSpace(value) != "" {
					prompt.AddCmd(strings.TrimSpace(value))
				}
			}
		case strings.HasPrefix(line, "download "):
			d, h, sf, dd, err := utils.GetDownload(hostgroupsMap, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			t := Task{
				Name:      d,
				Hostgroup: h,
				Subtasks: []Subtask{
					{
						Name:    "DEFAULT",
						Mode:    SFTP,
						FtpType: DOWNLOAD,
						SrcFile: sf,
						DesDir:  dd,
					},
				},
			}
			tasks.Ts = append(tasks.Ts, t)
			prompt.AddSrcFile(strings.TrimSpace(sf))
			prompt.AddDesDir(strings.TrimSpace(dd))
		case strings.HasPrefix(line, "upload "):
			u, h, sf, dd, err := utils.GetUpload(hostgroupsMap, line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				showInteractiveRunUsage()
				goto retry
			}
			t := Task{
				Name:      u,
				Hostgroup: h,
				Subtasks: []Subtask{
					{
						Name:    "DEFAULT",
						Mode:    SFTP,
						FtpType: UPLOAD,
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
		case line == "exit":
			return
		case strings.HasPrefix(line, "encrypt_aes "):
			pass, err := utils.GetEncypt(line)
			if err != nil {
				fmt.Printf("Encypt failed. %v\n", err)
			} else {
				pp, err := crypt.AesEncrypt(pass, cfg)
				if err != nil {
					fmt.Printf("Encypt failed. %v\n", err)
				} else {
					fmt.Println(pp)
				}
			}
			goto retry
		case strings.HasPrefix(line, "decrypt_aes "):
			pass, err := utils.GetEncypt(line)
			if err != nil {
				fmt.Printf("Decypt failed. %v\n", err)
			} else {
				pp, err := crypt.AesDecrypt(pass, cfg)
				if err != nil {
					fmt.Printf("Decypt failed. %v\n", err)
				} else {
					fmt.Println(pp)
				}
			}
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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(tasks.Ts) == 0 {
		return fmt.Errorf("%s", "Tasks empty.")
	}

	limit := make(chan bool, cfg.Concurrency)
	defer close(limit)

	for _, task := range tasks.Ts {
		if task.Name == "" || task.Hostgroup == "" {
			log.Fatal("The task's name or hostgroup empty.")
		}

		hg := hostgroupsMap[task.Hostgroup]
		if hg.Groupname == "" {
			return fmt.Errorf("%s", "The hostgroup not found.")
		}

		auth := authsMap[hg.Authmethod]
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
		if sudotype == "" {
			sudotype = "su"
		}
		sudopass := auth.Sudopass

		if cfg.Passcrypttype == AES {
			var perr error
			password, perr = crypt.AesDecrypt(password, cfg)
			if perr != nil {
				return fmt.Errorf("%s/%s/%s", task.Hostgroup, auth.Name, "password format not right.")
			}
			passphrase, perr = crypt.AesDecrypt(passphrase, cfg)
			if perr != nil {
				return fmt.Errorf("%s/%s/%s", task.Hostgroup, auth.Name, "passphrase format not right.")
			}
			sudopass, perr = crypt.AesDecrypt(sudopass, cfg)
			if perr != nil {
				return fmt.Errorf("%s/%s/%s", task.Hostgroup, auth.Name, "sudopass format not right.")
			}
		}

		if len(hg.Ips) == 0 {
			return fmt.Errorf("%s", "Hostgroup empty.")
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
		if len(task.Subtasks) == 0 {
			return fmt.Errorf("%s", "SSH or SFTP Tasks empty.")
		}

		taskchs := make(chan Hostresult, len(hg.Ips))
		defer close(taskchs)

		var taskresult Taskresult
		for _, host := range hg.Ips {
			limit <- true
			go func(host string, sshport int, username, password, privatekey, passphrase, sudotype, sudopass string, cipers []string, task Task) {
				var hostresult Hostresult
				hostresult.Hostaddr = host

				var stdout, stderr string
				var err error
				for _, item := range task.Subtasks {
					if item.Mode == SSH {
						if isScriptMode {
							hostresult.Stdout += fmt.Sprintf("SSH  [%-16s] -------------------------------------------------------\n", item.Name)
						}
						if len(item.Cmds) == 0 {
							hostresult.Error += "The SSH cmds empty.\n"
							break
						}
						for _, cstr := range item.Cmds {
							if filters.IsBlackCmd(cstr, cfg.BlackCmdList) {
								hostresult.Error += fmt.Sprintf("DANGEROUS: Get a blacklist command [%s], DO NOT RUN.\n", cstr)
								goto breakfor
							}
						}
						if item.Cmds[len(item.Cmds)-1] != "exit" && !strings.HasPrefix(item.Cmds[len(item.Cmds)-1], "exit ") {
							item.Cmds = append(item.Cmds, "exit")
						}
						if item.Sudo {
							item.Cmds = append([]string{sudotype, sudopass}, item.Cmds...)
							item.Cmds = append(item.Cmds, "exit")
						}
						stdout, stderr, err = lwssh.SSHShell(host, sshport, username, password, privatekey, passphrase, ciphers, item.Cmds)
						hostresult.Stdout += stdout
						if stderr != "" {
							if isScriptMode {
								hostresult.Stdout += "\n"
								hostresult.Stdout += "WARN : COMMANDS Return @ " + stderr + "\n"
							} else {
								hostresult.Stderr += stderr
							}
						}
						if err != nil {
							hostresult.Error += "Finish Now @ " + err.Error()
							break
						}
					} else if item.Mode == SFTP {
						if isScriptMode {
							hostresult.Stdout += fmt.Sprintf("FTP  [%-16s] -------------------------------------------------------\n", item.Name)
						}
						if item.FtpType == DOWNLOAD {
							err = lwssh.ScpDownload(host, sshport, username, password, privatekey, passphrase, ciphers, item.SrcFile, path.Join(item.DesDir, hg.Groupname))
							if err == nil {
								hostresult.Stdout += "DOWNLOAD Success [" + item.SrcFile + " -> " + path.Join(item.DesDir, hg.Groupname) + "]\n\n"
							} else {
								hostresult.Error += "Finish Now @ DOWNLOAD Failed [" + item.SrcFile + " -> " + path.Join(item.DesDir, hg.Groupname) + "] " + err.Error()
								break
							}
						} else if item.FtpType == UPLOAD {
							err = lwssh.ScpUpload(host, sshport, username, password, privatekey, passphrase, ciphers, item.SrcFile, item.DesDir)
							if err == nil {
								hostresult.Stdout += "UPLOAD Success [" + item.SrcFile + " -> " + item.DesDir + "]\n\n"
							} else {
								hostresult.Error += "Finish Now @ UPLOAD Failed [" + item.SrcFile + " -> " + item.DesDir + "] " + err.Error()
								break
							}
						} else {
							hostresult.Error += "Finish Now @ SFTP Failed. Not support ftp type[" + item.FtpType + "], not in [download/upload].\n"
							break
						}
					} else {
						hostresult.Error += "Finish Now @ TASK Failed. Not support task mode[" + item.Mode + "], not in [ssh/sftp].\n"
						break
					}
				}
			breakfor:

				taskchs <- hostresult
				<-limit
			}(host, sshport, username, password, privatekey, passphrase, sudotype, sudopass, ciphers, task)
		}

		if cfg.Outputintime {
			utils.OutputTaskHeader(task.Name + "@" + task.Hostgroup)
		}

		for i := 0; i < len(hg.Ips); i++ {
			taskresult.Name = task.Name + "@" + task.Hostgroup
			select {
			case res := <-taskchs:
				taskresult.Results = append(taskresult.Results, res)
				if cfg.Outputintime {
					utils.OutputHostResult(res)
				}
			case <-time.After(time.Duration(cfg.Tasktimeout) * time.Second):
				break
			}
		}

		m := make(map[string]Hostresult)
		for _, v := range taskresult.Results {
			m[v.Hostaddr] = v
		}

		var newtaskresult Taskresult
		newtaskresult.Name = taskresult.Name
		for _, h := range hg.Ips {
			if _, ok := m[h]; ok {
				newtaskresult.Results = append(newtaskresult.Results, m[h])
			} else {
				t := Hostresult{
					Hostaddr: h,
					Error:    "TIMEOUT",
					Stdout:   "",
					Stderr:   "",
				}
				newtaskresult.Results = append(newtaskresult.Results, t)
				if cfg.Outputintime {
					utils.OutputHostResult(t)
				}
			}
		}

		if !cfg.Outputintime {
			utils.Output(newtaskresult)
		}
	}

	return nil
}
