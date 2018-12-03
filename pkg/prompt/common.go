package prompt

import (
	"github.com/luckywinds/rshell/types"
	"github.com/scylladb/go-set/strset"
)

var kset = strset.New()

func AddKeyword(keyword string) {
	kset.Add(keyword)
}

var hset = strset.New()

func AddHostgroup(hostgroup string) {
	hset.Add(hostgroup)
}

var cset = strset.New()

func AddCmd(cmd string) {
	if cmd == "" {
		return
	}
	cset.Add(cmd + cfg.CmdSeparator)
}

var sset = strset.New()

func AddSrcFile(src string) {
	sset.Add(src)
}

var dset = strset.New()

func AddDesDir(des string) {
	dset.Add(des)
}

func init() {
	AddKeyword("do")
	AddKeyword("sudo")
	AddKeyword("upload")
	AddKeyword("download")
}

var commonCmd = []string{
	"date",
	"ls",
	"hostname",
	"whoami",
	"pwd",
}

var cfg types.Cfg
