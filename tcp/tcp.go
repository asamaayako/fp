package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {

	// 创建TCP监听器
	l, err := net.Listen("tcp", ":443") // 使用默认的HTTPS端口，或者你选择的其他端口
	if err != nil {
		log.Fatal(err)
	}

	// 将TCP连接转换为TLS连接
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

type ConnManager struct {
	net.ListenConfig
	Conns []net.Conn
}

func (c *ConnManager) Listen(listenAddr string) {

	// 创建TCP监听器
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on", listenAddr)

	for {
		// 接受新的连接请求
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取客户端发送的数据
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading:", err)
		return
	}

	// 打印接收到的数据
	fmt.Printf("Received from client: %s\n", buffer[:n])

	// 回复客户端
	message := "Hello, Client!"
	conn.Write([]byte(message))
}
