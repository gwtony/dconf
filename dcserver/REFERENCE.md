# API REFERENCE

## Service Management
### Add a service
```
POST /service/add HTTP/1.1
Content-Type: application/json
{
	"service": $service,
	"description": $description,
}
 
HTTP/1.1 200 OK
{
    "token": $token
}
```
> Caution: service SHOULD NOT contains "/"

### Delete a service

```  
POST /service/delete HTTP/1.1
Token: $token
Content-Type: application/json
{
    "service": $service
}

HTTP/1.1 200 OK
```

> Caution: delete will fail if service has any group("all" and "default" exclueded)

## Group Management
### Add a group
```
POST /group/add HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"description": $description
}
 
HTTP/1.1 200 OK
```
> Caution: group SHOULD NOT contains "/"

### Delete a group
```
POST /group/delete HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group
}
 
HTTP/1.1 200 OK
```
> Caution: delete "deault" group is forbidden

> Caution: SHOULD NOT delete non-empty group

### Update a group
Update the group desctiption

```
POST /group/delete HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"description": $description
}
 
HTTP/1.1 200 OK
```

### Read all group of specifed service

```
POST /group/read HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
}

HTTP/1.1 200 OK
{
	"result": [
	{
		"group": $group1,
		"description": $description
	}, ...,
	]
}
```

## Group Member Management
### Add a member
Add a member to specifed service, only add to "all" group 
> Caution: only be called by administrator
```
POST /member/add HTTP/1.1
Content-Type: application/json
{
	"service": $service,
	"ip": $ip,
}

HTTP/1.1 200 OK
```

### Delete a member
Delete a member from specifed service and group
> Caution: only be called by administrator
```
POST /member/delete HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"ip": $ip,
}

HTTP/1.1 200 OK
```

### Move a member
Move a member between groups
```
POST /member/move HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"from": $group1,
	"to": $group2,
	"ip": $ip,
}
 
HTTP/1.1 200 OK
```

### Read member
Read all member from a specifed group, if group is "*", it will return all members in all group of this service(group "all" exclueded)

```
POST /member/read HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
}
 
HTTP/1.1 200 OK
{
	"result": [{
		"group": $group1,
			"ip": [$ip1, ..., $ip2]
		}, ...
	]
}
```

## Config Management
### Add a config
Add a config to a specifed group of a service

```
POST /config/add HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"key": $key,
	"value": $value
}
 
HTTP/1.1 200 OK
```

> Caution: value NOT support binary, base64-encode if needed

> Caution: key SHOULD NOT contains "/"

### Delete a config
```
POST /config/delete HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"key": $key
}
 
HTTP/1.1 200 OK
```

### Update a config
```
POST /config/update HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"key": $key,
	"value": $value
}
 
HTTP/1.1 200 OK
```

### Read a config
If key is "*", it will read all config in group
```
POST /config/read HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"key": $key
}
 
HTTP/1.1 200 OK
{
	"result": [
		{
			"key": $key1,
			"value": $value1
		}, ...
	]
}
```
 
### Copy config
Copy configs between groups in same service
> If key is "*", it will copy all config in source group
```
POST /config/copy HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"from": $group_source,
	"to": $group_dest,
	"key": $key
} 
 
HTTP/1.1 200 OK
```
  
## Render Management
### Do render
Need to set a tag for this render
```
POST /render/do HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"group": $group,
	"key": $key,
	"tag": $tag
}
 
HTTP/1.1 200 OK
{
    "version": $version
}
```

### Read a render
If key is "*", it will read all keys in this condition(key, service)

```
POST /render/read HTTP/1.1
Token: $token
Content-Type: application/json
{
	"service": $service,
	"Ip": $ip,
	"key": $key,
}
 
HTTP/1.1 200 OK
{
	"result": [
	{
		"key": $key,
		"value": $value
	}, ...
	]
}
```

## Status Code 

* 200 - Success
* 204 - No Content (Not found record)
* 400 - Bad request error (Arguments invalid)
* 401 - Unauthrized
* 404 - Page not found (Location incorrect)
* 500 - Internal server error (Api server internal error)
* 502 - Bad Gateway (Backend error)
