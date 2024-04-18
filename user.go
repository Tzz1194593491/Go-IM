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
}

func NewUser(conn net.Conn) *User {
	userAddress := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddress,
		Address: userAddress,
		C:       make(chan string),
		conn:    conn,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
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
