package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

var txt = `
mqtt 客户端调试使用须知
运行程序：mqc.exe ip:port clientId msg subTopic pushTopic user pass
如： mqc.exe  "test.mosquitto.org:1883" xx1 "ping form xx1"  xx1111 xx2222
连接成功后会每隔 5s 发送一次消息 msg
`

func main() {
	parasm := os.Args[1:]
	if len(parasm) == 1 && parasm[0] == "-h" {
		fmt.Println(txt)
		return
	}
	if len(parasm) != 5 && len(parasm) != 7 {
		fmt.Println("参数格式不对，请执行：mqc.exe -h 查看详细说明")
		return
	}
	// client.go ip:port clientId msg subTopic pushTopic user pass
	//start("test.mosquitto.org:1883", "xx1", "ping form xx1", "xx1111", "xx2222", "", "")
	start(parasm[0], parasm[1], parasm[2], parasm[3], parasm[4], parasm[5], parasm[6])

}
func start(ipPort, clientId, msg, subTopic, pushTopic, user, pass string) {
	fmt.Printf("starting conecnt  %s", ipPort)
	// 创建MQTT客户端实例
	opts := MQTT.NewClientOptions()
	opts = opts.AddBroker("tcp://" + ipPort)
	opts.SetClientID(fmt.Sprintf("%s-%d", clientId, time.Now().Unix()))
	opts.SetUsername(user)
	opts.SetPassword(pass)
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)
	fmt.Printf("conecnt succ %s", ipPort)
	// 定义消息处理函数
	messageHandler := func(client MQTT.Client, msg MQTT.Message) {

		fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
	}

	// 订阅主题
	if token := client.Subscribe(subTopic, 0, messageHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic %s", subTopic)

	for {
		if client.IsConnected() {
			client.Publish(pushTopic, 0, false, msg)
		}
		time.Sleep(5 * time.Second)
	}
}
