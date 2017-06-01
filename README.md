# Dconf
A config management service

## Modules
* dcserver: config manage api server, publish config
* dcagent: client agent monitor
* etcd: config storage backend

## Architecture

            publish config(http api)
    	     /
         dcserver
            |
           etcd
            |
         dcagent
            | get config(http api)
          client
    



