package schema

import (
	"testing"
)

func TestColumnStruct(t *testing.T) {
	col := Column{
		Name:       "user_id",
		DataType:   "bigint",
		IsNullable: false,
		IsUnsigned: true,
		ColumnKey:  "PRI",
		Extra:      "auto_increment",
		Comment:    "Primary key",
	}

	if col.Name != "user_id" {
		t.Errorf("Name = %q, want %q", col.Name, "user_id")
	}
	if col.DataType != "bigint" {
		t.Errorf("DataType = %q, want %q", col.DataType, "bigint")
	}
	if col.IsNullable != false {
		t.Errorf("IsNullable = %v, want %v", col.IsNullable, false)
	}
	if col.IsUnsigned != true {
		t.Errorf("IsUnsigned = %v, want %v", col.IsUnsigned, true)
	}
	if col.ColumnKey != "PRI" {
		t.Errorf("ColumnKey = %q, want %q", col.ColumnKey, "PRI")
	}
	if col.Extra != "auto_increment" {
		t.Errorf("Extra = %q, want %q", col.Extra, "auto_increment")
	}
	if col.Comment != "Primary key" {
		t.Errorf("Comment = %q, want %q", col.Comment, "Primary key")
	}
}

func TestTableStruct(t *testing.T) {
	table := Table{
		Name: "users",
		Columns: []Column{
			{Name: "id", DataType: "bigint", IsNullable: false, ColumnKey: "PRI"},
			{Name: "username", DataType: "varchar", IsNullable: false},
			{Name: "email", DataType: "varchar", IsNullable: true},
		},
	}

	if table.Name != "users" {
		t.Errorf("Name = %q, want %q", table.Name, "users")
	}
	if len(table.Columns) != 3 {
		t.Errorf("len(Columns) = %d, want %d", len(table.Columns), 3)
	}

	// Test first column
	if table.Columns[0].Name != "id" {
		t.Errorf("Columns[0].Name = %q, want %q", table.Columns[0].Name, "id")
	}
	if table.Columns[0].ColumnKey != "PRI" {
		t.Errorf("Columns[0].ColumnKey = %q, want %q", table.Columns[0].ColumnKey, "PRI")
	}

	// Test nullable column
	if table.Columns[2].IsNullable != true {
		t.Errorf("Columns[2].IsNullable = %v, want %v", table.Columns[2].IsNullable, true)
	}
}

func TestTableWithEmptyColumns(t *testing.T) {
	table := Table{
		Name:    "empty_table",
		Columns: []Column{},
	}

	if table.Name != "empty_table" {
		t.Errorf("Name = %q, want %q", table.Name, "empty_table")
	}
	if len(table.Columns) != 0 {
		t.Errorf("len(Columns) = %d, want %d", len(table.Columns), 0)
	}
}

func TestColumnNullableField(t *testing.T) {
	tests := []struct {
		name       string
		isNullable bool
	}{
		{"nullable column", true},
		{"not nullable column", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := Column{
				IsNullable: tt.isNullable,
			}

			if col.IsNullable != tt.isNullable {
				t.Errorf("IsNullable = %v, want %v", col.IsNullable, tt.isNullable)
			}
		})
	}
}

func TestColumnUnsignedField(t *testing.T) {
	tests := []struct {
		name       string
		isUnsigned bool
	}{
		{"unsigned column", true},
		{"signed column", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := Column{
				IsUnsigned: tt.isUnsigned,
			}

			if col.IsUnsigned != tt.isUnsigned {
				t.Errorf("IsUnsigned = %v, want %v", col.IsUnsigned, tt.isUnsigned)
			}
		})
	}
}

func TestColumnDataTypes(t *testing.T) {
	dataTypes := []string{
		"tinyint", "smallint", "mediumint", "int", "integer", "bigint",
		"float", "double", "real", "decimal", "numeric",
		"char", "varchar", "text", "tinytext", "mediumtext", "longtext",
		"binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob",
		"datetime", "timestamp", "date", "time", "year",
		"enum", "set", "json", "bit",
	}

	for _, dt := range dataTypes {
		t.Run(dt, func(t *testing.T) {
			col := Column{
				DataType: dt,
			}

			if col.DataType != dt {
				t.Errorf("DataType = %q, want %q", col.DataType, dt)
			}
		})
	}
}

func TestColumnKeys(t *testing.T) {
	tests := []struct {
		name      string
		columnKey string
	}{
		{"primary key", "PRI"},
		{"unique key", "UNI"},
		{"multiple key", "MUL"},
		{"no key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := Column{
				ColumnKey: tt.columnKey,
			}

			if col.ColumnKey != tt.columnKey {
				t.Errorf("ColumnKey = %q, want %q", col.ColumnKey, tt.columnKey)
			}
		})
	}
}

func TestColumnExtra(t *testing.T) {
	tests := []struct {
		name  string
		extra string
	}{
		{"auto increment", "auto_increment"},
		{"on update", "on update CURRENT_TIMESTAMP"},
		{"default generated", "DEFAULT_GENERATED"},
		{"no extra", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := Column{
				Extra: tt.extra,
			}

			if col.Extra != tt.extra {
				t.Errorf("Extra = %q, want %q", col.Extra, tt.extra)
			}
		})
	}
}
