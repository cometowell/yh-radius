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
| authPort | 1812 | int |  authenticate port  |
| acctPort | 1813 | int |  accounting port  |
| encrypt.key | 16/24/32 length hex string | string |  used to encrypt passwords  |
| radius.session.timeout | 604800 | int | session duration, default: sec of a week  |
| limiter.limit | 100 | int | to limit the amount of goroutine |
| limiter.burst | 1000 | int | to limit the amount of goroutine  |

## database tables
the database table structure is defined in the database.sql


## third-party dependencies
``` 
github.com/gin-contrib/sse v0.0.0-20190125020943-a7658810eb74 // indirect
github.com/gin-gonic/gin v1.3.0
github.com/go-sql-driver/mysql v1.4.1
github.com/go-xorm/xorm v0.7.1
github.com/golang/protobuf v1.2.0 // indirect
github.com/gomodule/redigo v2.0.0+incompatible
github.com/mattn/go-isatty v0.0.4 // indirect
github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
github.com/sirupsen/logrus v1.3.0
github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43 // indirect
golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
gopkg.in/yaml.v2 v2.2.2 // indirect
```

## License
[MIT](http://www.baidu.com)