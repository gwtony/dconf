# Dconf
A config management service

## Modules
* dcserver: config manage api server, config file will be published with http api ([README](dcserver/README.md))
* dcagent: config agent, client can get config from dcagent with http api ([README](dcserver/README.md))
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
    



