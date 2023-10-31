package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var succIp = make(chan string)
var succIpdata = make([][]string, 0)
var errIp = make(chan string)
var wg sync.WaitGroup
var chanlNum = make(chan int, 20)
var txt = `
使用须知
1.运行当前程序需要加 ip 参数，如： scan_ip.exe 127.0.0.*，此处的符合 * 就是需要查找的 ip 
2.假设查找的 ip 端为多个，则参数为：scan_ip.exe 127.*.9.*，此时则会扫描 127.0.9.0->127.255.9.255期间的所有 ip
3.如果需要将扫描的结果保存下来，则可以在后面加上 s 参数，如：scan_ip.exe 127.0.0.1 s，则会在当前目录生成一个扫描成功的 scan_succ_ip.csv 文件
`

func main() {
	//var iptmp = "172.16.20.*"
	//var iptmp = "172.*.20.*"
	//var iptmp = "*.16.20.*"
	parasm := os.Args[1:]

	if len(parasm) == 0 {
		fmt.Println("\n请运行 scan_ip.exe -h 查看说明\n")
		return
	}
	if parasm[0] == "-h" {
		fmt.Println(txt)
		return
	}
	iptmp := parasm[0]
	if strings.Count(iptmp, ".") != 3 {
		fmt.Println(" \n ip 格式不对\n")
		return
	}
	go printSucc()
	go printErr()
	println("starting scan ip")
	for _, ipstr := range getCheckIpList(iptmp) {
		wg.Add(1)
		go addWork(ipstr, &wg)
	}
	wg.Wait()
	if len(parasm) == 2 && parasm[1] == "s" {
		saveCsv(succIpdata)
	}
}
func saveCsv(data [][]string) {

	file, err := os.Create("scan_succ_ip.csv") // 创建要写入的文件
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		err := writer.Write(row)
		if err != nil {
			panic(err)
		}
	}
}
func printSucc() {
	for {
		select {
		case v := <-succIp:
			t1 := time.Now().Format("2006-01-02 15:04:05")
			succIpdata = append(succIpdata, []string{fmt.Sprintf("%s  %s", t1, v)})
			fmt.Printf("\n%s ping ip=%s is ok", t1, v)
		}
	}
}
func printErr() {
	for {
		select {
		case _ = <-errIp:
			//fmt.Printf("\n err ip=%s", v)
		}
	}
}
func addWork(ipstr string, wg *sync.WaitGroup) {

	chanlNum <- 1
	defer func() {
		wg.Done()
		<-chanlNum
	}()
	var find = make(chan int)
	ipAddr := net.ParseIP(ipstr)
	if ipAddr == nil {
		fmt.Printf("\nInvalid IP address err:%s", ipstr)
		return
	}
	var AheadExitAll = make(chan int)
	go func() {
		select {
		case <-AheadExitAll:
			return
		default:
			err := pingIp(ipAddr)
			if err != nil {
				//fmt.Printf("\nping IP err:%s", err.Error())
				return
			}
			find <- 1
			return
		}
	}()
	select {
	case <-find:
		go func() { succIp <- ipstr }()
		return
	case <-time.After(time.Duration(3) * time.Second):
		go func() {
			errIp <- ipstr
			AheadExitAll <- 1
		}()
		return
	}
}

func getCheckIpList(ipTmp string) []string {
	ips := make([]string, 0)
	if strings.Count(ipTmp, "*") == 0 {
		ips = append(ips, ipTmp)
	} else {
		generateIP(&ips, "", strings.Split(ipTmp, "."))
	}
	return ips
}

func generateIP(ips *[]string, currentIP string, parts []string) {
	if len(parts) == 0 {
		*ips = append(*ips, currentIP)
		return
	}
	part := parts[0]
	nextParts := parts[1:]
	if part == "*" {
		for i := 0; i <= 255; i++ {
			cu := currentIP + strconv.Itoa(i) + "."
			if strings.Count(cu, ".") == 4 {
				cu = cu[0 : len(cu)-1]
			}
			generateIP(ips, cu, nextParts)
		}
	} else {
		cu := currentIP + part + "."
		if strings.Count(cu, ".") == 4 {
			cu = cu[0 : len(cu)-1]
		}
		generateIP(ips, cu, nextParts)
	}
}

func pingIp(ipAddress net.IP) error {

	conn, err := net.Dial("ip4:icmp", ipAddress.String())
	if err != nil {
		println(1111)
		return err
	}
	defer conn.Close()
	// 构建 ICMP 报文
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("hello"),
		},
	}
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		println(2222)
		return err
	}
	// 发送 ICMP 报文
	_, err = conn.Write(msgBytes)
	if err != nil {
		println(3333)
		return err
	}
	// 接收 ICMP 回复报文
	reply := make([]byte, 1500)
	_, err = conn.Read(reply)
	if err != nil {
		println(44444)
		return err
	}

	// 解析 ICMP 回复报文
	_, err = icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply)
	if err != nil {
		println(5555)
		return err
	}
	//if replyMsg.Type != ipv4.ICMPTypeEchoReply {
	//	println(6666)
	//	return fmt.Errorf("unexpected ICMP message type %v", replyMsg.Type)
	//}
	//fmt.Printf("9999 ->%s", ipAddress.String())
	return nil
}
