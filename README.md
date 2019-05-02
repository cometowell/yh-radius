#GO-RAD系统
go-rad 是一个使用golang语言开发的radius协议实现，支持华为，思科，中心，RouterOS(MikroTik)以及标准协议
## 编译安装

如若想自行编译系统，可以到这个网址下载golang语言安装包: https://golang.google.cn/.

你需要下载golang1.11以上版本，并设置好相应的环境变量(如果想直接使用可以下载dist目录下的相应安装包).

如果你对golang开发或者radius系统打包感兴趣的话，可以自行编译相应平台的安装包。

例如在windows平台下:
```  
    cd source_code_dir
    set CGO_ENABLED=0
    set GOOS=linux
    set GOARCH=amd64 
    go build
```

## 运行系统

你需要复制这些目录或者文件至你的目标目录: go-rad, attributes, config, startup.sh, shutdown.sh

文件目录结构如下:

    |___ attributes
  
    |___ config
  
    |__ go-rad
    
    |__ startup.sh
    
    |__ shutdown.sh

#### 在Linux系统中运行系统: 

> chmod +x startup.sh

> ./startup.sh

#### Linux系统中停止系统:

> chmod +x shutdown.sh

> ./shutdown.sh

## 配置文件解释

| 字段名 | 默认值 | 类型 | 描述 |
| ------| ------ | ------ | ----- |
| 认证端口 | 1812 | int |  radius认证端口  |
| 计费端口 | 1813 | int |  radius计费端口  |
| 密码加密秘钥 | 支持16,24,32长度的十六进制字符串 | string |  用于加密用户密码  |
| 默认的下发会话时长 | 604800 | int | 一周的秒数  |
| limiter.limit | 100 | int | 用于限制每次添加到令牌桶中的token数量，间接控制go协程并发数量 |
| limiter.burst | 1000 | int | 用于限制最多的可用token数量,间接控制go协程并发数量  |

## 数据库表结构
数据库表定义在radius-tables.sql中

## 使用radius-web管理平台
这里有一个可用的radius管理平台，实现了用户管理，套餐管理，nas管理，在线用户管理，管理员管理，角色管理等[RADIUD-WEB](https://github.com/cometowell/radius-web.git)
web平台默认的登陆用户: admin/123456

![首页](https://github.com/cometowell/go-rad/tree/master/document/index.png)
![用户管理](https://github.com/cometowell/go-rad/tree/master/document/user.png)
![用户续费](https://github.com/cometowell/go-rad/tree/master/document/continue.png)
![套餐管理](https://github.com/cometowell/go-rad/tree/master/document/product.png)
![在线用户](https://github.com/cometowell/go-rad/tree/master/document/online.png)
![NAS管理](https://github.com/cometowell/go-rad/tree/master/document/nas.png)


## 许可协议
[MIT](https://mit-license.org/)