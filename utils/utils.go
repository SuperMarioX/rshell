package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/types"
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
		if ret.Stdout == "" && ret.Stderr == "" && ret.Error == "" {
			fmt.Println()
		}
	}
}