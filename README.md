# rshell

## 功能说明

远程批量执行命令APP

- 支持文件执行模式和交互式命令行执行模式
- 支持ssh批量执行shell命令
- 支持sftp批量执行上传下载单文件
- 支持密码和key认证
- 支持自动切换root
- 支持命令关键字和历史自动提示
- 支持彩色输出
- 运行平台支持windows和linux

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
rshell.exe -f tasks.yaml
```

注：examples存放示例文件tasks.yaml

## 交互式命令行执行模式

```
rshell.exe
```

使用说明：
```
Usage:
    SSH Task: hostgroup cmd1; cmd2; cmd3
    SSH Task(sudo root): hostgroup sudo cmd1; cmd2; cmd3
    SFTP Task: hostgroup download/upload srcFile desDir
    Exit: Ctrl c
    Help: ?
```

## 输出说明

```
TASK [task name       ] ********************************************************
HOST [host address    ] --------------------------------------------------------
STDOUT =>

STDERR =>

ERROR =>

```

- STDOUT：命令标准输出
- STDERR：命令标准错误
- ERROR：系统错误

## 约束

- 同一task下的sshtasks和sftptasks无关联关系
- 同一sftptasks下的download和upload无关联关系

## 引用

- "github.com/c-bata/go-prompt"
- "github.com/deckarep/golang-set"
- "gopkg.in/yaml.v2"
- "github.com/fatih/color"