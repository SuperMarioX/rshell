package options

import (
	"flag"
	"fmt"
	"github.com/luckywinds/rshell/pkg/checkers"
	. "github.com/luckywinds/rshell/types"
	"gopkg.in/yaml.v2"
	"html/template"
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
	if cfg.Hostgroupsize == 0 {
		cfg.Hostgroupsize = 200
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

func incIp(s string) string {
	ip := net.ParseIP(s)
	for j := len(ip)-1; j>=0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
	return ip.String()
}

func parseHosts(hg Hostgroup) Hostgroup {
	for _, h := range hg.Hosts {
		if !checkers.IsIpv4(h) {
			log.Fatalf("IP illegal [%s/%s].", hg.Groupname, h)
		}
	}
	hg.Ips = append(hg.Ips, hg.Hosts...)
	return hg
}

func parseHostrange(hg Hostgroup) Hostgroup {
	for _, hr := range hg.Hostranges {
		if !checkers.IsIpv4(hr.From) || !checkers.IsIpv4(hr.To) {
			log.Fatalf("IP Range illegal [%s/%s-%s].", hg.Groupname, hr.From, hr.To)
		}
		temp := []string{hr.From}
		found := false
		nip := hr.From
		count := 0
		for {
			count++
			nip = incIp(nip)
			if nip == hr.To {
				found = true
				temp = append(temp, nip)
				break
			}
			if count > cfg.Hostgroupsize && !found {
				log.Fatalf("Too Large Not Found. IP Range illegal [%s/%s-%s].", hg.Groupname, hr.From, hr.To)
			}
			temp = append(temp, nip)
		}
		if found {
			hg.Ips = append(hg.Ips, temp...)
		} else {
			log.Fatalf("IP Range illegal [%s/%s-%s].", hg.Groupname, hr.From, hr.To)
		}
	}
	return hg
}

func parseHostgroups(hgs Hostgroups) Hostgroups {
	var tmphg = Hostgroups{}
	for _, hg := range hgs.Hgs {
		hg = parseHosts(hg)
		hg = parseHostrange(hg)
		tmphg.Hgs = append(tmphg.Hgs, hg)
	}

	var hgmap = make(map[string]Hostgroup)
	for _, value := range tmphg.Hgs {
		hgmap[value.Groupname] = value
	}
	if len(tmphg.Hgs) != len(hgmap) {
		log.Fatal("There is duplicate hostgroup.")
	}

	var rethg = Hostgroups{}
	for _, hg := range tmphg.Hgs {
		for _, g := range hg.Groups {
			if hgmap[g].Groupname == "" {
				log.Fatalf("Not found. Groups illegal [%s/%s].", hg.Groupname, g)
			}
			hg.Ips = append(hg.Ips, hgmap[g].Ips...)
		}
		if checkers.IsDuplicate(hg.Ips) {
			log.Fatalf("IP Duplicate. Hostgroup illegal [%s].", hg.Groupname)
		}

		if !checkers.CheckHostgroupSize(hg, cfg.Hostgroupsize) {
			log.Fatalf("Too large. IP Range illegal [%s] > [%d].", hg.Groupname, cfg.Hostgroupsize)
		}

		rethg.Hgs = append(rethg.Hgs, hg)
	}

	return rethg
}

func GetHostgroups() (Hostgroups, map[string]Hostgroup) {
	initHostgroups()
	if len(hostgroups.Hgs) == 0 {
		log.Fatal("The hostgroups empty.")
	}
	hostgroups = parseHostgroups(hostgroups)
	var ret = make(map[string]Hostgroup)
	for _, value := range hostgroups.Hgs {
		ret[value.Groupname] = value
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

var tempScript = ".rshell/rshell.tmp"

func templateScript(script string) {
	t, err := template.ParseFiles(script)
	if err != nil {
		log.Fatalf("Parser template script file [%s] failed.", script)
	}

	f, err := os.Create(tempScript)
	if err != nil {
		log.Fatal("Create temp script file failed.")
	}
	defer f.Close()

	if err := t.Execute(f, tasks.Env); err != nil {
		log.Fatal("Parser template script file task failed.")
	}
}

func GetTasks() (Tasks, bool) {
	if *script != "" {
		initTasks(*script)
		if len(tasks.Ts) == 0 {
			log.Fatal("The tasks empty.")
		}
		if len(tasks.Env) != 0 {
			templateScript(*script)
			initTasks(tempScript)
		}
		return tasks, true
	} else {
		return tasks, false
	}
}
