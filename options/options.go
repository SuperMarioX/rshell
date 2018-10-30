package options

import (
	"flag"
	"fmt"
	. "github.com/luckywinds/rshell/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
)

var (
	cfgFile   = flag.String("cfg", path.Join(".rshell", "cfg.yaml"), "System Config file to read, Default: "+path.Join(".rshell", "cfg.yaml"))
	hostsFile = flag.String("hosts", path.Join(".rshell", "hosts.yaml"), "Hosts Config file to read, Default: "+path.Join(".rshell", "hosts.yaml"))
	authFile  = flag.String("auth", path.Join(".rshell", "auth.yaml"), "Auth Config file to read, Default: "+path.Join(".rshell", "hosts.yaml"))
	script    = flag.String("f", "", "The script yaml.")
)

func init() {
	os.Mkdir(".rshell", os.ModeDir)
	initFlag()
}

func initFlag() {
	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`Usage:`, cmd, `[<options>]

Options:`)
		flag.PrintDefaults()
	}
	flag.Parse()
}

var cfg Cfg

func initCfg() {
	c, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatalf("Can not find cfg file[%s].", *cfgFile)
	}
	err = yaml.Unmarshal(c, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
func GetCfg() Cfg {
	initCfg()
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 6
	}
	if cfg.Tasktimeout == 0 {
		cfg.Tasktimeout = 300
	}
	if cfg.CmdSeparator == "" {
		cfg.CmdSeparator = ";"
	}
	if cfg.PromptString == "" {
		cfg.PromptString = "rshell: "
	}
	return cfg
}

var hostgroups Hostgroups

func initHostgroups() {
	h, err := ioutil.ReadFile(*hostsFile)
	if err != nil {
		log.Fatalf("Can not find hosts file[%s].", *hostsFile)
	}
	err = yaml.Unmarshal(h, &hostgroups)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, hg := range hostgroups.Hgs {
		for _, h := range hg.Hosts {
			if net.ParseIP(h) == nil {
				log.Fatalf("IP illegal [%s/%s].", hg.Groupname, h)
			}
		}
	}
}
func GetHostgroups() (Hostgroups, map[string]Hostgroup) {
	initHostgroups()
	if len(hostgroups.Hgs) == 0 {
		log.Fatal("The hostgroups empty.")
	}
	var ret = make(map[string]Hostgroup)
	for _, value := range hostgroups.Hgs {
		ret[value.Groupname] = value
	}
	if len(hostgroups.Hgs) != len(ret) {
		log.Fatal("There is duplicate hostgroup.")
	}
	return hostgroups, ret
}

var auths Auths

func initAuths() {
	a, err := ioutil.ReadFile(*authFile)
	if err != nil {
		log.Fatalf("Can not find auth file[%s].", *authFile)
	}
	err = yaml.Unmarshal(a, &auths)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
func GetAuths() (Auths, map[string]Auth) {
	initAuths()
	if len(auths.As) == 0 {
		log.Fatal("The auths empty.")
	}
	var ret = make(map[string]Auth)
	for _, value := range auths.As {
		ret[value.Name] = value
	}
	if len(auths.As) != len(ret) {
		log.Fatal("There is duplicate auth.")
	}
	return auths, ret
}

var tasks Tasks

func initTasks(scriptFile string) {
	p, err := ioutil.ReadFile(scriptFile)
	if err != nil {
		log.Fatalf("Can not find script file[%s].", scriptFile)
	}

	err = yaml.Unmarshal(p, &tasks)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
func GetTasks() (Tasks, bool) {
	if *script != "" {
		initTasks(*script)
		if len(tasks.Ts) == 0 {
			log.Fatal("The tasks empty.")
		}
		return tasks, true
	} else {
		return tasks, false
	}
}
