package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/types"
	"strings"
)

func ChooseHostgroups(hgs types.Hostgroups, groupname string) (types.Hostgroup) {
	var result types.Hostgroup
	for _, hg := range hgs.Hgs {
		if hg.Groupname == groupname {
			result = hg
			break
		}
	}
	return result
}

func ChooseAuthmethod(as types.Auths, methodname string) (types.Auth) {
	var result types.Auth
	for _, a := range as.As {
		if a.Name == methodname {
			result = a
			break
		}
	}
	return result
}

func Output(result types.Taskresult) {
	color.Yellow("TASK [%-16s] ********************************************************\n", result.Name)
	for _, ret := range result.Results {
		color.Green("HOST [%-16s] --------------------------------------------------------\n", ret.Hostaddr)
		color.Magenta("%s\n","STDOUT =>")
		if ret.Stdout != "" {
			fmt.Printf("%s\n", ret.Stdout)
		}
		if ret.Stderr != "" {
			color.Red("%s\n", "STDERR =>")
			fmt.Printf("%s\n", ret.Stderr)
		}
		if ret.Error != "" {
			color.Red("%s\n", "ERROR =>")
			fmt.Printf("%s\n", ret.Error)
		}
	}
}

func CheckDo(hgs types.Hostgroups, line string) error {
	if len(line) <= 3 {
		return fmt.Errorf("do need arguments")
	}
	ls := strings.SplitN(line, " ", 3)
	if len(ls) != 3 {
		return fmt.Errorf("do need arguments")
	}
	hg := ChooseHostgroups(hgs, ls[1])
	if hg.Groupname == "" {
		return fmt.Errorf("do need hostgroup")
	}
	return nil
}

func CheckSudo(hgs types.Hostgroups, line string) error {
	if len(line) <= 5 {
		return fmt.Errorf("sudo need arguments")
	}
	ls := strings.SplitN(line, " ", 3)
	if len(ls) != 3 {
		return fmt.Errorf("sudo need arguments")
	}
	hg := ChooseHostgroups(hgs, ls[1])
	if hg.Groupname == "" {
		return fmt.Errorf("sudo need hostgroup")
	}
	return nil
}

func CheckDownload(hgs types.Hostgroups, line string) error {
	if len(line) <= 9 {
		return fmt.Errorf("download need arguments")
	}
	ls := strings.Split(line, " ")
	if len(ls) != 4 {
		return fmt.Errorf("download need arguments")
	}
	hg := ChooseHostgroups(hgs, ls[1])
	if hg.Groupname == "" {
		return fmt.Errorf("download need hostgroup")
	}
	return nil
}

func CheckUpload(hgs types.Hostgroups, line string) error {
	if len(line) <= 7 {
		return fmt.Errorf("upload need arguments")
	}
	ls := strings.Split(line, " ")
	if len(ls) != 4 {
		return fmt.Errorf("upload need arguments")
	}
	hg := ChooseHostgroups(hgs, ls[1])
	if hg.Groupname == "" {
		return fmt.Errorf("upload need hostgroup")
	}
	return nil
}
