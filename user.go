package main

import (
	"fmt"
	"net"
)

type User struct {
	Name    string
	Address string
	C       chan string
	conn    net.Conn
	server  *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddress := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddress,
		Address: userAddress,
		C:       make(chan string),
		conn:    conn,
		server:  server,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

func (_this *User) Online() {
	//加入到在线用户列表
	_this.server.mapLock.Lock()
	_this.server.OnlineMap[_this.Address] = _this
	_this.server.mapLock.Unlock()
	//广播消息
	_this.server.BroadCast(_this, "login")
}

func (_this *User) Offline() {
	//加入到在线用户列表
	_this.server.mapLock.Lock()
	delete(_this.server.OnlineMap, _this.Address)
	_this.server.mapLock.Unlock()
	//广播消息
	_this.server.BroadCast(_this, "login")
}

func (_this *User) DoMessage(msg string) {
	//广播消息
	_this.server.BroadCast(_this, msg)
}

func (_this *User) ListenMessage() {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(_this.conn)

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(_this.conn)

	for {
		msg := <-_this.C
		_, err := _this.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println(_this.Name, "has", "error", err)
			return
		}
	}
}
