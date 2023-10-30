package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"time"
)

func main() {
	parasm := os.Args[1:]
	if len(parasm) != 2 {
		fmt.Println("参数不正确，执行示例：wsc.exe ws://127.0.0.1:9008/api ping ")
		return
	}
	fmt.Println("正在连接：" + parasm[0])
	// 建立 WebSocket 连接
	conn, _, err := websocket.DefaultDialer.Dial(parasm[0], nil)
	if err != nil {
		fmt.Println("无法建立 WebSocket 连接：", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("websocket连接成功~~")
	// 接收消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("读取消息出错：" + err.Error())
				return
			}
			fmt.Printf("\n%s 接收消息：%s ", time.Now().Format("2006-01-02 15:04:05"), string(message))
		}
	}()
	for {
		time.Sleep(5 * time.Second)
		// 发送消息
		message := []byte(parasm[1])
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("发送消息出错：" + err.Error())
			return
		} else {
			fmt.Printf("\n%s 发送消息成功：%s ", time.Now().Format("2006-01-02 15:04:05"), string(message))
		}
	}

}
