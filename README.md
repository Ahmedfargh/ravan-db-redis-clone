# Raven DB - Redis Clone in Go

**Raven DB** is a high-performance in-memory key-value database engine written in Go. It features a **configurable worker pool architecture**, **goroutine query channels**, **thread-safe data storage**, **Cobra & Viper CLI/Config management**, and a **dynamic programming-language-style AST Command Parser & Evaluator Engine**.

---

## 🌟 Architectural Features

* **Configurable Worker Pool Architecture:** Client TCP connections submit raw query jobs to a buffered `JobQueue`. Worker count is fully configurable via CLI flags, environment variables, or config files (default: 4 workers).
* **Per-Query Goroutine Response Channels:** Each query contains a per-request response channel (`chan Response`), ensuring strict isolation and response routing back to the client TCP connection.
* **Thread-Safe Storage Engine:** Memory operations are guarded by fine-grained `sync.RWMutex` locks, preventing data races across concurrent workers.
* **Struct-bound AST Lexer, Parser & Evaluator:** Uses object-oriented `Lexer`, `Parser`, and `Evaluator` structs capable of parsing:
  - Quoted strings with spaces: `SET title "Raven Database Engine"`
  - Nested command expressions: `SET backup (GET title)`
  - Number & String literals
* **Interface-based Command System:** Commands implement a clean `Command` interface (`type Command interface { Execute(args []string) (string, error) }`), each residing in its own modular file inside `internals/commands/`.
* **Cobra & Viper CLI Configuration:** Fully supports configuration via **CLI flags**, **Environment variables** (`RAVEN_WORKERS`, `RAVEN_PORT`, `RAVEN_QUEUE`), and **Config Files** (`config.yaml`).

---

## 🏗️ System Architecture Flow

```mermaid
graph TD
    Client["Client (raven-cli / TCP Socket)"] -->|Query Line Stream| ConnHandler["TcpIpConnectionHandler"]
    ConnHandler -->|Job{Query, ResponseChan}| JobQueue["Configurable Worker Job Queue"]
    JobQueue --> Worker["Worker (1 of N Goroutines)"]
    Worker --> Lexer["Lexer.Tokenize()"]
    Lexer --> Parser["Parser.Parse() (CommandExpr)"]
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

### Quoted Strings
```text
raven> SET greeting "Hello World from Raven DB"
OK
raven> GET greeting
Hello World from Raven DB
```

### Nested Sub-Command Expressions
Dynamic sub-expressions inside parentheses `(...)` evaluate recursively first and pass results into outer commands:
```text
raven> SET default_role "Administrator"
OK
raven> SET user_role (GET default_role)
OK
raven> GET user_role
Administrator
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

### 2. Start Raven DB Server
Launch with default settings (4 workers, port 7777, queue size 100):
```bash
./ravendb
```

Launch with custom flags:
```bash
./ravendb -w 8 -p 7777 -q 200
```

Launch using environment variables:
```bash
RAVEN_WORKERS=16 RAVEN_PORT=8888 ./ravendb
```

Launch using a configuration file:
```bash
./ravendb --config config.example.yaml
```

### 3. Connect via Raven CLI
```bash
./raven-cli
# Or connect to a custom port
./raven-cli -p 8888
```

---

## 🧪 Running Unit & Race Detector Tests

Run the full test suite with Go's data race detector:
```bash
go test -v -race ./...
```
