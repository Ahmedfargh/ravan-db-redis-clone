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

const (
	banner = `
                  /\
                 /  \
                / /  \
               | (o)  |======>
                \    /
               /      \
              /  /\    \
             /  /  \    \
            /  /    \    \
           /  /      \    \
          |  |        |    |
          |__|        |____|
       ========================
            R A V E N  D B
  High Performance In-Memory Store
`
	colorReset  = "\033[0m"
	colorCyan   = "\033[1;36m"
	colorGreen  = "\033[1;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[1;34m"
)

var rootCmd = &cobra.Command{
	Use:   "ravendb",
	Short: "Raven DB - High Performance In-Memory Key-Value Store",
	Long:  banner + "\nRaven DB is an in-memory key-value database engine featuring a worker pool architecture, goroutine response channels, and dynamic AST command parsing.",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Raven DB server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Raven DB engine version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Raven DB Engine v1.0.0 (AST Parser + Worker Pool architecture)")
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().IntVarP(&workerCount, "workers", "w", 4, "Number of worker goroutines in worker pool")
	rootCmd.PersistentFlags().IntVarP(&queueSize, "queue", "q", 100, "Buffer size of the query job queue")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 7777, "Port number for the TCP server")

	viper.BindPFlag("workers", rootCmd.PersistentFlags().Lookup("workers"))
	viper.BindPFlag("queue", rootCmd.PersistentFlags().Lookup("queue"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	viper.SetDefault("workers", 4)
	viper.SetDefault("queue", 100)
	viper.SetDefault("port", 7777)

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
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
		fmt.Printf("%s[CONFIG]%s Loaded configuration file: %s\n", colorBlue, colorReset, viper.ConfigFileUsed())
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

	fmt.Print(colorCyan + banner + colorReset)
	fmt.Printf("%s[INIT]%s Initializing thread-safe DataStore...\n", colorGreen, colorReset)
	data_store.InitiatDataStore()

	fmt.Printf("%s[WORKERS]%s Spawning %d worker goroutines (Job queue capacity: %d)...\n", colorYellow, colorReset, workersNum, qSize)
	workerPool := workers.NewWorkerPool(workersNum, qSize)

	fmt.Printf("%s[SERVER]%s Starting TCP Server listener on port %d...\n", colorBlue, colorReset, portNum)
	tcp_server := servers.NewTcpServer(workerPool, portNum)
	if err := tcp_server.HandleConnections(); err != nil {
		fmt.Printf("\n%s[ERROR]%s Could not start Raven DB server: %v\n", "\033[1;31m", colorReset, err)
		fmt.Printf("%s[TIP]%s Port %d may already be in use. Use '-p <port>' or '--port <port>' to specify a different port (e.g. ./ravendb -p 7778)\n\n", colorYellow, colorReset, portNum)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
