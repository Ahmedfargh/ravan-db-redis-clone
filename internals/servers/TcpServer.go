package servers

import (
	tcp_connection_handler "Raven/internals/connections/tcp"
	"Raven/internals/workers"
	"fmt"
	"net"
)

type TcpServer struct {
	listener   net.Listener
	workerPool *workers.WorkerPool
}

func NewTcpServer(wp *workers.WorkerPool) *TcpServer {
	return &TcpServer{
		workerPool: wp,
	}
}

func (tcp_server *TcpServer) initatServer() {
	address := "localhost:7777"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	tcp_server.listener = listener
	fmt.Println("Server initiated successfully on address:", listener.Addr())
	fmt.Printf("Worker Pool active: %d worker goroutines listening on query queue\n", tcp_server.workerPool.WorkerCount)
}

func (tcp_server *TcpServer) HandleConnections() {
	tcp_server.initatServer()
	tcp_handler := tcp_connection_handler.NewTcpConnectionHandler(tcp_server.workerPool)

	for {
		conn, err := tcp_server.listener.Accept()
		if err != nil {
			fmt.Println("failed to establish connection:", err)
			continue
		}

		fmt.Println("Connection established with:", conn.RemoteAddr())
		go tcp_handler.HandleConnection(conn)
	}
}
