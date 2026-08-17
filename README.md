# Raven DB - Redis Clone in Go

```text
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
```

**Raven DB** is a high-performance in-memory key-value database engine written in Go. It features a **configurable worker pool architecture**, **goroutine query channels**, **thread-safe data storage**, **Cobra & Viper CLI/Config management**, **interactive REPL & CLI subcommands**, and a **dynamic programming-language-style AST Command Parser & Evaluator Engine** supporting list literals and nested command expressions.

---

## 🌟 Architectural Features

* **Configurable Worker Pool Architecture:** Client TCP connections submit raw query jobs to a buffered `JobQueue`. Worker count is fully configurable via Cobra CLI flags, environment variables, or config files (default: 4 workers).
* **Per-Query Goroutine Response Channels:** Each query contains a per-request response channel (`chan Response`), ensuring strict isolation and response routing back to the client TCP connection.
* **Thread-Safe Storage Engine:** Memory operations are guarded by fine-grained `sync.RWMutex` locks, preventing data races across concurrent workers.
* **Cobra & Viper CLI System**:
  * **Server (`ravendb`)**: Supports `start`, `version`, and flags (`-w`/`--workers`, `-p`/`--port`, `-q`/`--queue`, `-c`/`--config`).
  * **Client (`raven-cli`)**: Supports interactive REPL mode, single-command execution mode (`raven-cli exec`), `ping`, `version`, and host/port flags (`-H`/`--host`, `-p`/`--port`).
* **Struct-bound AST Lexer, Parser & Evaluator:** Uses object-oriented `Lexer`, `Parser`, and `Evaluator` structs capable of parsing:
  * Quoted strings with spaces: `SET title "Raven Database Engine"`
  * Nested command expressions: `SET backup (GET title)`
  * List Literals & Brackets: `SET tags [ "server" (GET env) 8080 ]`
  * Number & String literals
* **Interface-based Command System:** Commands implement a clean `Command` interface (`type Command interface { Execute(args []string) (string, error) }`), each residing in its own modular file inside `internals/commands/`.
* **Rich Colorized CLI Visualization**: Interactive REPL features ANSI colorized status outputs (bold green `OK`, bold yellow `(integer) N`, dim gray `(nil)`, bold red `ERR`), multi-line array item numbering, and a custom Raven Bird ASCII logo.

---

## 🏗️ System Architecture Flow

```mermaid
graph TD
    Client["Client (raven-cli / TCP Socket)"] -->|Query Line Stream| ConnHandler["TcpIpConnectionHandler"]
    ConnHandler -->|Job{Query, ResponseChan}| JobQueue["Configurable Worker Job Queue"]
    JobQueue --> Worker["Worker (1 of N Goroutines)"]
    Worker --> Lexer["Lexer.Tokenize()"]
    Lexer --> Parser["Parser.Parse() (CommandExpr / ListLiteral)"]
    Parser --> Evaluator["Evaluator.EvaluateQuery()"]
    Evaluator <--> Registry["Dynamic Command Registry"]
    Registry <--> DataStore["Thread-Safe DataStore (sync.RWMutex)"]
    Evaluator -- Result String --> ConnHandler
    ConnHandler -- Response Bytes --> Client
```

---

## 📖 Supported Commands Reference

| Command | Usage | Description | Example |
| :--- | :--- | :--- | :--- |
| `PING` | `PING [msg]` | Health check | `PING` -> `PONG` |
| `ECHO` | `ECHO msg` | Prints input message | `ECHO "Hello"` |
| `SET` | `SET key value [ttl_seconds]` | Sets key to value with optional TTL | `SET msg "Hello World"` |
| `GET` | `GET key` | Retrieves key value | `GET msg` |
| `DEL` | `DEL key [key2 ...]` | Deletes one or more keys | `DEL msg` -> `(integer) 1` |
| `EXISTS` | `EXISTS key [key2 ...]` | Checks if keys exist | `EXISTS msg` -> `(integer) 1` |
| `KEYS` | `KEYS` | Lists all active keys | `KEYS` |
| `EXPIRE` | `EXPIRE key seconds` | Sets expiration TTL in seconds | `EXPIRE msg 60` |
| `TTL` | `TTL key` | Returns remaining TTL in seconds | `TTL msg` -> `(integer) 58` |
| `INCR` | `INCR key` | Increments integer key by 1 | `INCR counter` -> `(integer) 1` |
| `DECR` | `DECR key` | Decrements integer key by 1 | `DECR counter` -> `(integer) 0` |
| `MGET` | `MGET key1 key2 ...` | Gets multiple key values | `MGET k1 k2` |
| `MSET` | `MSET k1 v1 k2 v2 ...` | Sets multiple key-value pairs | `MSET k1 "v1" k2 "v2"` |

---

## 🚀 Dynamic Language Syntax Examples

### 1. Quoted Strings
```text
raven> SET greeting "Hello World from Raven DB"
OK
raven> GET greeting
Hello World from Raven DB
```

### 2. Nested Sub-Command Expressions
Dynamic sub-expressions inside parentheses `(...)` evaluate recursively first and pass results into outer commands:
```text
raven> SET default_role "Administrator"
OK
raven> SET user_role (GET default_role)
OK
raven> GET user_role
Administrator
```

### 3. List Literals with Nested Sub-Commands
List literals `[...]` parse square brackets and recursively evaluate nested expressions:
```text
raven> SET env "production"
OK
raven> SET server_config [ "server_v1" (GET env) 8080 ]
OK
raven> GET server_config
[server_v1 production 8080]
```

---

## ⚙️ How to Register Custom Commands

Commands implement the `Command` interface and are registered dynamically in the `CommandRegistry`:

```go
package main

import (
	"Raven/internals/commands"
	"strings"
)

// Define a custom command struct
type UppercaseCommand struct{}

func (c *UppercaseCommand) Execute(args []string) (string, error) {
	return strings.ToUpper(strings.Join(args, " ")) + "\n", nil
}

func main() {
	// Register the custom command
	commands.GlobalRegistry.Register("UPPERCASE", &UppercaseCommand{})
}
```

---

## 🛠️ Building & Running

### 1. Build Binaries
```bash
go build -o ravendb ./cmd/database
go build -o raven-cli ./cmd/cli
```

### 2. Start Raven DB Server (`ravendb`)
Launch with default settings (4 workers, port 7777, queue size 100):
```bash
./ravendb start
```

Launch with custom Cobra CLI flags:
```bash
./ravendb start -w 8 -p 7778 -q 200
```

Launch using environment variables:
```bash
RAVEN_WORKERS=16 RAVEN_PORT=8888 ./ravendb start
```

Launch using a configuration file:
```bash
./ravendb start --config config.example.yaml
```

### 3. Connect via Raven CLI (`raven-cli`)
Start interactive REPL mode:
```bash
./raven-cli
# Or connect to custom host/port
./raven-cli -H 127.0.0.1 -p 7778
```

Execute a single command directly:
```bash
./raven-cli exec "SET title 'Raven DB'"
./raven-cli exec "GET title"
./raven-cli ping
```

---

## 🧪 Running Unit & Race Detector Tests

Run the complete unit and feature test suite across all packages:
```bash
go test -count=1 -v ./...
```

RunWith Go's data race detector:
```bash
go test -v -race ./...
```
