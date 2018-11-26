# rshell

## 功能说明

远程批量执行命令APP

- 简单化，单文件运行，无外部依赖
- 跨平台，运行支持Win和Linux平台
- 双模式，支持文件编排和命令行交互操作
- 双类型，支持ssh命令和ftp上传下载文件
- 双认证，支持密码和key认证
- 自切换，支持自动切换root用户
- 高安全，支持高危命令黑名单，密码支持加密
- 智能化，支持自动提示补全，历史搜索
- 定制化，支持提示符、分隔符、超时等定制
- 模板化，文件编排支持变量自定义

## 应用安装

```
go get github.com/luckywinds/rshell
```

## 应用构建

```
go build rshell.go
```

## 配置说明

默认路径：.rshell

文件说明：

- cfg.yaml：系统配置项
- auth.yaml：认证信息配置
- hosts.yaml：主机列表配置

## 文件执行模式

```
rshell -f examples/test.yaml -v examples/values.yaml
```

注：examples存放示例文件test.yaml，使用说明见文件注释说明

## 交互式命令行执行模式

```
rshell
```

使用说明：
```
Usage: <keywords> <hostgroup> <agruments>

do hostgroup cmd1; cmd2; cmd3
    --- Run cmds on hostgroup use normal user
sudo hostgroup sudo cmd1; cmd2; cmd3
    --- Run cmds on hostgroup use root which auto change from normal user
download hostgroup srcFile desDir
    --- Download srcFile from hostgroup to local desDir
upload hostgroup srcFile desDir
    --- Upload srcFile from local to hostgroup desDir

encrypt_aes cleartext_password
    --- Encrypt cleartext_password with aes 256 cfb
decrypt_aes ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb

exit
    --- Exit rshell
?
    --- Help
```

## 输出说明

```
TASK [task name       ] ********************************************************
HOST [host address    ] --------------------------------------------------------

STDERR =>

SYSERR =>

```

- STDERR：命令标准错误
- SYSERR：系统错误


## 引用

- "github.com/chzyer/readline"
- "github.com/luckywinds/lwssh"
- "gopkg.in/yaml.v2"
- "github.com/fatih/color"
- "github.com/scylladb/go-set/strset"
- "github.com/peterh/liner"

变更说明：
