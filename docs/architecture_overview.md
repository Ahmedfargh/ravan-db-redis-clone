# Raven DB - Architecture & Behavioral Specification

This document provides a comprehensive technical breakdown of Raven DB's finite state machines, query processing lifecycle, and object-oriented class architecture.

---

## 📁 Draw.io Diagram File Reference Index

All diagrams are provided in native Draw.io XML format under the [`docs/`](file:///home/ahmed/golang/redis-clone/docs/) directory. They can be opened and edited directly via [Draw.io / diagrams.net](https://app.diagrams.net/):

- 🏆 **Unified Master Workbook**: [`docs/raven_db_architecture_all.drawio`](file:///home/ahmed/golang/redis-clone/docs/raven_db_architecture_all.drawio) (Contains all diagrams across tabbed pages: Master Overview, UML Architecture, Syntax Evaluation Sequence, Job Lifecycle FSM, Lexer FSM, AST Parser FSM).

### Individual Subsystem Files:
1. 🟢 [`docs/syntax_evaluation_lifecycle.drawio`](file:///home/ahmed/golang/redis-clone/docs/syntax_evaluation_lifecycle.drawio): End-to-end query processing sequence flow.
2. 🟡 [`docs/lexer_fsm.drawio`](file:///home/ahmed/golang/redis-clone/docs/lexer_fsm.drawio): Lexer finite state machine (character stream to token scanner).
3. 🔵 [`docs/parser_fsm.drawio`](file:///home/ahmed/golang/redis-clone/docs/parser_fsm.drawio): AST Parser finite state machine (token stream to `CommandExpr` AST).
4. 🟣 [`docs/job_lifecycle_fsm.drawio`](file:///home/ahmed/golang/redis-clone/docs/job_lifecycle_fsm.drawio): Concurrency worker pool query job lifecycle state machine.
5. 🔴 [`docs/system_architecture_uml.drawio`](file:///home/ahmed/golang/redis-clone/docs/system_architecture_uml.drawio): Complete system class and component UML diagram.

---

## 1. 🔤 Lexer Finite State Machine ([`lexer_fsm.drawio`](file:///home/ahmed/golang/redis-clone/docs/lexer_fsm.drawio))

The [`Lexer`](file:///home/ahmed/golang/redis-clone/internals/parser/lexer.go) scans incoming raw query strings character-by-character (`ch byte`) and emits a stream of typed tokens (`[]Token`).

### States & Transitions

| State | Trigger / Condition | Action / Output | Next State |
| :--- | :--- | :--- | :--- |
| **`Initial`** | `NewLexer(input)` | Sets position pointers, reads first char | `ReadNextChar` |
| **`SkipWhitespace`** | `ch` is `' '`, `\t`, `\n`, `\r` | Advances pointer `readChar()` | `ReadNextChar` |
| **`StateIdent`** | `ch` is letter, `_`, `:`, `*`, `-` | Accumulates bytes until non-ident character | Emits `TOKEN_IDENT` |
| **`StateNumber`** | `ch` is digit `'0'`-`'9'` | Accumulates digits and decimals `.` | Emits `TOKEN_NUMBER` |
| **`StateString`** | `ch` is quote `"` or `'` | Skips starting quote, handles escape `\`, skips ending quote | Emits `TOKEN_STRING` |
| **`StateParen`** | `ch` is `(` or `)` | Identifies structural parenthesis | Emits `TOKEN_LPAREN` / `TOKEN_RPAREN` |
| **`StateEOF`** | `ch == 0` | End of string input | Emits `TOKEN_EOF` |
| **`StateIllegal`** | Unrecognized byte | Marks character as invalid | Emits `TOKEN_ILLEGAL` |

---

## 2. 🌲 AST Parser Finite State Machine ([`parser_fsm.drawio`](file:///home/ahmed/golang/redis-clone/docs/parser_fsm.drawio))

The [`Parser`](file:///home/ahmed/golang/redis-clone/internals/parser/parser.go) processes token streams and constructs an Abstract Syntax Tree (AST) rooted at `*CommandExpr`.

### Parsing Algorithm Flow

1. **Check Parentheses**: If `curToken` is `TOKEN_LPAREN` (`(`), set flag `isEnclosedInParen = true` and consume token.
2. **Command Identification**: `curToken` MUST be `TOKEN_IDENT` (e.g., `SET`, `GET`, `PING`). Converts command name to uppercase.
3. **Expression Loop**: Loops over remaining tokens until `TOKEN_EOF` or matching `TOKEN_RPAREN` (`)`):
   - **Literals** (`TOKEN_STRING`, `TOKEN_NUMBER`, `TOKEN_IDENT`): Instantiates `StringLiteral` or `NumberLiteral` node and appends to `CommandExpr.Args`.
   - **Nested Commands** (`TOKEN_LPAREN`): Recursively invokes `p.ParseCommand()` to construct a child `*CommandExpr` node and appends it to `CommandExpr.Args`.
4. **Completion**: Returns complete `*CommandExpr` node tree.

---

## 3. ⚡ Query Job Lifecycle State Machine ([`job_lifecycle_fsm.drawio`](file:///home/ahmed/golang/redis-clone/docs/job_lifecycle_fsm.drawio))

Demonstrates the lifecycle of a query request from the client TCP connection through the goroutine worker pool.

```text
[StateSocketRead] ──(push to queue)──> [StateJobQueued] ──(pop by worker)──> [StateWorkerPicked]
                                                                                      │
                                                                           (NewParserFromInput)
                                                                                      │
[StateResponseDispatched] <──(send result)── [StateDataStoreMutating] <──(Evaluate)── [StateEvaluating]
```

1. **`StateSocketRead`**: `TcpIpConnectionHandler` receives line from TCP connection.
2. **`StateJobQueued`**: Instantiates `Job{Query, ResponseChan}` and pushes to buffered `JobQueue`.
3. **`StateWorkerPicked`**: One of $N$ worker goroutines pops `Job` from `JobQueue`.
4. **`StateLexingParsing`**: Worker invokes `NewParserFromInput(query).Parse()` to obtain AST.
5. **`StateEvaluating`**: `Evaluator.EvaluateQuery()` recursively evaluates sub-expressions.
6. **`StateDataStoreMutating`**: Target `Command.Execute()` reads/mutates `DataStore` under `sync.RWMutex`.
7. **`StateResponseDispatched`**: Worker sends `Response{Result}` back on `job.ResponseChan` to write to TCP socket.

---

## 4. 🏗️ Class Architecture UML ([`system_architecture_uml.drawio`](file:///home/ahmed/golang/redis-clone/docs/system_architecture_uml.drawio))

### Key System Interfaces & Structs

- **`Command` Interface**:
  ```go
  type Command interface {
      Execute(args []string) (string, error)
  }
  ```
  Implemented by concrete structs: `PingCommand`, `SetCommand`, `GetCommand`, `DelCommand`, `ExistsCommand`, `KeysCommand`, `ExpireCommand`, `TtlCommand`, `IncrCommand`, `DecrCommand`, `MgetCommand`, `MsetCommand`, `EchoCommand`.

- **`CommandRegistry`**: Holds `map[string]Command` guarded by `sync.RWMutex`. Thread-safe dynamic lookup and registration.
- **`Evaluator`**: Dependency-injected evaluator executing queries against a `CommandRegistry`.
- **`DataStore`**: Thread-safe storage engine holding `map[string]*DataValue` guarded by `sync.RWMutex`.
