package net

import (
	"os"
	"fmt"
	"net"
	"sync"
	"strconv"
	"io"
	"log"
)

const (
	NETWORK string = "tcp"
	LADDR   string = ":8080"
	LOGFILE string = "log.txt"
)

func main(){
	var mu sync.Mutex
	var  count int = 0
	f,err := os.OpenFile(LOGFILE,os.O_RDWR | os.O_APPEND|os.O_CREATE,0666)
	if err != nil{
		fmt.Printf("打开日志文件失败!error:%s\n",err)
		os.Exit(1)
	}
	defer f.Close()
	logger := log.New(f,"",1)
	listener,err := net.Listen(NETWORK,LADDR)
	if err != nil {
		logger.Printf("监听端口失败!error:%s",err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn,err := listener.Accept()

		if err != nil {
			logger.Printf("创建连接失败!error:err%s",err)
			os.Exit(1)
		}
		go connHandle(conn,logger,&count,mu)
	}

}

func connHandle(conn net.Conn,logger *log.Logger,countPtr *int,mu sync.Mutex){
	mu.Lock()
	defer mu.Unlock()
	defer conn.Close()
	*countPtr++
	fileName := "记录"+strconv.Itoa(*countPtr)
	f,err := os.OpenFile(fileName,os.O_RDWR | os.O_APPEND|os.O_CREATE,0666)
	if err != nil{
		logger.Printf("打开记录文件失败!error:%s\n",err)
		return
	}
	defer f.Close()
	Raddr := conn.RemoteAddr().String()
	f.WriteString("客户端:"+Raddr+"已连接")
	var buf []byte = make([]byte,4096)
	for{
		n,err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				f.WriteString("socket连接已关闭!\n")
				break
			}else {
				f.WriteString("写入数据失败!error:"+err.Error()+"\n")
				break
			}
		}
		f.WriteString(string(buf[:n])+"\n")
	}
}