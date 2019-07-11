#yh-RADIUS系统
yh-radius 是一个使用golang语言开发的radius协议实现，目前已适配华为，思科，中心，RouterOS(MikroTik)以及标准协议
## 编译安装

可以使用github上已经发布的release版本

也可以自行编译相应平台的安装包：
例如在windows平台下:
```  
    cd source_code_dir
    set CGO_ENABLED=0
    set GOOS=linux
    set GOARCH=amd64 
    go build
```

### release版本:

系统采用前后端分离的方式开发： radius server后端 + 管理系统前端

管理系统前端需要运行在web服务器环境(nginx, tomcat等)，radius server后端是编译后的二进制版本可按照下述方式运行

## yh-radius系统介绍

编译完成，复制以下目录或者文件至你的运行目录: yh-radius, attributes, config, startup.sh, shutdown.sh

目录结构如下:
 
    yh-radius
        |___ attributes
  
        |___ config
  
        |__ yh-radius
    
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
| auth.port | 1812 | int |  radius认证端口  |
| acct.port | 1813 | int |  radius计费端口  |
| encrypt.key | 支持16,24,32长度的十六进制字符串 | string |  用于加密用户密码  |
| radius.session.timeout | 604800 | int | 默认一周的秒数  |
| limiter.limit | 100 | int | 用于限制每次添加到令牌桶中的token数量，间接控制go协程并发数量, 服务器环境可根据实际情况调整 |
| limiter.burst | 1000 | int | 用于限制最多的可用token数量,间接控制go协程并发数量,服务器环境可根据实际情况调整  |
| product.stage | debug | string | 控制gin日志，sql显示；可选值：test,debug,release 发布生产环境时请修改此配置为：release  |

## 数据库表结构
数据库表定义在radius-v2.sql中

## 使用radius-web管理平台
这里有一个可用的radius管理平台，实现了用户管理，套餐管理，nas管理，在线用户管理，管理员管理，角色管理等[yh-radius-web](https://github.com/cometowell/radius-web.git)
web平台默认的登陆用户: admin/123456

![首页](https://github.com/cometowell/yh-radius/raw/master/document/index.png)
![用户管理](https://github.com/cometowell/yh-radius/raw/master/document/user.png)
![用户续费](https://github.com/cometowell/yh-radius/raw/master/document/continue.png)
![套餐管理](https://github.com/cometowell/yh-radius/raw/master/document/product.png)
![在线用户管理](https://github.com/cometowell/yh-radius/raw/master/document/online.png)
![NAS管理](https://github.com/cometowell/yh-radius/raw/master/document/nas.png)


## 许可协议
[MIT](https://mit-license.org/)