package schema

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Column struct {
	Name       string
	DataType   string
	IsNullable bool
	ColumnKey  string
	Extra      string
	Comment    string
}

type Table struct {
	Name    string
	Columns []Column
}

type Reader struct {
	db *sql.DB
}

func NewReader(dsn string) (*Reader, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Reader{db: db}, nil
}

func (r *Reader) Close() error {
	return r.db.Close()
}

func (r *Reader) GetTables(database string) ([]string, error) {
	query := `SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?`
	rows, err := r.db.Query(query, database)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func (r *Reader) GetTableSchema(database, tableName string) (*Table, error) {
	query := `
		SELECT
			COLUMN_NAME,
			DATA_TYPE,
			IS_NULLABLE,
			IFNULL(COLUMN_KEY, ''),
			IFNULL(EXTRA, ''),
			IFNULL(COLUMN_COMMENT, '')
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`
	rows, err := r.db.Query(query, database, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	table := &Table{Name: tableName}
	for rows.Next() {
		var col Column
		var isNullable string
		if err := rows.Scan(&col.Name, &col.DataType, &isNullable, &col.ColumnKey, &col.Extra, &col.Comment); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		col.IsNullable = isNullable == "YES"
		table.Columns = append(table.Columns, col)
	}
	return table, rows.Err()
}
