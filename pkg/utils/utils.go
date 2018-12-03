package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/types"
	"strings"
)

func OutputTaskHeader(name string) {
	color.Yellow("TASK [%-30s] *****************************************\n", name)
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

func GetSSHArgs(line string) (string, string, []string, error) {
	ks := strings.Fields(line)
	if len(ks) < 3 {
		return "", "", []string{}, fmt.Errorf("ssh arguments not enough")
	}
	cfg := options.GetCfg()
	if cfg.CmdSeparator == " " {
		return ks[0], ks[1], ks[2:], nil
	} else {
		l := strings.Join(ks[2:], " ")
		return ks[0], ks[1], strings.Split(l, cfg.CmdSeparator), nil
	}
}

func GetSFTPArgs(line string) (string, string, string, string, error) {
	ks := strings.Fields(line)
	if len(ks) != 4 {
		return "", "", "", "", fmt.Errorf("sftp arguments not enough")
	}
	return ks[0], ks[1], ks[2], ks[3], nil
}

func GetCryptArgs(line string) (string, string, error) {
	ks := strings.Fields(line)
	if len(ks) != 2 {
		return "", "", fmt.Errorf("crypt arguments illegal")
	}
	return ks[0], ks[1], nil
}
