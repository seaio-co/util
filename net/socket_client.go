package net

import (
	"fmt"
	"os"
	"net"
	"time"
)


const(
	RADDR = "127.0.0.1:8080"  //服务器地址
)

func main(){
	conn,err := net.DialTimeout(NETWORK, RADDR,5*time.Second)  //创建套接字,连接服务器,设置超时时间

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	for i := 0; i< 100;i++{
		conn.Write([]byte("hello socket\n"))     //发送数据给服务器端
		time.Sleep(1*time.Second)
	}
}