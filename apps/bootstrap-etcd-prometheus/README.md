> 注意Windows和Linux文件路径分隔符的区别

# 命令

```
# rshell.exe -f apps\bootstrap-etcd-prometheus\bootstrap-etcd-prometheus.yaml -v apps\bootstrap-etcd-prometheus\values.yaml
```

# 初始化Grafana

- 登陆：admin/admin/Admin
- 创建：Datasource/etcd

```
Name:   etcd
Type:   Prometheus
Url:    http://localhost:9090
```

- 导入Dashboard：etcd.json
