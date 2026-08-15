package main

import (
	data_store "Raven/internals/database"
	servers "Raven/internals/servers"
	"Raven/internals/workers"
)

func main() {
	// 1. Initialize thread-safe data store map
	data_store.InitiatDataStore()

	// 2. Spawn 4 worker goroutines listening on the query job queue (size 100)
	workerPool := workers.NewWorkerPool(4, 100)

	// 3. Start TCP server with worker pool
	tcp_server := servers.NewTcpServer(workerPool)
	tcp_server.HandleConnections()
}
