package prompt

import (
	"github.com/luckywinds/rshell/types"
	"github.com/peterh/liner"
	"log"
	"os"
	"sort"
	"strings"
)

func keywordCompleter(line string) (c []string) {
	l := kset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func hostgroupCompleter(line string) (c []string) {
	l := hset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func cmdCompleter(line string) (c []string) {
	l := cset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func srcCompleter(line string) (c []string) {
	l := sset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func desCompleter(line string) (c []string) {
	l := dset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func wordCompleter(line string, pos int) (head string, completions []string, tail string) {
	head = string([]rune(line)[:pos])
	tail = string([]rune(line)[pos:])

	keyword := strings.SplitN(strings.TrimLeft(line, " "), " ", 2)
	if len(keyword) == 1 {
		return "", keywordCompleter(keyword[0]), tail
	} else {
		hostgroup := strings.SplitN(strings.TrimLeft(keyword[1], " "), " ", 2)
		if len(hostgroup) == 1 {
			return keyword[0] + " ", hostgroupCompleter(hostgroup[0]), tail
		} else {
			if keyword[0] == "do" || keyword[0] == "sudo" {
				cmds := strings.Split(strings.TrimLeft(hostgroup[1], " "), cfg.CmdSeparator)
				if len(cmds) == 0 {
					return keyword[0] + " " + hostgroup[0] + " ", cmdCompleter(""), tail
				} else if len(cmds) == 1 {
					return keyword[0] + " " + hostgroup[0] + " ", cmdCompleter(cmds[len(cmds)-1]), tail
				} else {
					return keyword[0] + " " + hostgroup[0] + " " + strings.Join(cmds[0:len(cmds)-1], cfg.CmdSeparator) + cfg.CmdSeparator, cmdCompleter(cmds[len(cmds)-1]), tail
				}
			} else if keyword[0] == "upload" || keyword[0] == "download" {
				srcfile := strings.SplitN(strings.TrimLeft(hostgroup[1], " "), " ", 2)
				if len(srcfile) == 1 {
					return keyword[0] + " " + hostgroup[0] + " ", srcCompleter(srcfile[0]), tail
				} else {
					desfile := strings.SplitN(strings.TrimLeft(srcfile[1], " "), " ", 2)
					if len(desfile) == 1 {
						return keyword[0] + " " + hostgroup[0] + " " + srcfile[0] + " ", desCompleter(desfile[0]), tail
					}
				}
			}
		}
	}
	return "", []string{}, tail
}

func New(c types.Cfg, hostgroups types.Hostgroups) (*liner.State, error) {
	cfg = c

	for _, value := range hostgroups.Hgs {
		hset.Add(value.Groupname)
	}
	for _, value := range commonCmd {
		cset.Add(value + c.CmdSeparator)
	}

	line := liner.NewLiner()
	line.SetWordCompleter(wordCompleter)
	line.SetTabCompletionStyle(liner.TabPrints)

	if f, err := os.Open(c.HistoryFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	return line, nil
}

func Prompt(l *liner.State, c types.Cfg) (string, error) {
	name, err := l.Prompt(c.PromptString)
	if err == nil {
		l.AppendHistory(name)
	} else if err == liner.ErrPromptAborted {
		return name, ErrPromptAborted
	} else {
		return "", ErrPromptError
	}

	if f, err := os.Create(c.HistoryFile); err != nil {
		log.Print("Error writing history file: ", err)
		return name, ErrPromptError
	} else {
		l.WriteHistory(f)
		f.Close()
	}

	return name, nil
}