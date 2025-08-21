package server

import (
	"net"
)

type User struct {
	conn net.Conn // 链接信息
	// 用户名
	Name string
	// 用户id
	Id int64
	// 缓冲通道，用来接受消息
	C chan string
	// ip地址
	Addr string
}

// NewUser 创建一个新用户
func NewUser(con net.Conn, id int64) *User {
	addr := con.RemoteAddr().String()
	user := &User{
		conn: con,
		Name: addr,
		Id:   id,
		Addr: addr,
		C:    make(chan string),
	}

	// 启动监听用户消息的go程
	go user.ListenerMessage()

	return user
}

// ListenerMessage 监听用户的信息
func (this *User) ListenerMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
