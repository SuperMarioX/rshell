package commands

import (
	"fmt"
	"github.com/luckywinds/rshell/modes/sftp"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/outputs"
	"github.com/luckywinds/rshell/types"
	"path"
	"strings"
)

func Download(actionname, groupname string, srcFilePath, desDirPath string) error {
	if srcFilePath == "" || desDirPath == "" {
		return fmt.Errorf("srcFilePath[%s] or desDirPath[%s] empty", srcFilePath, desDirPath)
	}

	hg, au, err := getGroupAuthbyGroupname(groupname)
	if err != nil {
		return err
	}

	cfg := options.GetCfg()

	if cfg.Passcrypttype != "" {
		au.Password, err = getPlainPass(au.Password, cfg)
		if err != nil {
			return fmt.Errorf("get plain password error [%v] crypt type is [%s]", err, cfg.Passcrypttype)
		}
		au.Passphrase, err = getPlainPass(au.Passphrase, cfg)
		if err != nil {
			return fmt.Errorf("get plain password error [%v] crypt type is [%s]", err, cfg.Passcrypttype)
		}
		au.Sudopass, err = getPlainPass(au.Sudopass, cfg)
		if err != nil {
			return fmt.Errorf("get plain password error [%v] crypt type is [%s]", err, cfg.Passcrypttype)
		}
	}

	limit := make(chan bool, cfg.Concurrency)
	defer close(limit)

	taskchs := make(chan types.Hostresult, len(hg.Ips))
	defer close(taskchs)

	for _, ip := range hg.Ips {
		limit <- true
		go func(groupname, host string, port int, user, pass, keyname, passphrase string, timeout int, ciphers []string, srcFilePath, desDirPath string) {
			//fmt.Printf("%v %v %v %v %v %v %v %v %v %v\n", groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers, cmds)
			if !strings.HasSuffix(desDirPath, "/") {
				desDirPath = desDirPath + "/"
			}
			err := sftp.Download(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers, srcFilePath, desDirPath)
			var result types.Hostresult
			result.Actionname = actionname
			result.Actiontype = "download"
			result.Groupname = groupname
			result.Hostaddr = host
			if err == nil {
				result.Stdout = "DOWNLOAD Success [" + srcFilePath + " -> " + path.Join(desDirPath, hg.Groupname) + "]"
			} else {
				result.Stderr = "DOWNLOAD Failed [" + srcFilePath + " -> " + path.Join(desDirPath, hg.Groupname) + "] " + err.Error()
			}
			taskchs <- result
			<-limit
		}(groupname, ip, hg.Sshport, au.Username, au.Password, au.Privatekey, au.Passphrase, cfg.Tasktimeout, []string{}, srcFilePath, desDirPath)
	}

	outputs.Output(taskchs, hg)
	return nil
}
