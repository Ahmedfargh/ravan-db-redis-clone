package main

import (
	"fmt"
	"os"
	"strings"

	data_store "Raven/internals/database"
	servers "Raven/internals/servers"
	"Raven/internals/workers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	workerCount int
	queueSize   int
	port        int
)

var rootCmd = &cobra.Command{
	Use:   "ravendb",
	Short: "Raven DB - High Performance In-Memory Key-Value Store",
	Long:  `Raven DB is an in-memory key-value database engine featuring a worker pool architecture, goroutine response channels, and dynamic AST command parsing.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")

	rootCmd.Flags().IntVarP(&workerCount, "workers", "w", 4, "Number of worker goroutines in worker pool")
	rootCmd.Flags().IntVarP(&queueSize, "queue", "q", 100, "Buffer size of the query job queue")
	rootCmd.Flags().IntVarP(&port, "port", "p", 7777, "Port number for the TCP server")

	viper.BindPFlag("workers", rootCmd.Flags().Lookup("workers"))
	viper.BindPFlag("queue", rootCmd.Flags().Lookup("queue"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	viper.SetDefault("workers", 4)
	viper.SetDefault("queue", 100)
	viper.SetDefault("port", 7777)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("RAVEN")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func runServer() {
	workersNum := viper.GetInt("workers")
	qSize := viper.GetInt("queue")
	portNum := viper.GetInt("port")

	if workersNum <= 0 {
		fmt.Fprintf(os.Stderr, "Error: worker count must be greater than 0\n")
		os.Exit(1)
	}

	// 1. Initialize thread-safe data store map
	data_store.InitiatDataStore()

	// 2. Spawn worker pool with configured worker count and queue size
	workerPool := workers.NewWorkerPool(workersNum, qSize)

	// 3. Start TCP server with worker pool and configured port
	tcp_server := servers.NewTcpServer(workerPool, portNum)
	tcp_server.HandleConnections()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
