# 使用说明

## 版本说明

- Linux：rshell
- Windows：rshell.exe

## 命令行交互执行模式

### 启动

```
#rshell
 ______     ______     __  __     ______     __         __
/\  == \   /\  ___\   /\ \_\ \   /\  ___\   /\ \       /\ \
\ \  __<   \ \___  \  \ \  __ \  \ \  __\   \ \ \____  \ \ \____
 \ \_\ \_\  \/\_____\  \ \_\ \_\  \ \_____\  \ \_____\  \ \_____\
  \/_/ /_/   \/_____/   \/_/\/_/   \/_____/   \/_____/   \/_____/
------ Rshell @5.0 Type "?" or "help" for more information. -----
rshell:
```

### 应用场景

#### 获取帮助

```
rshell: ?
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

#### 执行远程Shell命令

- 以普通用户执行单条命令

```
rshell: do test-group01 whoami
TASK [do@test-group01               ] *****************************************
HOST [192.168.31.37   ] =======================================================
ttt

HOST [192.168.31.63   ] =======================================================
ttt

```

- 以普通用户执行多条命令

```
rshell: do test-group01 whoami;date;pwd;echo 123;date;
TASK [do@test-group01               ] *****************************************
HOST [192.168.31.37   ] =======================================================
ttt
Mon Dec  3 20:25:40 CST 2018
/home/ttt
123
Mon Dec  3 20:25:40 CST 2018

HOST [192.168.31.63   ] =======================================================
ttt
Mon Dec  3 20:22:38 CST 2018
/home/ttt
123
Mon Dec  3 20:22:38 CST 2018

```

- 自动切换root用户并执行单条命令

```
rshell: sudo test-group01 whoami
TASK [sudo@test-group01             ] *****************************************
HOST [192.168.31.37   ] =======================================================
root
Last login: Mon Dec  3 20:26:42 CST 2018

STDERR =>
Password:
HOST [192.168.31.63   ] =======================================================
root
Last login: Mon Dec  3 20:23:38 CST 2018

STDERR =>
Password:
```

- 自动切换root用户并执行多条命令

```
rshell: sudo test-group01 whoami;date;pwd;echo 123;date;
TASK [sudo@test-group01             ] *****************************************
HOST [192.168.31.63   ] =======================================================
root
Mon Dec  3 20:24:44 CST 2018
/root
123
Mon Dec  3 20:24:44 CST 2018
Last login: Mon Dec  3 20:24:02 CST 2018

STDERR =>
Password:
HOST [192.168.31.37   ] =======================================================
root
Mon Dec  3 20:27:47 CST 2018
/root
123
Mon Dec  3 20:27:47 CST 2018
Last login: Mon Dec  3 20:27:06 CST 2018

STDERR =>
Password:
```

#### 执行远程上传下载

> 上传下载操作采用SSH通道，不支持自动切换root用户，可以通过执行命令方式实现文件权限及路径调整动作

- 上传文件到用户目录当前目录

> 文件需要已经存在

```
rshell: upload test-group01 LICENSE .
TASK [upload@test-group01           ] *****************************************
HOST [192.168.31.37   ] =======================================================
UPLOAD Success [LICENSE -> ./]
HOST [192.168.31.63   ] =======================================================
UPLOAD Success [LICENSE -> ./]
```

- 上传文件到特定目录

> 文件需要已经存在，目录需要已经存在并且具有用户操作权限

```
rshell: upload test-group01 LICENSE /home/ttt/tst
TASK [upload@test-group01           ] *****************************************
HOST [192.168.31.63   ] =======================================================
UPLOAD Success [LICENSE -> /home/ttt/tst/]
HOST [192.168.31.37   ] =======================================================
UPLOAD Success [LICENSE -> /home/ttt/tst/]
rshell: do test-group01 ls /home/ttt/tst
TASK [do@test-group01               ] *****************************************
HOST [192.168.31.37   ] =======================================================
LICENSE

HOST [192.168.31.63   ] =======================================================
LICENSE

