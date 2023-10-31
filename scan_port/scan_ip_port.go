package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var succIpPort = make(chan string)
var succIpPortdata = make([][]string, 0)
var errIpPort = make(chan string)
var wg sync.WaitGroup
var chanlNum = make(chan int, 20)
var txt = `
使用须知
1.端口范围：0-65535
1.假设我们需要扫描 127.0.0.1 的 3000，则执行： scan_port.exe 127.0.0.1 3000，则会去扫描 ip 为 127.0.0.1 的 3000 的端口
2.假设我们需要扫描 127.0.0.1 端口范围是3000-5000，则执行： scan_port.exe 127.0.0.1 3000 5000
3.如果需要将扫描的结果保存下来，则可以在后面加上 s 参数，如：scan_port.exe 127.0.0.1 3000 5000 s，则会在当前目录生成一个扫描成功的 scan_succ_ip_port.csv 文件
`

func main() {

	//var iptmp = "172.16.9.103:8001"
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
		fmt.Println(" \n ip 格式不对 ，请运行 scan_ip.exe -h 查看说明 \n")
		return
	}
	var port1 int64
	var port2 int64
	if len(parasm) >= 2 {
		port1, _ = strconv.ParseInt(parasm[1], 10, 64)
		if port1 < 0 || port1 > 65535 {
			fmt.Println("\n请运行 scan_ip.exe -h 查看说明\n")
			return
		}
	}
	if len(parasm) >= 3 {
		port2, _ = strconv.ParseInt(parasm[2], 10, 64)
		if port2 < 0 || port2 > 65535 || port1 >= port2 {
			fmt.Println("\n请运行 scan_ip.exe -h 查看说明\n")
			return
		}
	}
	go printSucc()
	go printErr()
	println("starting scan ip port")
	var checkList = make([]string, 0)
	if len(parasm) == 2 {
		checkList = append(checkList, fmt.Sprintf("%s:%d", iptmp, port1))
	} else {
		for i := port1; i <= port2; i++ {
			checkList = append(checkList, fmt.Sprintf("%s:%d", iptmp, i))
		}
	}
	for _, ipstr := range checkList {
		//println(ipstr)
		wg.Add(1)
		go addWork(ipstr, &wg)
	}
	wg.Wait()
	if len(parasm) == 4 && parasm[3] == "s" {
		saveCsv(succIpPortdata)
	}
}
func saveCsv(data [][]string) {

	file, err := os.Create("scan_succ_ip_port.csv") // 创建要写入的文件
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		err = writer.Write(row)
		if err != nil {
			panic(err)
		}
	}
}
func printSucc() {
	for {
		select {
		case v := <-succIpPort:
			t1 := time.Now().Format("2006-01-02 15:04:05")
			succIpPortdata = append(succIpPortdata, []string{fmt.Sprintf("%s  %s", t1, v)})
			fmt.Printf("\n%s ping ip=%s is ok", t1, v)
		}
	}
}
func printErr() {
	for {
		select {
		case _ = <-errIpPort:
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

	var AheadExitAll = make(chan int)
	go func() {
		select {
		case <-AheadExitAll:
			return
		default:
			err := pingIpPort(ipstr)
			if err != nil {
				//fmt.Printf("\ntelnet IP port err:%s", err.Error())
				return
			}
			find <- 1
			return
		}
	}()
	select {
	case <-find:
		go func() { succIpPort <- ipstr }()
		return
	case <-time.After(time.Duration(3) * time.Second):
		go func() {
			errIpPort <- ipstr
			AheadExitAll <- 1
		}()
		return
	}
}

func pingIpPort(ipPort string) error {
	// 设置连接超时时间
	conn, err := net.DialTimeout("tcp", ipPort, time.Duration(10)*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
