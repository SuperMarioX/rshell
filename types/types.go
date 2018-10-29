package types

//Default config
type Cfg struct {
	Concurrency  int      `yaml:"concurrency,omitempty"`
	Tasktimeout  int      `yaml:"tasktimeout,omitempty"`
	CmdSeparator string   `yaml:"cmdseparator,omitempty"`
	BlackCmdList []string `yaml:"blackcmdlist,omitempty"`
	PromptString string   `yaml:"promptstring,omitempty"`
}

//Hosts config
type Hostgroup struct {
	Groupname  string   `yaml:"groupname,omitempty"`
	Authmethod string   `yaml:"authmethod,omitempty"`
	Sshport    int      `yaml:"sshport,omitempty"`
	Hosts      []string `yaml:"hosts,omitempty"`
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
	Name    string 	 `yaml:"name,omitempty"`
	Mode    string   `yaml:"mode,omitempty"`
	Sudo    bool     `yaml:"sudo,omitempty"`
	Cmds    []string `yaml:"cmds,omitempty"`
	FtpType string   `yaml:"ftptype,omitempty"`
	SrcFile string   `yaml:"srcfile,omitempty"`
	DesDir  string   `yaml:"desdir,omitempty"`
}

type Task struct {
	Taskname   string `yaml:"taskname,omitempty"`
	Hostgroups string `yaml:"hostgroups,omitempty"`
	Subtasks   []Subtask `yaml:"subtasks,omitempty"`
}

type Tasks struct {
	Ts []Task `yaml:"tasks,omitempty"`
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
