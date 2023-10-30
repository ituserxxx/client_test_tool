# client_test_tool

## mqtt 客户端调试工具

### 说明
连接 mqtt 服务命令参数：
```shell
mqc.exe  ip:port clientId msg subTopic pushTopic user pass
```
如连接 mqtt 服务端 test.mosquitto.org:1883，并且发送消息 ping form xx1，则执行：
```shell
# 不带账号密码
mqc.exe "test.mosquitto.org:1883" xx1 "ping form xx1"  xx11
#带账号密码
mqc.exe "test.mosquitto.org:1883" xx1 "ping form xx1"  xx11 user passwd

```
每5秒向服务端发送消息 "ping form xx1"，如果收到消息则会输出

#### 下载
linux 端下载 ：mqc_linux

windows 端下载：mqc.exe

#### 源码

windows 打包

```shell
go build -o mqc.exe client.go
```

mac 打包
```shell
set GOOS=darwin
set GOARCH=amd64
go build -o mqc_mac client.go
```

linux 打包
```shell
set GOARCH=amd64
set GOOS=linux
go build -o mqc_linux client.go
```
