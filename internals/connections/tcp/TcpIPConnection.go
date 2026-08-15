package tcp

import (
	"Raven/internals/workers"
	"bufio"
	"fmt"
	"net"
	"strings"
)

type TcpIpConnectionHandler struct {
	workerPool *workers.WorkerPool
	Ttl        uint32 //in mellisecends
	Conn       net.Conn
}

func NewTcpConnectionHandler(wp *workers.WorkerPool) *TcpIpConnectionHandler {
	return &TcpIpConnectionHandler{
		workerPool: wp,
		Ttl:        3000,
	}
}

func (tcp_ip_handler *TcpIpConnectionHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	tcp_ip_handler.Conn = conn
	// Greeting message to connected client
	conn.Write([]byte("Connected to Raven DB Server (Dynamic AST Parser & 4-Worker Engine)\n"))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		resChan := make(chan workers.Response, 1)
		job := workers.Job{
			Query:        line,
			ResponseChan: resChan,
		}

		// Submit raw query string into the 4-worker pool queue
		tcp_ip_handler.workerPool.JobQueue <- job

		// Wait for response from worker goroutine
		res := <-resChan

		if res.Err != nil {
			conn.Write([]byte(fmt.Sprintf("%v\n", res.Err)))
		} else {
			conn.Write([]byte(res.Result))
		}
	}
}
