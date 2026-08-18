# Raven Query Language (RQL) - Specification & Syntax Guide

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
            R A V E N  Q L
  Raven Query Language Specification
```

**Raven Query Language (RQL)** is the formal query language of Raven DB. It combines **Redis command semantics** with a **Lisp/S-expression AST evaluation model**, supporting quoted literals, nested sub-queries, and array list expressions.

---

## 📖 1. Formal Grammar & Syntax Rules

### Rule 1: Command Expression Structure
Queries follow an imperative command-first pattern consisting of a **Command Name** followed by zero or more space-separated arguments:
```text
COMMAND_NAME [arg1] [arg2] ... [argN]
```

### Rule 2: Case Insensitivity
Command names are case-insensitive. `set`, `Set`, and `SET` evaluate identically:
```text
set key1 "value"
SET key1 "value"
```

### Rule 3: Literal Types
1. **Identifiers / Unquoted Strings**: Continuous alphanumeric sequences (supports `_`, `:`, `*`, `-`).
2. **Quoted Strings**: Enclosed in single `'...'` or double `"..."` quotes. Supports spaces and escape sequences (`\"`, `\\`).
3. **Number Literals**: Continuous digits and decimal points (e.g. `100`, `3.14`).

### Rule 4: Sub-Command Expressions `(...)`
Sub-queries enclosed in parentheses `( ... )` evaluate **inside-out** recursively before the outer command executes:
```text
SET user_role (GET default_role)
```

### Rule 5: List Literals `[...]`
Values enclosed in square brackets `[ ... ]` construct structured lists containing strings, numbers, or sub-queries:
```text
SET server_config [ "production" (GET port_num) 8080 ]
```

---

## 🌲 2. Abstract Syntax Tree (AST) Specification

RQL queries are parsed into a tree of typed nodes implementing the `Node` interface ([`ast.go`](file:///home/ahmed/golang/redis-clone/internals/parser/ast.go)):

| AST Node Type | Code Struct | Description | Example |
| :--- | :--- | :--- | :--- |
| **`CommandExpr`** | [`*parser.CommandExpr`](file:///home/ahmed/golang/redis-clone/internals/parser/ast.go#L40-L43) | Command invocation node | `SET key val` |
| **`StringLiteral`** | [`*parser.StringLiteral`](file:///home/ahmed/golang/redis-clone/internals/parser/ast.go#L12-L14) | Text or identifier literal | `"Hello World"` |
| **`NumberLiteral`** | [`*parser.NumberLiteral`](file:///home/ahmed/golang/redis-clone/internals/parser/ast.go#L20-L22) | Numeric literal | `42` |
| **`ListLiteral`** | [`*parser.ListLiteral`](file:///home/ahmed/golang/redis-clone/internals/parser/ast.go#L28-L30) | Array list literal | `[item1 item2]` |

### Visual AST Representation

For the RQL query: `SET server_config [ "production" (GET port_num) 8080 ]`

```mermaid
graph TD
    Root["CommandExpr: SET"]
    Arg1["StringLiteral: server_config"]
    List["ListLiteral: [ ... ]"]
    
    Item1["StringLiteral: production"]
    SubCmd["CommandExpr: GET"]
    SubArg["StringLiteral: port_num"]
    Item3["NumberLiteral: 8080"]

    Root --> Arg1
    Root --> List
    List --> Item1
    List --> SubCmd
    List --> Item3
    SubCmd --> SubArg
```

---

## ⚙️ 3. Finite State Machines (FSM)

### Lexer FSM ([`lexer.go`](file:///home/ahmed/golang/redis-clone/internals/parser/lexer.go))
Scans character streams (`ch byte`) into tokens (`[]Token`):

```mermaid
stateDiagram-v2
    [*] --> ReadNextChar
    ReadNextChar --> SkipWhitespace : Space / Tab / Newline
    SkipWhitespace --> ReadNextChar
    ReadNextChar --> ReadIdent : Letter, _, :, *, -
    ReadIdent --> EmitTokenIdent : Non-ident char
    ReadNextChar --> ReadNumber : Digit 0-9
    ReadNumber --> EmitTokenNumber : Non-digit char
    ReadNextChar --> ReadString : Quote " or '
    ReadString --> EmitTokenString : Closing quote
    ReadNextChar --> ReadParen : ( or )
    ReadParen --> EmitTokenParen
    ReadNextChar --> ReadBracket : [ or ]
    ReadBracket --> EmitTokenBracket
    ReadNextChar --> EmitEOF : EOF / \0
    EmitTokenIdent --> ReadNextChar
    EmitTokenNumber --> ReadNextChar
    EmitTokenString --> ReadNextChar
    EmitTokenParen --> ReadNextChar
    EmitTokenBracket --> ReadNextChar
    EmitEOF --> [*]
```

---

### Parser FSM ([`parser.go`](file:///home/ahmed/golang/redis-clone/internals/parser/parser.go))
Consumes tokens to construct AST trees:

```mermaid
stateDiagram-v2
    [*] --> CheckParen
    CheckParen --> ParseCommandName : Consume LPAREN if present
    ParseCommandName --> ExpressionLoop : Valid Command IDENT
    ExpressionLoop --> ParseLiteral : String / Number / Ident
    ParseLiteral --> ExpressionLoop
    ExpressionLoop --> RecurseCommand : LPAREN (
    RecurseCommand --> ExpressionLoop : Returns Child CommandExpr
    ExpressionLoop --> ParseList : Open Bracket [
    ParseList --> ExpressionLoop : Returns ListLiteral
    ExpressionLoop --> CompleteAST : RPAREN ) or EOF
    CompleteAST --> [*]
```

---

## 📋 4. Supported Commands Quick Reference

| Command | Syntax | Description |
| :--- | :--- | :--- |
| `PING` | `PING [message]` | Returns `PONG` or input message |
| `ECHO` | `ECHO <message>` | Echoes input back to client |
| `SET` | `SET <key> <value> [ttl_seconds]` | Sets key value with optional TTL |
| `GET` | `GET <key>` | Gets value for key |
| `DEL` | `DEL <key1> [key2 ...]` | Deletes one or more keys |
| `EXISTS` | `EXISTS <key1> [key2 ...]` | Checks existence of keys |
| `KEYS` | `KEYS` | Lists active database keys |
| `EXPIRE` | `EXPIRE <key> <seconds>` | Sets expiration TTL on key |
| `TTL` | `TTL <key>` | Returns remaining TTL seconds |
| `INCR` | `INCR <key>` | Increments numeric value by 1 |
| `DECR` | `DECR <key>` | Decrements numeric value by 1 |
| `MGET` | `MGET <key1> <key2> ...` | Retrieves multiple keys |
| `MSET` | `MSET <k1> <v1> <k2> <v2> ...` | Sets multiple key-value pairs |
| `VALUEBYINDEX` | `VALUEBYINDEX <key> <index>` | Retrieves element at the specified 0-based index from stored `StringLiteral` or `ListLiteral` |
| `UPDATEINDEX` | `UPDATEINDEX <key> <index> <value>` | Updates element or character at the specified 0-based index with the new value |
