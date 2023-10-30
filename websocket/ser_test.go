package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

func TestServ(t *testing.T) {
	http.HandleFunc("/", handleWebSocket)
	log.Println("WebSocket server started")
	http.ListenAndServe(":9008", nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected:", conn.RemoteAddr())

	for {
		// 读取客户端发送的消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from client:", err)
			break
		}

		log.Println("Received message from client:", string(message))

		// 发送响应给客户端
		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, client!,I am websocket server"))
		if err != nil {
			log.Println("Failed to send response to client:", err)
			break
		}
	}
}
