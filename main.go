package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/luckywinds/lwssh"
	. "github.com/luckywinds/rshell/types"
	"github.com/luckywinds/rshell/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Interactive string = "interactive"
	SCRIPT      string = "script"
)

var (
	cfgFile   = flag.String("cfg", "cfg.yaml", "System Config file to read")
	hostsFile = flag.String("hosts", "hosts.yaml", "Hosts Config file to read")
	authFile  = flag.String("auth", "auth.yaml", "Auth Config file to read")
	mode      = flag.String("mode", "interactive", "The Run mode, options[interactive/script], DEFAULT is interactive.")
	script    = flag.String("f", "", "The script yaml.")
	hg        = flag.String("g", "", "The target hostgroup.")
	cmd       = flag.String("c", "", "The commands split by ;")
)

func initFlag() {
	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`Usage:`, cmd, `[<options>] <playbook>

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

var hostgroups Hostgroups

func initHosts() {
	h, err := ioutil.ReadFile(*hostsFile)
	if err != nil {
		log.Fatalf("Can not find hosts file[%s].", *hostsFile)
	}
	err = yaml.Unmarshal(h, &hostgroups)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
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

var tasks Tasks

func initArgs(playbook string) {
	p, err := ioutil.ReadFile(playbook)
	if err != nil {
		log.Fatalf("Can not find script file[%s].", playbook)
	}

	err = yaml.Unmarshal(p, &tasks)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func hgCompleter(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(ps, d.GetWordBeforeCursor(), true)
}

var ps []prompt.Suggest

func main() {
	initFlag()
	initCfg()
	initHosts()
	initAuths()

	if len(hostgroups.Hgs) == 0 {
		fmt.Println("The hostgroup empty.")
	}

	if *mode == Interactive {
		interactiveRun()
	} else if *mode == SCRIPT {
		scriptRun()
	} else {
		log.Fatal("The mode not supported. not in [interactive/script]")
	}
}

func interactiveRun() {
	for _, value := range hostgroups.Hgs {
		p := prompt.Suggest{
			Text:        value.Groupname,
			Description: "",
		}
		ps = append(ps, p)
	}
	fmt.Println("Please select hostgroup and input commmands.")
	for {
		v := prompt.Input(">>> ", hgCompleter)
		if strings.Trim(string(v), " ") == "" {
			continue
		}
		vs := strings.SplitN(string(v), " ", 2)
		if len(vs) != 2 {
			fmt.Println("The hostgroup and commands needed.")
			continue
		}
		if vs[0] == "" || vs[1] == "" {
			fmt.Println("The hostgroup and commands needed.")
			continue
		}
		t := Task{
			Taskname:   "DEFAULT",
			Hostgroups: vs[0],
			Sshtasks:   strings.Split(vs[1], ";"),
		}
		tasks = Tasks{}
		tasks.Ts = append(tasks.Ts, t)
		err := run()
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}

func scriptRun() {
	if *script == "" {
		fmt.Println("The script yaml needed.")
	}
	initArgs(*script)
	err := run()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func run() error {
	if len(tasks.Ts) == 0 {
		return fmt.Errorf("%s", "Tasks empty.")
	}

	limit := make(chan bool, cfg.Concurrency)
	defer close(limit)

	for _, task := range tasks.Ts {
		hg := utils.ChooseHostgroups(hostgroups, task.Hostgroups)
		if hg.Groupname == "" {
			return fmt.Errorf("%s", "The hostgroup not found.")
		}

		auth := utils.ChooseAuthmethod(auths, hg.Authmethod)
		if auth.Name == "" {
			return fmt.Errorf("%s", "The auth method not found.")
		}

		username := auth.Username
		password := auth.Password
		privatekey := auth.Privatekey
		passphrase := auth.Passphrase
		ciphers := auth.Ciphers
		sshport := hg.Sshport
		sudotype := auth.Sudotype
		sudopass := auth.Sudopass

		if len(hg.Hosts) == 0 {
			return fmt.Errorf("%s", "Hosts empty.")
		}
		if sshport <= 0 {
			return fmt.Errorf("%s", "SSH Port not right.")
		}
		if username == "" {
			return fmt.Errorf("%s", "Username empty.")
		}
		if password == "" && privatekey == "" {
			return fmt.Errorf("%s", "Password and PrivateKey empty.")
		}
		if len(task.Sshtasks) == 0 {
			return fmt.Errorf("%s", "SSH Tasks empty.")
		}
		if !strings.HasPrefix(task.Sshtasks[len(task.Sshtasks)-1], "exit") {
			task.Sshtasks = append(task.Sshtasks, "exit")
		}
		if sudotype != "" && sudopass != "" {
			task.Sshtasks = append([]string{sudotype, sudopass}, task.Sshtasks...)
			task.Sshtasks = append(task.Sshtasks, "exit")
		}

		taskchs := make(chan Hostresult, len(hg.Hosts))
		var taskresult Taskresult

		for _, host := range hg.Hosts {
			limit <- true
			go func(host string, sshport int, username, password, privatekey, passphrase, sudotype, sudopass string, cipers []string, task Task) {
				var hostresult Hostresult
				hostresult.Hostaddr = host

				var stdout, stderr string
				var err error
				stdout, stderr, err = lwssh.SSHShell(host, sshport, username, password, privatekey, passphrase, ciphers, task.Sshtasks)
				if err != nil {
					hostresult.Error = err.Error()
				} else {
					hostresult.Stdout = stdout
					hostresult.Stderr = stderr
				}
				taskchs <- hostresult
				<-limit
			}(host, sshport, username, password, privatekey, passphrase, sudotype, sudopass, ciphers, task)
		}

		for i := 0; i < len(hg.Hosts); i++ {
			taskresult.Name = task.Taskname
			select {
			case res := <-taskchs:
				taskresult.Results = append(taskresult.Results, res)
			case <-time.After(time.Duration(cfg.Tasktimeout) * time.Second):
				taskresult.Results = append(taskresult.Results, Hostresult{})
			}
		}
		utils.Output(taskresult)
		close(taskchs)
	}

	return nil
}
