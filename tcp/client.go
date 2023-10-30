package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var txt = `
tcp 客户端调试使用须知
运行程序：tc.exe ip port msg，如： tc.exe 127.0.0.1 9008 ping
`

func main() {
	parasm := os.Args[1:]
	if len(parasm) == 0 || len(parasm) != 3 {
		fmt.Println("请运行 tc.exe -h 查看说明")
		return
	}
	iptmp := parasm[0]
	if strings.Count(iptmp, ".") != 3 {
		fmt.Println("  ip 格式不对 ，请运行 tc.exe -h 查看说明 ")
		return
	}
	port1, _ := strconv.ParseInt(parasm[1], 10, 64)
	if port1 < 0 || port1 > 65535 {
		fmt.Println("  port 格式不对 ，请运行 tc.exe -h 查看说明 ")
		return
	}
	ip := net.ParseIP(iptmp)
	if ip == nil {
		fmt.Println("  ip 格式不对 ，请运行 tc.exe -h 查看说明 ")
		return
	}
	startTcp(fmt.Sprintf("%s:%d", ip.String(), port1), parasm[2])
}
func startTcp(ipPort, msg string) {
	// 建立与服务器的TCP连接
	conn, err := net.Dial("tcp", ipPort)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()
	for {
		// 发送数据
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Failed to send data to server:", err)
			return
		}
		// 接收服务器的响应
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Failed to receive data from server:", err)
			return
		}
		response := string(buffer[:n])
		fmt.Println("Server response:", response)
		time.Sleep(5 * time.Second)
	}
}
