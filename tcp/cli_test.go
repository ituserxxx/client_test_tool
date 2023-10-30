package main

import (
	"fmt"
	"net"
	"testing"
)

func TestTcpServer(t *testing.T) {
	ipPort := "127.0.0.1:9008"
	startTcpServer(ipPort)
}
func startTcpServer(ipPort string) {
	// 监听指定端口
	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started, listening on port 8080")

	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		// 启动一个独立的goroutine处理客户端请求
		go handleClient(conn)
	}
}

// 处理客户端请求
func handleClient(conn net.Conn) {
	defer conn.Close()

	for {
		// 接收客户端发送的消息
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read data from client:", err)
			return
		}
		message := string(buffer[:n])
		fmt.Println("Received message from client:", message)
		// 发送响应给客户端
		response := "Hello, client, I am server!"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Failed to send response to client:", err)
			return
		}
	}
}
