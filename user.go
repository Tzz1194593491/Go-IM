package main

import (
	"fmt"
	"net"
	"strings"
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
	_this.server.OnlineMap[_this.Name] = _this
	_this.server.mapLock.Unlock()
	//广播消息
	_this.server.BroadCast(_this, "login")
}

func (_this *User) Offline() {
	//加入到在线用户列表
	_this.server.mapLock.Lock()
	delete(_this.server.OnlineMap, _this.Name)
	_this.server.mapLock.Unlock()
	//广播消息
	_this.server.BroadCast(_this, "logout")
}

func (_this User) SendMsg(msg string) {
	_, write := _this.conn.Write([]byte(msg + "\n"))
	if write != nil {
		fmt.Println("SendMsg error occur", write)
		return
	}
}

func (_this *User) DoMessage(msg string) {
	if msg == "who" {
		_this.who()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		_this.rename(strings.Split(msg, "|")[1])
	} else {
		//广播消息
		_this.server.BroadCast(_this, msg)
	}
}

func (_this *User) ListenMessage() {
	for {
		msg := <-_this.C
		if msg == "" {
			return
		}
		_this.SendMsg(msg)
	}
}

func (_this *User) who() {
	msg := ""
	_this.server.mapLock.RLock()
	for name, _ := range _this.server.OnlineMap {
		msg += name + " online"
	}
	_this.server.mapLock.RUnlock()
	fmt.Println(msg)
	_this.SendMsg(msg)
}

func (_this *User) rename(newName string) {
	user := _this.server.OnlineMap[_this.Name]
	_this.server.mapLock.Lock()
	delete(_this.server.OnlineMap, user.Name)
	_this.server.OnlineMap[newName] = user
	_this.Name = newName
	_this.server.mapLock.Unlock()
	_this.SendMsg("you have used new name")
}
