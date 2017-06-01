# Dconf Server
Dconf api server to manage config record with etcd

## Config Example
```
[default]
http_addr: 0.0.0.0:10010

log: ../log/dcserver.log
level: debug

[dcserver]
admin_token: AdMiN
etcd_addr: 127.0.0.1:2379
api_location: /dconf
etcd_root: /dconf
```
* http_addr: http listen address
* log: log file
* level: log level
* admin_token: admin token
* etcd_addr: etcd address, example: 1.1.1.1:2379,1.1.1.2:2379,1.1.1.3:2379
* etcd_root: root prefix in etcd, should match with dcagent

## Usage
* -f config file
* -h help
* -v version

## Schema
* [schema infomation](SCHEMA.md)

## Dependency
* [gapi](https://github.com/gwtony/gapi)
* [etcd clientv3](http://github.com/coreos/etcd/clientv3)
* [etcd mvccpb](http://github.com/coreos/etcd/mvcc/mvccpb)
* [etcd v3 rpc type](http://github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes)
