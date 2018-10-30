package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/types"
	"strings"
)

func OutputTaskHeader(name string) {
	color.Yellow("TASK [%-16s] *******************************************************\n", name)
}
func OutputHostResult(result types.Hostresult) {
	color.Green("HOST [%-16s] =======================================================\n", result.Hostaddr)
	if result.Stdout != "" {
		fmt.Printf("%s\n", result.Stdout)
	}
	if result.Stderr != "" {
		color.Red("%s\n", "STDERR =>")
		fmt.Printf("%s\n", result.Stderr)
	}
	if result.Error != "" {
		color.Red("%s\n", "SYSERR =>")
		fmt.Printf("%s\n", result.Error)
	}
	if result.Stdout == "" && result.Stderr == "" && result.Error == "" {
		fmt.Println()
	}
}

func Output(result types.Taskresult) {
	OutputTaskHeader(result.Name)
	for _, ret := range result.Results {
		OutputHostResult(ret)
	}
}

func GetDo(hgm map[string]types.Hostgroup, line string) (do, hostgroup, cmd string, err error) {
	ds := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(ds) != 2 {
		return "", "", "", fmt.Errorf("%s", "do need arguments")
	}
	do = ds[0]

	hs := strings.SplitN(strings.TrimSpace(ds[1]), " ", 2)
	if len(hs) != 2 {
		return "", "", "", fmt.Errorf("%s", "do need arguments")
	}
	hg := hgm[hs[0]]
	if hg.Groupname == "" {
		return "", "", "", fmt.Errorf("do need hostgroup")
	}
	hostgroup = hs[0]

	cmd = strings.TrimSpace(hs[1])
	if cmd == "" {
		return "", "", "", fmt.Errorf("do need command")
	}

	return do, hostgroup, cmd, nil
}

func GetSudo(hgm map[string]types.Hostgroup, line string) (sudo, hostgroup, cmd string, err error) {
	ss := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(ss) != 2 {
		return "", "", "", fmt.Errorf("%s", "sudo need arguments")
	}
	sudo = ss[0]

	hs := strings.SplitN(strings.TrimSpace(ss[1]), " ", 2)
	if len(hs) != 2 {
		return "", "", "", fmt.Errorf("%s", "sudo need arguments")
	}
	hg := hgm[hs[0]]
	if hg.Groupname == "" {
		return "", "", "", fmt.Errorf("sudo need hostgroup")
	}
	hostgroup = hs[0]

	cmd = strings.TrimSpace(hs[1])
	if cmd == "" {
		return "", "", "", fmt.Errorf("sudo need command")
	}

	return sudo, hostgroup, cmd, nil
}

func GetDownload(hgm map[string]types.Hostgroup, line string) (download, hostgroup, src, des string, err error) {
	ds := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(ds) != 2 {
		return "", "", "", "", fmt.Errorf("%s", "download need arguments")
	}
	download = ds[0]

	hs := strings.SplitN(strings.TrimSpace(ds[1]), " ", 2)
	if len(hs) != 2 {
		return "", "", "", "", fmt.Errorf("%s", "download need arguments")
	}
	hg := hgm[hs[0]]
	if hg.Groupname == "" {
		return "", "", "", "", fmt.Errorf("download need hostgroup")
	}
	hostgroup = hs[0]

	ss := strings.SplitN(strings.TrimSpace(hs[1]), " ", 2)
	if len(ss) != 2 {
		return "", "", "", "", fmt.Errorf("%s", "download need srdFile")
	}
	src = ss[0]

	des = strings.TrimSpace(ss[1])
	if des == "" {
		return "", "", "", "", fmt.Errorf("download need desDir")
	}

	if len(strings.Split(des, " ")) != 1 {
		return "", "", "", "", fmt.Errorf("download too manay arguments")
	}
	return download, hostgroup, src, des, nil
}

func GetUpload(hgm map[string]types.Hostgroup, line string) (upload, hostgroup, src, des string, err error) {
	us := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(us) != 2 {
		return "", "", "", "", fmt.Errorf("%s", "upload need arguments")
	}
	upload = us[0]

	hs := strings.SplitN(strings.TrimSpace(us[1]), " ", 2)
	if len(hs) != 2 {
		return "", "", "", "", fmt.Errorf("%s", "upload need arguments")
	}
	hg := hgm[hs[0]]
	if hg.Groupname == "" {
		return "", "", "", "", fmt.Errorf("upload need hostgroup")
	}
	hostgroup = hs[0]

	ss := strings.SplitN(strings.TrimSpace(hs[1]), " ", 2)
	if len(ss) != 2 {
		return "", "", "", "", fmt.Errorf("%s", "upload need srdFile")
	}
	src = ss[0]

	des = strings.TrimSpace(ss[1])
	if des == "" {
		return "", "", "", "", fmt.Errorf("upload need desDir")
	}

	if len(strings.Split(des, " ")) != 1 {
		return "", "", "", "", fmt.Errorf("upload too many arguments")
	}
	return upload, hostgroup, src, des, nil
}
