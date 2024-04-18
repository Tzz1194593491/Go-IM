package main

import (
	"fmt"
	"net"
)

// Server server对象
type Server struct {
	ip   string
	port int
}

// NewServer 构造server对象
func NewServer(ip string, port int) *Server {
	server := &Server{ip: ip, port: port}
	return server
}

func (_this *Server) Handler(accept net.Conn) {
	fmt.Println(accept, "链接建立成功")
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

	//accept
	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept error", err)
			return
		}
		//do handler
		_this.Handler(accept)
	}
}
