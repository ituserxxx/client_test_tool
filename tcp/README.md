
## tcp 客户端调试工具

### 说明
连接 tcp 服务命令参数：
```shell
tx.exe ip port msg
```
如连接 tcp 服务端 127.0.0.1:9008，并且发送消息 ping，则执行：
```shell
tx.exe 127.0.0.1 9008 ping
```
每5秒向服务端发送消息 ping，如果收到服务端消息则会输出

#### 下载
linux 端下载 ：tc_linux

windows 端下载：tc.exe

#### 源码

windows 打包

```shell
go build -o tc.exe client.go
```

mac 打包
```shell
set GOOS=darwin
set GOARCH=amd64
go build -o tc_mac client.go
```

linux 打包
```shell
set GOARCH=amd64
set GOOS=linux
go build -o tc_linux client.go
```
