# sqlgen

A fast and simple CLI tool that generates Go structs from MySQL table schemas.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/coverage-96.5%25-brightgreen.svg)]()

## Features

- Generate Go structs from MySQL table schemas
- Automatic type mapping from MySQL to Go types
- Support for `UNSIGNED` integer types (maps to Go `uint*` types)
- Support for nullable columns with pointer types
- Custom `db` tags for database field mapping
- Generate single table or all tables at once
- Clean, formatted Go code output
- Zero external dependencies (except MySQL driver)

## Installation

```bash
go install github.com/ttaatoo/sqlgen@latest
```

Or build from source:

```bash
git clone https://github.com/ttaatoo/sqlgen.git
cd sqlgen
go build -o sqlgen .
```

## Quick Start

```bash
# Generate structs for all tables in a database
sqlgen -U root -p secret -db myapp -o ./models

# Generate struct for a specific table
sqlgen -U root -p secret -db myapp -table users -o ./models
```

## Usage

```
sqlgen [options]

Options:
  -H string
        MySQL host (default "localhost")
  -P int
        MySQL port (default 3306)
  -U string
        MySQL user (default "root")
  -p string
        MySQL password
  -db string
        MySQL database name (required)
  -table string
        Table name (optional, generates all tables if empty)
  -o string
        Output directory (required)
  -f
        Force overwrite existing files without confirmation

Examples:
  sqlgen -U root -p secret -db myapp -o ./models
  sqlgen -U root -p secret -db myapp -table users -o ./models
  sqlgen -H 192.168.1.100 -P 3306 -U admin -p pass -db myapp -o ./models -f
```

## Example

Given a MySQL table:

```sql
CREATE TABLE user_accounts (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    balance DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_active TINYINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    deleted_at DATETIME
);
```

Running:

```bash
sqlgen -U root -p secret -db myapp -table user_accounts -o ./models
```

Generates `models/user_accounts.go`:

```go
package models

import (
	"time"
)

type UserAccounts struct {
	Id        uint64     `db:"id"`
	Username  string     `db:"username"`
	Email     string     `db:"email"`
	AvatarUrl *string    `db:"avatar_url"`
	Balance   float64    `db:"balance"`
	IsActive  uint8      `db:"is_active"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
```

## Type Mapping

| MySQL Type | Go Type | Nullable Go Type |
|------------|---------|------------------|
| `TINYINT` | `int8` | `*int8` |
| `TINYINT UNSIGNED` | `uint8` | `*uint8` |
| `SMALLINT` | `int16` | `*int16` |
| `SMALLINT UNSIGNED` | `uint16` | `*uint16` |
| `INT`, `MEDIUMINT` | `int32` | `*int32` |
| `INT UNSIGNED`, `MEDIUMINT UNSIGNED` | `uint32` | `*uint32` |
| `BIGINT` | `int64` | `*int64` |
| `BIGINT UNSIGNED` | `uint64` | `*uint64` |
| `FLOAT` | `float32` | `*float32` |
| `DOUBLE`, `DECIMAL` | `float64` | `*float64` |
| `CHAR`, `VARCHAR`, `TEXT` | `string` | `*string` |
| `BLOB`, `BINARY` | `[]byte` | `[]byte` |
| `DATETIME`, `TIMESTAMP`, `DATE`, `TIME` | `time.Time` | `*time.Time` |
| `JSON` | `string` | `*string` |

## Naming Conventions

- **File names**: `snake_case.go` (e.g., `user_accounts.go`)
- **Struct names**: `PascalCase` (e.g., `UserAccounts`)
- **Field names**: `PascalCase` (e.g., `AvatarUrl`)
- **DB tags**: Original column name (e.g., `` `db:"avatar_url"` ``)

## Use with Go Standard Library

The generated structs are fully compatible with Go's standard `database/sql` package:

```go
import "database/sql"

db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/myapp")

// Query single row
var user UserAccounts
row := db.QueryRow("SELECT id, username, email, created_at FROM user_accounts WHERE id = ?", 1)
err := row.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)

// Query multiple rows
rows, _ := db.Query("SELECT id, username, email FROM user_accounts")
defer rows.Close()

var users []UserAccounts
for rows.Next() {
    var u UserAccounts
    rows.Scan(&u.Id, &u.Username, &u.Email)
    users = append(users, u)
}
```

### Works with sqlx

The `db` tags are compatible with [sqlx](https://github.com/jmoiron/sqlx) for automatic struct scanning:

```go
import "github.com/jmoiron/sqlx"

var user UserAccounts
err := db.Get(&user, "SELECT * FROM user_accounts WHERE id = ?", 1)

var users []UserAccounts
err := db.Select(&users, "SELECT * FROM user_accounts")
```

## Supported Features

| Feature | Status |
|---------|--------|
| MySQL table schema reading | ✅ |
| Go struct generation | ✅ |
| Automatic MySQL → Go type mapping | ✅ |
| `UNSIGNED` integer type support | ✅ |
| Nullable column support (pointer types) | ✅ |
| `db` struct tags | ✅ |
| Single table generation | ✅ |
| Batch generation (all tables) | ✅ |
| Custom output directory | ✅ |
| Auto package name from output directory | ✅ |
| `snake_case` file naming | ✅ |
| `PascalCase` struct/field naming | ✅ |
| `go fmt` formatted output | ✅ |
| `time.Time` for datetime types | ✅ |
| `[]byte` for binary/blob types | ✅ |
| Overwrite confirmation prompt | ✅ |
| Force overwrite with `-f` flag | ✅ |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
