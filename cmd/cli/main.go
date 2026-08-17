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

const (
	cliBanner = `
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
            R A V E N  C L I
`
	colorReset   = "\033[0m"
	colorRed     = "\033[1;31m"
	colorGreen   = "\033[1;32m"
	colorYellow  = "\033[1;33m"
	colorBlue    = "\033[1;34m"
	colorMagenta = "\033[1;35m"
	colorCyan    = "\033[1;36m"
	colorDim     = "\033[2;37m"
	colorBold    = "\033[1m"
)

var rootCmd = &cobra.Command{
	Use:   "raven-cli",
	Short: "Raven DB Interactive CLI Client",
	Long:  cliBanner + "\nInteractive CLI client for connecting to Raven DB key-value database.",
	Run: func(cmd *cobra.Command, args []string) {
		runREPL()
	},
}

var execCmd = &cobra.Command{
	Use:   "exec [command]",
	Short: "Execute a single command directly and print the response",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		runSingleCommand(query)
	},
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the Raven DB server",
	Run: func(cmd *cobra.Command, args []string) {
		runSingleCommand("PING")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Raven DB CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Raven DB CLI v1.0.0")
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Server host address")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 7777, "Server port")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", 7777)

	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(pingCmd)
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

	_ = viper.ReadInConfig()
}

func formatResponse(response string) string {
	trimmed := strings.TrimRight(response, "\r\n")

	// OK
	if trimmed == "OK" {
		return fmt.Sprintf("%sOK%s", colorGreen, colorReset)
	}

	// Null
	if trimmed == "(nil)" {
		return fmt.Sprintf("%s(nil)%s", colorDim, colorReset)
	}

	// Empty set/list
	if trimmed == "(empty list or set)" {
		return fmt.Sprintf("%s(empty list or set)%s", colorYellow, colorReset)
	}

	// Integers
	if strings.HasPrefix(trimmed, "(integer)") {
		return fmt.Sprintf("%s%s%s", colorYellow, trimmed, colorReset)
	}

	// Errors
	if strings.HasPrefix(trimmed, "ERR") || strings.HasPrefix(trimmed, "(error)") {
		return fmt.Sprintf("%s%s%s", colorRed, trimmed, colorReset)
	}

	// Multi-line list/array responses
	lines := strings.Split(trimmed, "\n")
	if len(lines) > 1 {
		var formatted []string
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			formatted = append(formatted, fmt.Sprintf("%s%d)%s %s%s%s", colorCyan, i+1, colorReset, colorBold, line, colorReset))
		}
		return strings.Join(formatted, "\n")
	}

	return fmt.Sprintf("%s%s%s", colorCyan, trimmed, colorReset)
}

func connectServer() (net.Conn, *bufio.Reader, error) {
	serverHost := viper.GetString("host")
	serverPort := viper.GetInt("port")
	serverAddr := fmt.Sprintf("%s:%d", serverHost, serverPort)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Raven DB server at %s: %w", serverAddr, err)
	}

	connReader := bufio.NewReader(conn)
	// Consume greeting line from server
	_, _ = connReader.ReadString('\n')

	return conn, connReader, nil
}

func runSingleCommand(query string) {
	conn, connReader, err := connectServer()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(query + "\n"))
	if err != nil {
		fmt.Printf("%s[ERROR]%s Failed to send command: %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}

	response, err := connReader.ReadString('\n')
	if err != nil {
		fmt.Printf("%s[ERROR]%s Server closed connection: %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}

	fmt.Println(formatResponse(response))
}

func runREPL() {
	serverHost := viper.GetString("host")
	serverPort := viper.GetInt("port")

	conn, connReader, err := connectServer()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Print(colorCyan + cliBanner + colorReset)
	fmt.Printf("%sConnected to Raven DB Server at %s:%d%s\n", colorGreen, serverHost, serverPort, colorReset)
	fmt.Printf("%sType 'help' for command syntax, 'exit' or 'quit' to exit.%s\n\n", colorDim, colorReset)

	prompt := fmt.Sprintf("%sraven [%s:%d]> %s", colorGreen, serverHost, serverPort, colorReset)
	stdinScanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)

	for stdinScanner.Scan() {
		input := strings.TrimSpace(stdinScanner.Text())
		if input == "" {
			fmt.Print(prompt)
			continue
		}

		lower := strings.ToLower(input)
		if lower == "exit" || lower == "quit" {
			fmt.Printf("%sBye!%s\n", colorYellow, colorReset)
			break
		}

		if lower == "clear" {
			fmt.Print("\033[H\033[2J")
			fmt.Print(prompt)
			continue
		}

		// Send query to server
		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Printf("%s[ERROR]%s Failed to write to server: %v\n", colorRed, colorReset, err)
			break
		}

		// Read response
		response, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("%s[ERROR]%s Server connection lost: %v\n", colorRed, colorReset, err)
			break
		}

		fmt.Println(formatResponse(response))
		fmt.Print(prompt)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
