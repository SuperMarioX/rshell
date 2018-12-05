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

## 使用文档

- [配置说明](docs/config.md)
- [使用说明](docs/usage.md)

## 应用示例

- [Auto Bootstrap A ETCD Cluster](apps/bootstrap-etcd-cluster/README.md)
- [Auto Bootstrap Prometheus Grafana And Monitor ETCD Cluster](apps/bootstrap-etcd-prometheus/README.md)
