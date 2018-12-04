# 配置说明

## 配置文件路径

> 当前目录下.rshell文件夹

## 系统配置

> .rshell/cfg.yaml

```
concurrency: 6                                                       #The num of concurrency goroutine for hosts, Default: 6
tasktimeout: 300                                                     #The total timeout [second] of everyone goroutine, Default: 300
cmdseparator: ";"                                                    #The command separator, Default: ";"
promptstring: "rshell: "                                             #The system prompt string, Default: "rshell: "
outputintime: true                                                   #The output in time, Default: false
hostgroupsize: 200                                                   #The max number of ip in hostgroup, Default: 200
passcrypttype: "aes"                                                 #The password crypt type, support [aes], Default: ""
passcryptkey: "I`a!M@a#H$a%P^p&Y*c(R)y_P+t,K.eY"                     #The aes crypt key, length must be 32. Not be empty if passcrypttype = aes.
outputtype: "text"                                                   #The result output format, Default: "text"
blackcmdlist:                                                        #The dangerous black command list, Default if =cmd or ~=cmdprefix can not run.
- cmd: rm -rf /                                                      
  cmdprefix: rm -rf /                                                
mostusedcmds:                                                        #The most used command list
- pwd                                                                
- date                                                               
- whoami                                                             
updateserver:                                                        #The auto update server address list
- https://github.com/luckywinds/rshell/raw/master/releases
```

## 认证信息配置

> .rshell/auth.yaml

```
authmethods:                                                         #The authmethod list
- name: alpha-env-ttt-key-su                                         #The auth method name
  username: ttt                                                      #The username of ssh login
  privatekey: ""                                                     #First auth method option[privatekey], default: ""
  passphrase: "aes"                                                  #First auth method option[passphrase], default: ""
  password: 56efc9c43425221af1ab2ce8d6d52a6d2926ff608d56be0601d6     #Second auth method option[password], default: ""
  sudotype: "su -"                                                   #Optional：Change root command if needed, default: "su"
  sudopass: b6c2e5ba01a5ba9b1b1eae53f839c4bf13d541da17829d8408       #Optional：root password if needed, default: ""

```

> privatekey/password至少选1项即可

## 主机信息配置

> .rshell/hosts.yaml

```
hostgroups:                                 #The hostgroup list
- groupname: test-group03                   #The hostgroup name    
  authmethod: alpha-env-root-pass           #The auth method name which must in authmethods list
  sshport: 22                               #The ssh connect port, default: ""
  groups:                                   #The group list, item must in hostgroups list
  - hostgroup02
  - hostgroup03
  hosts:                                    #The host list, ip format must be ipv4
  - 192.168.31.78
  hostranges:                               #The host ip range list
  - from: 192.168.31.90
    to: 192.168.31.92
```

> groups/hosts/hostranges至少选1项即可

## 脚本配置

> 示例：example\test.yaml

```
tasks:                                      #The task list
- name: test all                            #The task name
  hostgroup: test-group01                   #The hostgroup name
  subtasks:                                 #The sub task list
  - name: test do                           #The sub task name
    mode: ssh                               #The sub task mode ssh/sftp
    cmds:                                   #The command list
    - whoami
    - date
    - echo {{.group1.abc}}
    - echo {{.group2.abc}}

  - name: test sudo                                                
    mode: ssh
    sudo: true
    cmds:
    - whoami
    - date

  - name: test upload
    mode: sftp
    ftptype: upload                         #The ftp type upload/download
    srcfile: examples/test.yaml             #The source file which will be ftp from
    desdir: .                               #The destination dir which wile be ftp to

  - name: check upload
    mode: ssh
    cmds:
    - ls -alh test.yaml

  - name: test download
    mode: sftp
    ftptype: download
    srcfile: test.yaml
    desdir: . 

```

## 变量配置

> 示例：example\values.yaml

```
group1:                                     #The var group name                          
  abc: 123                                  #The var item name
group2:
  abc: efg
```

## 执行逻辑

1. 获取系统配置
2. 获取主机配置
3. 关联认证配置

命令行交互执行模式
1. 执行相关命令

脚本编排执行模式
1. 获取变量配置，如果需要
2. 解析脚本配置，关联变量，如果需要
3. 执行相关任务命令
