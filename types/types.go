package types

//Default config

type BlackCmd struct {
	Cmd       string `yaml:"cmd,omitempty"`
	CmdPrefix string `yaml:"cmdprefix,omitempty"`
}

type Cfg struct {
	Concurrency   int        `yaml:"concurrency,omitempty"`
	Tasktimeout   int        `yaml:"tasktimeout,omitempty"`
	CmdSeparator  string     `yaml:"cmdseparator,omitempty"`
	BlackCmdList  []BlackCmd `yaml:"blackcmdlist,omitempty"`
	PromptString  string     `yaml:"promptstring,omitempty"`
	Outputintime  bool       `yaml:"outputintime,omitempty"`
	Hostgroupsize int        `yaml:"hostgroupsize,omitempty"`
	Passcrypttype string     `yaml:"passcrypttype,omitempty"`
	Passcryptkey  string     `yaml:"passcryptkey,omitempty"`
}

//Hosts config
type Hostrange struct {
	From string `yaml:"from,omitempty"`
	To   string `yaml:"to,omitempty"`
}

type Hostgroup struct {
	Groupname  string      `yaml:"groupname,omitempty"`
	Authmethod string      `yaml:"authmethod,omitempty"`
	Sshport    int         `yaml:"sshport,omitempty"`
	Hosts      []string    `yaml:"hosts,omitempty"`
	Groups     []string    `yaml:"groups,omitempty"`
	Hostranges []Hostrange `yaml:"hostranges,omitempty"`
	Ips        []string
}
type Hostgroups struct {
	Hgs []Hostgroup `yaml:"hostgroups,omitempty"`
}

//Auth config
type Auth struct {
	Name       string   `yaml:"name,omitempty"`
	Username   string   `yaml:"username,omitempty"`
	Password   string   `yaml:"password,omitempty"`
	Privatekey string   `yaml:"privatekey,omitempty"`
	Passphrase string   `yaml:"passphrase,omitempty"`
	Ciphers    []string `yaml:"ciphers,omitempty"`
	Sudotype   string   `yaml:"sudotype,omitempty"`
	Sudopass   string   `yaml:"sudopass,omitempty"`
}
type Auths struct {
	As []Auth `yaml:"authmethods,omitempty"`
}

//Task config
type Subtask struct {
	Name    string   `yaml:"name,omitempty"`
	Mode    string   `yaml:"mode,omitempty"`
	Sudo    bool     `yaml:"sudo,omitempty"`
	Cmds    []string `yaml:"cmds,omitempty"`
	FtpType string   `yaml:"ftptype,omitempty"`
	SrcFile string   `yaml:"srcfile,omitempty"`
	DesDir  string   `yaml:"desdir,omitempty"`
}

type Task struct {
	Name      string    `yaml:"name,omitempty"`
	Hostgroup string    `yaml:"hostgroup,omitempty"`
	Subtasks  []Subtask `yaml:"subtasks,omitempty"`
}

type Tasks struct {
	Env map[string]interface{} `yaml:"env,omitempty"`
	Ts  []Task                 `yaml:"tasks,omitempty"`
}

//Result
type Hostresult struct {
	Hostaddr string `yaml:"hostaddr,omitempty"`
	Error    string `yaml:"error,omitempyt"`
	Stdout   string `yaml:"stdout,omitempty"`
	Stderr   string `yaml:"stderr,omitempty"`
}
type Taskresult struct {
	Name    string       `yaml:"name,omitempty"`
	Results []Hostresult `yaml:"results,omitempty"`
}
type Tasksresults struct {
	Results []Taskresult `yaml:"results,omitempty"`
}
