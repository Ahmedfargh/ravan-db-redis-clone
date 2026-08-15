# Raven DB - Redis Clone in Go

**Raven DB** is a high-performance in-memory key-value database engine written in Go. It features a **4-worker pool architecture**, **goroutine query channels**, **thread-safe data storage**, and a **dynamic programming-language-style AST Command Parser & Evaluator Engine**.

---

## 🌟 Architectural Features

* **4-Worker Pool Architecture:** Client TCP connections submit raw query jobs to a buffered `JobQueue`. Four dedicated worker goroutines execute queries concurrently.
* **Per-Query Goroutine Response Channels:** Each query contains a per-request response channel (`chan Response`), ensuring strict isolation and response routing back to the client TCP connection.
* **Thread-Safe Storage Engine:** Memory operations are guarded by fine-grained `sync.RWMutex` locks, preventing data races across concurrent workers.
* **Dynamic Programming-Language Command Parser:** Uses a Lexer, AST Parser, and Recursive Evaluator capable of parsing:
  - Quoted strings with spaces: `SET title "Raven Database Engine"`
  - Nested command expressions: `SET backup (GET title)`
  - Number & String literals
* **Pluggable Command Registry:** Commands are registered dynamically (`commands.GlobalRegistry.Register("MYCMD", handler)`), eliminating hardcoded switch blocks.

---

## 🏗️ System Architecture Flow

```mermaid
graph TD
    Client["Client (raven-cli / TCP Socket)"] -->|Query Line Stream| ConnHandler["TcpIpConnectionHandler"]
    ConnHandler -->|Job{Query, ResponseChan}| JobQueue["4-Worker Job Queue"]
    JobQueue --> Worker["Worker (1 of 4 Goroutines)"]
    Worker --> Lexer["Lexer (Token Stream)"]
    Lexer --> Parser["AST Parser (CommandExpr)"]
    Parser --> Evaluator["AST Evaluator"]
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
| `SET` | `SET key value [ttl_seconds]` | Sets key to value with optional TTL | `SET msg "Hello World"` |
| `GET` | `GET key` | Retrieves key value | `GET msg` |
| `DEL` | `DEL key [key2 ...]` | Deletes one or more keys | `DEL msg` |
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

## ⚙️ How to Register Custom Commands Dynamically

You can dynamically extend Raven DB with custom commands without modifying the core worker or parser logic:

```go
package main

import (
	"Raven/internals/commands"
	"strings"
)

func main() {
	// Register a custom dynamic command "UPPERCASE"
	commands.GlobalRegistry.Register("UPPERCASE", func(args []string) (string, error) {
		return strings.ToUpper(strings.Join(args, " ")) + "\n", nil
	})
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
```bash
./ravendb
```
*Output:*
```text
Server initiated successfully on address: 127.0.0.1:7777
Worker Pool active: 4 worker goroutines listening on query queue
```

### 3. Connect via Raven CLI
```bash
./raven-cli
```
*Output:*
```text
Connected to Raven DB Server (Dynamic AST Parser & 4-Worker Engine)
raven> SET name "Raven"
OK
raven> GET name
Raven
```

---

## 🧪 Running Unit & Race Detector Tests

Run the full test suite with Go's data race detector:
```bash
go test -v -race ./...
```
# ravan-db-redis-clone
