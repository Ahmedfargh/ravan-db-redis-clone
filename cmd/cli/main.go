package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	host    string
	port    int
)

var rootCmd = &cobra.Command{
	Use:   "raven-cli",
	Short: "Raven DB Interactive CLI Client",
	Run: func(cmd *cobra.Command, args []string) {
		runCLI()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")
	rootCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "Server host address")
	rootCmd.Flags().IntVarP(&port, "port", "p", 7777, "Server port")

	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	viper.SetDefault("host", "127.0.0.1")
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
		// Loaded config file
	}
}

func runCLI() {
	serverHost := viper.GetString("host")
	serverPort := viper.GetInt("port")
	serverAddr := fmt.Sprintf("%s:%d", serverHost, serverPort)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to Raven DB server at %s: %v\n", serverAddr, err)
		os.Exit(1)
	}
	defer conn.Close()

	connReader := bufio.NewReader(conn)
	// Read server welcome message
	greeting, err := connReader.ReadString('\n')
	if err == nil {
		fmt.Print(greeting)
	}

	stdinScanner := bufio.NewScanner(os.Stdin)
	fmt.Print("raven> ")

	for stdinScanner.Scan() {
		input := strings.TrimSpace(stdinScanner.Text())
		if input == "" {
			fmt.Print("raven> ")
			continue
		}
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Bye!")
			break
		}

		// Send command to server
		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Printf("Error writing to server: %v\n", err)
			break
		}

		// Read response from server
		response, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Connection closed by server: %v\n", err)
			break
		}

		fmt.Print(response)
		fmt.Print("raven> ")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