```

- 下载文件到当前目录

> 文件需要已经存在

```
rshell: download test-group01 tst/LICENSE .
TASK [download@test-group01         ] *****************************************
HOST [192.168.31.37   ] =======================================================
DOWNLOAD Success [tst/LICENSE -> test-group01]
HOST [192.168.31.63   ] =======================================================
DOWNLOAD Success [tst/LICENSE -> test-group01]
```

- 下载文件到特定目录

> 文件需要已经存在

```
rshell: download test-group01 tst/LICENSE ./tst
TASK [download@test-group01         ] *****************************************
HOST [192.168.31.37   ] =======================================================
DOWNLOAD Success [tst/LICENSE -> tst/test-group01]
HOST [192.168.31.63   ] =======================================================
DOWNLOAD Success [tst/LICENSE -> tst/test-group01]
```

#### 加解密密码

> 需要在cfg.yaml中配置加密类型和key值

- 加密

```
rshell: encrypt_aes cleartext_password
4ee4b855095dd6607a463b24cfc9438654d4e9b0af4b845ffc9340e70d8d98b2d587
```

- 解密

```
rshell: decrypt_aes 4ee4b855095dd6607a463b24cfc9438654d4e9b0af4b845ffc9340e70d8d98b2d587
cleartext_password
```

#### 退出

```
rshell: exit

```

## 脚本任务编排执行模式

> 如果脚本含有变量，需要提前准备包含变量配置的yaml

- 执行无变量脚本

```
# rshell.exe -f examples\test.yaml
TASK [test all/test do@test-group01                     ] *********************
HOST [192.168.31.37   ] =======================================================
Last failed login: Mon Dec  3 19:05:33 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:15:46 CST 2018
{{.group1.abc}}
{{.group2.abc}}

HOST [192.168.31.63   ] =======================================================
Last failed login: Mon Dec  3 19:02:29 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:12:43 CST 2018
{{.group1.abc}}
{{.group2.abc}}

TASK [test all/test sudo@test-group01                   ] *********************
HOST [192.168.31.37   ] =======================================================
root
Mon Dec  3 21:15:47 CST 2018
Last login: Mon Dec  3 21:13:34 CST 2018

STDERR =>
Password:
HOST [192.168.31.63   ] =======================================================
root
Mon Dec  3 21:12:44 CST 2018
Last login: Mon Dec  3 21:10:31 CST 2018

STDERR =>
Password:
TASK [test all/test upload@test-group01                 ] *********************
HOST [192.168.31.37   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
HOST [192.168.31.63   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
TASK [test all/check upload@test-group01                ] *********************
HOST [192.168.31.37   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:15 test.yaml

HOST [192.168.31.63   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:12 test.yaml

TASK [test all/test download@test-group01               ] *********************
HOST [192.168.31.37   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]
HOST [192.168.31.63   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]

```

- 执行有变量脚本

```
# rshell.exe -f examples\test.yaml -v examples\values.yaml
TASK [test all/test do@test-group01                     ] *********************
HOST [192.168.31.63   ] =======================================================
Last failed login: Mon Dec  3 19:02:29 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:14:06 CST 2018
123
efg

HOST [192.168.31.37   ] =======================================================
Last failed login: Mon Dec  3 19:05:33 CST 2018 from 192.168.1.22 on ssh:notty
There were 2 failed login attempts since the last successful login.
ttt
Mon Dec  3 21:17:09 CST 2018
123
efg

TASK [test all/test sudo@test-group01                   ] *********************
HOST [192.168.31.37   ] =======================================================
root
Mon Dec  3 21:17:10 CST 2018
Last login: Mon Dec  3 21:15:47 CST 2018

STDERR =>
Password:
HOST [192.168.31.63   ] =======================================================
root
Mon Dec  3 21:14:07 CST 2018
Last login: Mon Dec  3 21:12:44 CST 2018

STDERR =>
Password:
TASK [test all/test upload@test-group01                 ] *********************
HOST [192.168.31.37   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
HOST [192.168.31.63   ] =======================================================
UPLOAD Success [examples/test.yaml -> ./]
TASK [test all/check upload@test-group01                ] *********************
HOST [192.168.31.37   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:17 test.yaml

HOST [192.168.31.63   ] =======================================================
-rw-rw-r--. 1 ttt ttt 542 Dec  3 21:14 test.yaml

TASK [test all/test download@test-group01               ] *********************
HOST [192.168.31.37   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]
HOST [192.168.31.63   ] =======================================================
DOWNLOAD Success [test.yaml -> test-group01]

```
