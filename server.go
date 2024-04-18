package main

import (
	"fmt"
	"net"
	"sync"
)

// Server server对象
type Server struct {
	//Server基本信息
	ip   string
	port int
	//在线用户信息维护
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息管道
	Message chan string
}

// NewServer 构造server对象
func NewServer(ip string, port int) *Server {
	server := &Server{
		ip:        ip,
		port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessage 读取消息队列中的消息，并将其写入用户的信箱
func (_this *Server) ListenMessage() {
	for {
		msg := <-_this.Message
		_this.mapLock.RLock()
		for _, cli := range _this.OnlineMap {
			cli.C <- msg
		}
		_this.mapLock.RUnlock()
	}
}

func (_this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg
	_this.Message <- sendMsg
}

func (_this *Server) Handler(accept net.Conn) {
	user := NewUser(accept, _this)

	//用户上线
	user.Online()

	//监听客户端写
	go func() {
		data := make([]byte, 4096)
		for {
			read, err := user.conn.Read(data)
			if err != nil {
				fmt.Println("occur error:", err)
				return
			}
			if read == 0 {
				user.Offline()
				return
			}
			//提取消息
			_msg := string(data[:read])
			user.DoMessage(_msg)
		}
	}()
	//阻塞handler
	select {}
}

// Start 使用TCP监听指定的网络地址
func (_this *Server) Start() {
	//socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", _this.ip, _this.port))
	fmt.Println("address :", fmt.Sprintf("%s:%d", _this.ip, _this.port))
	fmt.Println("listen port is", _this.port)
	if err != nil {
		fmt.Println("net.Listen error", err)
		return
	}
	//close listen socket
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	//start listen messages
	go _this.ListenMessage()

	//accept
	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept error", err)
			return
		}
		//do handler
		go _this.Handler(accept)
	}
}
