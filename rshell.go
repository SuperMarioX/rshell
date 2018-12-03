package main

import (
	"fmt"
	"github.com/luckywinds/rshell/commands"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/pkg/prompt"
	"github.com/luckywinds/rshell/pkg/update"
	"github.com/luckywinds/rshell/pkg/utils"
	"log"
	"strings"
)

const (
	Interactive string = "interactive"
	SCRIPT      string = "script"
	DOWNLOAD    string = "download"
	UPLOAD      string = "upload"
	SSH         string = "ssh"
	SFTP        string = "sftp"
	AES         string = "aes"
	DO          string = "do"
	SUDO        string = "sudo"
)

var cfg = options.GetCfg()
var auths, authsMap = options.GetAuths()
var hostgroups, hostgroupsMap = options.GetHostgroups()
var tasks, isScriptMode = options.GetTasks()

func main() {
	go update.Update(cfg, version)

	if !isScriptMode {
		interactiveRun()
	} else {
		scriptRun()
	}
}

var version = "5.0"
func showIntro() {
	fmt.Println(`
 ______     ______     __  __     ______     __         __
/\  == \   /\  ___\   /\ \_\ \   /\  ___\   /\ \       /\ \
\ \  __<   \ \___  \  \ \  __ \  \ \  __\   \ \ \____  \ \ \____
 \ \_\ \_\  \/\_____\  \ \_\ \_\  \ \_____\  \ \_____\  \ \_____\
  \/_/ /_/   \/_____/   \/_/\/_/   \/_____/   \/_____/   \/_____/
------ Rshell @`+version+` Type "?" or "help" for more information. -----`)
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
			_, h, c, err := utils.GetSSHArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			if err = commands.DO("do", h, c); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			for _, value := range c {
				prompt.AddCmd(strings.Trim(value, " "))
			}
		case strings.HasPrefix(line, "sudo "):
			_, h, c, err := utils.GetSSHArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			if err = commands.SUDO("sudo", h, c); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			for _, value := range c {
				prompt.AddCmd(strings.Trim(value, " "))
			}
		case strings.HasPrefix(line, "download "):
			_, h, sf, dd, err := utils.GetSFTPArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			if err = commands.Download("download", h, sf, dd); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			prompt.AddSrcFile(strings.Trim(sf, " "))
			prompt.AddDesDir(strings.Trim(dd, " "))
		case strings.HasPrefix(line, "upload "):
			_, h, sf, dd, err := utils.GetSFTPArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			if err = commands.Upload("upload", h, sf, dd); err != nil {
				fmt.Printf("%v\n", err.Error())
				commands.Help()
				goto retry
			}
			prompt.AddSrcFile(strings.Trim(sf, " "))
			prompt.AddDesDir(strings.Trim(dd, " "))
		case strings.HasPrefix(line, "encrypt_aes "):
			_, t, err := utils.GetCryptArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.Help()
				goto retry
			}
			p, err := commands.Encrypt(t)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.Help()
				goto retry
			}
			fmt.Println(p)
		case strings.HasPrefix(line, "decrypt_aes "):
			_, t, err := utils.GetCryptArgs(line)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.Help()
				goto retry
			}
			p, err := commands.Decrypt(t)
			if err != nil {
				fmt.Printf("%v\n", err)
				commands.Help()
				goto retry
			}
			fmt.Println(p)
		case line == "?":
			commands.Help()
			goto retry
		case line == "":
			goto retry
		case line == "exit":
			return
		default:
			commands.Help()
			goto retry
		}
	}
}

func scriptRun() {
	for _, task := range tasks.Ts {
		if task.Name == "" || task.Hostgroup == "" {
			log.Fatal("The task's name or hostgroup empty.")
		}

		if len(task.Subtasks) == 0 {
			log.Fatal("SSH or SFTP Tasks empty.")
		}

		for _, stask := range task.Subtasks {
			name := task.Name + "/" + stask.Name
			if stask.Mode == SSH {
				if stask.Sudo {
					if err := commands.SUDO(name, task.Hostgroup, stask.Cmds); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, SUDO, err)
					}
				} else {
					if err := commands.DO(name, task.Hostgroup, stask.Cmds); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, DO, err)
					}
				}
			} else if stask.Mode == SFTP {
				if stask.FtpType == DOWNLOAD {
					if err := commands.Download(name, task.Hostgroup, stask.SrcFile, stask.DesDir); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, DOWNLOAD, err)
					}
				} else if stask.FtpType == UPLOAD {
					if err := commands.Upload(name, task.Hostgroup, stask.SrcFile, stask.DesDir); err != nil {
						log.Fatalf("ERROR: %s/%s/%s/%v", name, task.Hostgroup, UPLOAD, err)
					}
				} else {
					log.Fatalf("ERROR: %s/%s/%s/%s", name, task.Hostgroup, stask.FtpType, "Not support")
				}
			} else {
				log.Fatalf("ERROR: %s/%s/%s/%s", name, task.Hostgroup, stask.Mode, "Not support")
			}
		}
	}
}
