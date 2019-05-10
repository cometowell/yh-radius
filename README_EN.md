#GO-RAD
go-rad is a Radius protocol implement written in Golang. It supports NAS devices from Huawei, Cisco, ZTE, MikroTik, etc.

## Installation

How to install golang?  please visit the golang's website: https://golang.google.cn/.

you should download golang1.11 or higher and set related environment variables for golang.

then compile source code according to different platforms

for example:
``` 
 On the windows platform
 for linux: 
    cd source_code_dir
    set CGO_ENABLED=0
    set GOOS=linux
    set GOARCH=amd64 
    go build
```

## run this application

copy there files or dirs to the target dir: go-rad, attributes, config, startup.sh, shutdown.sh

the target directory structure like below:

    |___ attributes
  
    |___ config
  
    |__ go-rad
    
    |__ startup.sh
    
    |__ shutdown.sh

#### run application on linux: 

> chmod +x startup.sh

> ./startup.sh

#### stop application on linux:

> chmod +x shutdown.sh

> ./shutdown.sh

## configuration file
config.json file in config directory, you can modify the config item.

| name | default | type | desc |
| ------| ------ | ------ | ----- |
| auth.port | 1812 | int |  authenticate port  |
| acct.port | 1813 | int |  accounting port  |
| encrypt.key | 16/24/32 length hex string | string |  used to encrypt passwords  |
| radius.session.timeout | 604800 | int | session duration, default: sec of a week  |
| limiter.limit | 100 | int | to limit the amount of goroutine |
| limiter.burst | 1000 | int | to limit the amount of goroutine  |

## database tables
the database table structure is defined in the radius-tables.sql

## use go-rad
here is a simple web management system to use[RADIUD-WEB](https://github.com/cometowell/radius-web.git)

default account for login: admin/123456

## License
[MIT](https://mit-license.org/)