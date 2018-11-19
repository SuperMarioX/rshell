package prompt

import (
	"github.com/chzyer/readline"
	"github.com/luckywinds/rshell/types"
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

func listCmds() func(string) []string {
	return func(string) []string {
		return cset.List()
	}
}

func listSrcs() func(string) []string {
	return func(string) []string {
		return sset.List()
	}
}

func listDess() func(string) []string {
	return func(string) []string {
		return dset.List()
	}
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
