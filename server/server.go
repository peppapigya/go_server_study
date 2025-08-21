package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip       string
	Port     int
	OlineMap map[string]*User
	MapLock  sync.RWMutex
	Message  chan string
}

// NewServer 建立连接
func NewServer(Ip string, Port int) *Server {
	fmt.Println("建立连接..............")
	return &Server{
		Ip:       Ip,
		Port:     Port,
		OlineMap: make(map[string]*User),
		Message:  make(chan string),
	}
}

// Handler 消息处理器
func (server *Server) Handler(con net.Conn) {
	fmt.Println("读取对应的数据........")
	user := NewUser(con, 1)
	// 上锁
	server.MapLock.Lock()
	// 用户上线将信息设置到上线用户列表中
	server.OlineMap[user.Name] = user
	server.MapLock.Unlock()

	// 启动 go 程去发送消息
	go server.SendMessage(con, "上线")

	//接收客户端发来的消息,模拟用户去发送消息
	go server.ReceiveMessage(con)
	// 阻塞当前 go 程，防止当前go程退出所有的信息都结束了
	select {}
}

// MessageListener 监听Message 的go程，一但有消息就直接广播给其他用户
func (server *Server) MessageListener() {
	for {
		msg := <-server.Message
		server.MapLock.Lock()
		for _, user := range server.OlineMap {
			user.C <- msg
		}
		server.MapLock.Unlock()

	}
}

// SendMessage 发送消息
func (server *Server) SendMessage(conn net.Conn, msg string) {
	addr := conn.RemoteAddr().String()
	server.Message <- "[" + addr + "]:" + msg
}

// ReceiveMessage 接收客户端消息
func (server *Server) ReceiveMessage(conn net.Conn) {
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Conn Read Error", err)
			return
		}
		if n == 0 {
			fmt.Println("客户端退出")
			return
		}
		// 将消息解析成字符串
		msg := string(buf[:n-1])
		server.SendMessage(conn, msg)
	}
}

// Start 开启服务器
func (server *Server) Start() {
	// 建立连接
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("与服务器建立连接失败")
		return
	}
	// 监听链接
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("监听关闭失败.......")
		}
	}(listener)
	// 启动用户消息监听
	go server.MessageListener()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("与客户端建立连接失败")
			continue
		}

		// 获取信息
		go server.Handler(conn)
	}
}
