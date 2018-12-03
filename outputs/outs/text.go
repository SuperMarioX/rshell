package outs

import (
	"fmt"
	"github.com/fatih/color"
)

type TEXT struct {
	Actionname string
	Actiontype string
	Groupname string
	Address string
	Stdout string
	Stderr string
	Syserr string
}

func (t TEXT) Print() {
	fmt.Printf("%+v", t)
}

func (t TEXT) PrintSimple() {
	color.Green("HOST [%-16s] =======================================================\n", t.Address)
	if t.Stdout != "" {
		fmt.Printf("%s\n", t.Stdout)
	}
	if t.Stderr != "" {
		color.Red("%s\n", "STDERR =>")
		fmt.Printf("%s\n", t.Stderr)
	}
	if t.Syserr != "" {
		color.Red("%s\n", "SYSERR =>")
		fmt.Printf("%s\n", t.Syserr)
	}
	if t.Stdout == "" && t.Stderr == "" && t.Syserr == "" {
		fmt.Println()
	}
}