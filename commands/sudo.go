package commands

import (
	"fmt"
	"github.com/luckywinds/rshell/modes/ssh"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/outputs"
	"github.com/luckywinds/rshell/pkg/filters"
	"github.com/luckywinds/rshell/types"
)

func SUDO(actionname, groupname string, cmds []string) error {
	if len(cmds) == 0 {
		return fmt.Errorf("cmds[%v] empty", cmds)
	}

	hg, au, err := getGroupAuthbyGroupname(groupname)
	if err != nil {
		return err
	}

	cfg := options.GetCfg()

	for _, c := range cmds {
		if filters.IsBlackCmd(c, cfg.BlackCmdList) {
			return fmt.Errorf("DANGER: [%s] is in black command list.", c)
		}
	}

	if cfg.Passcrypttype != "" {
		au.Password, err = getPlainPass(au.Password, cfg)
		if err != nil {
			return err
		}
		au.Passphrase, err = getPlainPass(au.Passphrase, cfg)
		if err != nil {
			return err
		}
		au.Sudopass, err = getPlainPass(au.Sudopass, cfg)
		if err != nil {
			return err
		}
	}

	limit := make(chan bool, cfg.Concurrency)
	defer close(limit)

	taskchs := make(chan types.Hostresult, len(hg.Ips))
	defer close(taskchs)

	for _, ip := range hg.Ips {
		limit <- true
		go func(groupname, host string, port int, user, pass, keyname, passphrase, sudotype, sudopass string, timeout int, ciphers, cmds []string) {
			//fmt.Printf("%v %v %v %v %v %v %v %v %v %v\n", groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers, cmds)
			stdout, stderr, err := ssh.SUDO(groupname, host, port, user, pass, keyname, passphrase, sudotype, sudopass, timeout, ciphers, cmds)
			var result types.Hostresult
			result.Actionname = actionname
			result.Actiontype = "sudo"
			result.Groupname = groupname
			result.Hostaddr = ip
			result.Stdout = stdout
			result.Stderr = stderr
			if err != nil {
				result.Error = err.Error()
			}
			taskchs <- result
			<-limit
		}(groupname, ip, hg.Sshport, au.Username, au.Password, au.Privatekey, au.Passphrase, au.Sudotype, au.Sudopass, cfg.Tasktimeout, []string{}, cmds)
	}

	outputs.Output(taskchs, hg)
	return nil
}
