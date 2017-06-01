# Dconf Agent

Dconf agent will watch config change from etcd

## Config
Some config in file:
* etcd_root: etcd root path, this should be matched with dcserver
*api_location: agent http api location, default is '/dconf'

# Build
```
make
cd dist/bin
./dcagent
```

# Read config
```
Request
POST /config/get HTTP/1.1
Content-Type: application/json
{
	"service": $service,
	"key": $key 
}

Response
{
	"result": [
		{
			"key": $key,
			"value": $value
		}
	]
}
```
service: service name
key: config name
value: config value

