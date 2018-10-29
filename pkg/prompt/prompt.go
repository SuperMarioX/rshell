package prompt

import (
	"github.com/chzyer/readline"
	"github.com/luckywinds/rshell/types"
	"github.com/scylladb/go-set/strset"
)

func listHgs(h types.Hostgroups) func(string) []string {
	return func(string) []string {
		hgs := []string{}
		for _, value := range h.Hgs {
			hgs = append(hgs, value.Groupname)
		}
		return hgs
	}
}

var cset = strset.New()

func listCmds() func(string) []string {
	return func(string) []string {
		return cset.List()
	}
}
func AddCmd(cmd string) {
	cset.Add(cmd + cfg.CmdSeparator)
}

var sset = strset.New()

func listSrcs() func(string) []string {
	return func(string) []string {
		return sset.List()
	}
}
func AddSrcFile(src string) {
	sset.Add(src)
}

var dset = strset.New()

func listDess() func(string) []string {
	return func(string) []string {
		return dset.List()
	}
}
func AddDesDir(des string) {
	dset.Add(des)
}

func newCompleter(hostgroups types.Hostgroups) *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("do",
			readline.PcItemDynamic(listHgs(hostgroups),
				readline.PcItemDynamic(listCmds(),
					readline.PcItemDynamic(listCmds(),
						readline.PcItemDynamic(listCmds(),
							readline.PcItemDynamic(listCmds(),
								readline.PcItemDynamic(listCmds(),
									readline.PcItemDynamic(listCmds()))))))),
		),
		readline.PcItem("sudo",
			readline.PcItemDynamic(listHgs(hostgroups),
				readline.PcItemDynamic(listCmds(),
					readline.PcItemDynamic(listCmds(),
						readline.PcItemDynamic(listCmds(),
							readline.PcItemDynamic(listCmds(),
								readline.PcItemDynamic(listCmds(),
									readline.PcItemDynamic(listCmds()))))))),
		),
		readline.PcItem("download",
			readline.PcItemDynamic(listHgs(hostgroups),
				readline.PcItemDynamic(listSrcs(),
					readline.PcItemDynamic(listDess()))),
		),
		readline.PcItem("upload",
			readline.PcItemDynamic(listHgs(hostgroups),
				readline.PcItemDynamic(listSrcs(),
					readline.PcItemDynamic(listDess()))),
		),
	)
}

var commonCmd = []string{
	"date",
	"free",
	"hostname",
	"whoami",
	"pwd",
}

var cfg types.Cfg

func NewReadline(c types.Cfg, hostgroups types.Hostgroups) (*readline.Instance, error) {
	cfg = c

	for _, value := range commonCmd {
		cset.Add(value + cfg.CmdSeparator)
	}

	return readline.NewEx(&readline.Config{
		Prompt:            cfg.PromptString,
		HistoryFile:       ".rshell/rshell.history",
		AutoComplete:      newCompleter(hostgroups),
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
}
