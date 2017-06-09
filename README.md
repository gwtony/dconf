# Dconf
A config management service

## Modules
* dcserver: config manage api server, config file will be published with http api
* dcagent: config agent, client can get config from dcagent with http api
* etcd: config storage backend

## Architecture

           publish config(http api)
             /
         dcserver
            |
           etcd
            |
         dcagent-(optional persistent)
            | get config(http api)
          client
    



