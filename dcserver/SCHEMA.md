# Table in Etcd

* service index prefix： /srv_view
```
key: /srv_view/{service}/{group}/{key}
value: {value}
```

* group index prefix： /group_view
```
key：/group_view/{service}/{group}/{ip}
value: "alive"
```

* host index prefix：/host_view
```
key：/host_view/{ip}/{srv}/{key}
value: {value}
```

* service meta index
```
key: /srv_meta/{service}
value: json({service, description, token})
```

* group meta index
```
key: /group_meta/{service}/{group}
value: json({service, group, description})
```
