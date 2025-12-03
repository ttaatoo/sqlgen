package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ttaatoo/sqlgen/internal/schema"
)

func TestMysqlTypeToGo(t *testing.T) {
	tests := []struct {
		name     string
		dataType string
		nullable bool
		want     string
	}{
		// Integer types
		{"tinyint", "tinyint", false, "int8"},
		{"tinyint nullable", "tinyint", true, "*int8"},
		{"smallint", "smallint", false, "int16"},
		{"smallint nullable", "smallint", true, "*int16"},
		{"mediumint", "mediumint", false, "int32"},
		{"int", "int", false, "int32"},
		{"integer", "integer", false, "int32"},
		{"int nullable", "int", true, "*int32"},
		{"bigint", "bigint", false, "int64"},
		{"bigint nullable", "bigint", true, "*int64"},

		// Float types
		{"float", "float", false, "float32"},
		{"float nullable", "float", true, "*float32"},
		{"double", "double", false, "float64"},
		{"double nullable", "double", true, "*float64"},
		{"real", "real", false, "float64"},
		{"decimal", "decimal", false, "float64"},
		{"numeric", "numeric", false, "float64"},

		// String types
		{"char", "char", false, "string"},
		{"char nullable", "char", true, "*string"},
		{"varchar", "varchar", false, "string"},
		{"varchar nullable", "varchar", true, "*string"},
		{"text", "text", false, "string"},
		{"tinytext", "tinytext", false, "string"},
		{"mediumtext", "mediumtext", false, "string"},
		{"longtext", "longtext", false, "string"},
		{"enum", "enum", false, "string"},
		{"set", "set", false, "string"},
		{"json", "json", false, "string"},

		// Binary types
		{"binary", "binary", false, "[]byte"},
		{"binary nullable", "binary", true, "[]byte"},
		{"varbinary", "varbinary", false, "[]byte"},
		{"blob", "blob", false, "[]byte"},
		{"tinyblob", "tinyblob", false, "[]byte"},
		{"mediumblob", "mediumblob", false, "[]byte"},
		{"longblob", "longblob", false, "[]byte"},
		{"bit", "bit", false, "[]byte"},

		// Time types
		{"datetime", "datetime", false, "time.Time"},
		{"datetime nullable", "datetime", true, "*time.Time"},
		{"timestamp", "timestamp", false, "time.Time"},
		{"timestamp nullable", "timestamp", true, "*time.Time"},
		{"date", "date", false, "time.Time"},
		{"time", "time", false, "time.Time"},
		{"year", "year", false, "int16"},

		// Unknown type defaults to string
		{"unknown", "unknown_type", false, "string"},
		{"unknown nullable", "unknown_type", true, "*string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mysqlTypeToGo(tt.dataType, tt.nullable)
			if got != tt.want {
				t.Errorf("mysqlTypeToGo(%q, %v) = %q, want %q", tt.dataType, tt.nullable, got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user", "User"},
		{"user_id", "UserId"},
		{"user_account", "UserAccount"},
		{"created_at", "CreatedAt"},
		{"id", "Id"},
		{"USER", "USER"},
		{"user_name_test", "UserNameTest"},
		{"", ""},
		{"a", "A"},
		{"a_b_c", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"User", "user"},
		{"UserId", "user_id"},
		{"UserAccount", "user_account"},
		{"CreatedAt", "created_at"},
		{"ID", "i_d"},
		{"user", "user"},
		{"userNameTest", "user_name_test"},
		{"", ""},
		{"A", "a"},
		{"ABC", "a_b_c"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGeneratorNew(t *testing.T) {
	gen := New("models", "/tmp/output")

	if gen.packageName != "models" {
		t.Errorf("packageName = %q, want %q", gen.packageName, "models")
	}
	if gen.outputDir != "/tmp/output" {
		t.Errorf("outputDir = %q, want %q", gen.outputDir, "/tmp/output")
	}
	if gen.force != false {
		t.Errorf("force = %v, want %v", gen.force, false)
	}
	if gen.confirmFunc != nil {
		t.Error("confirmFunc should be nil by default")
	}
}

func TestGeneratorNewWithOptions(t *testing.T) {
	confirmCalled := false
	confirmFn := func(filename string) bool {
		confirmCalled = true
		return true
	}

	gen := New("models", "/tmp/output",
		WithForce(true),
		WithConfirmFunc(confirmFn),
	)

	if gen.packageName != "models" {
		t.Errorf("packageName = %q, want %q", gen.packageName, "models")
	}
	if gen.outputDir != "/tmp/output" {
		t.Errorf("outputDir = %q, want %q", gen.outputDir, "/tmp/output")
	}
	if gen.force != true {
		t.Errorf("force = %v, want %v", gen.force, true)
	}
	if gen.confirmFunc == nil {
		t.Error("confirmFunc should not be nil")
	}

	// Test that confirmFunc works
	gen.confirmFunc("test.go")
	if !confirmCalled {
		t.Error("confirmFunc was not called")
	}
}

func TestGenerateStruct(t *testing.T) {
	gen := New("models", "/tmp/output")

	table := &schema.Table{
		Name: "user_accounts",
		Columns: []schema.Column{
			{Name: "id", DataType: "bigint", IsNullable: false, ColumnKey: "PRI"},
			{Name: "username", DataType: "varchar", IsNullable: false},
			{Name: "email", DataType: "varchar", IsNullable: true},
			{Name: "created_at", DataType: "datetime", IsNullable: false},
			{Name: "deleted_at", DataType: "datetime", IsNullable: true},
		},
	}

	code := gen.generateStruct(table)

	// Check package declaration
	if !strings.Contains(code, "package models") {
		t.Error("generated code should contain package declaration")
	}

	// Check struct name (PascalCase)
	if !strings.Contains(code, "type UserAccounts struct") {
		t.Error("generated code should contain struct declaration with PascalCase name")
	}

	// Check field names (PascalCase)
	if !strings.Contains(code, "Id int64") {
		t.Error("generated code should contain Id field")
	}
	if !strings.Contains(code, "Username string") {
		t.Error("generated code should contain Username field")
	}
	if !strings.Contains(code, "Email *string") {
		t.Error("generated code should contain Email field with pointer for nullable")
	}
	if !strings.Contains(code, "CreatedAt time.Time") {
		t.Error("generated code should contain CreatedAt field")
	}
	if !strings.Contains(code, "DeletedAt *time.Time") {
		t.Error("generated code should contain DeletedAt field with pointer for nullable")
	}

	// Check db tags
	if !strings.Contains(code, "`db:\"id\"`") {
		t.Error("generated code should contain db tag for id")
	}
	if !strings.Contains(code, "`db:\"username\"`") {
		t.Error("generated code should contain db tag for username")
	}

	// Check imports
	if !strings.Contains(code, `"time"`) {
		t.Error("generated code should import time package")
	}
}

func TestCollectImports(t *testing.T) {
	gen := New("models", "/tmp/output")

	tests := []struct {
		name        string
		table       *schema.Table
		wantImports []string
	}{
		{
			name: "no imports needed",
			table: &schema.Table{
				Columns: []schema.Column{
					{DataType: "varchar", IsNullable: false},
					{DataType: "int", IsNullable: false},
				},
			},
			wantImports: nil,
		},
		{
			name: "time import",
			table: &schema.Table{
				Columns: []schema.Column{
					{DataType: "datetime", IsNullable: false},
				},
			},
			wantImports: []string{"time"},
		},
		{
			name: "nullable datetime only needs time",
			table: &schema.Table{
				Columns: []schema.Column{
					{DataType: "datetime", IsNullable: true},
				},
			},
			wantImports: []string{"time"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.collectImports(tt.table)

			if len(tt.wantImports) == 0 && len(got) == 0 {
				return
			}

			for _, want := range tt.wantImports {
				found := false
				for _, g := range got {
					if g == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("collectImports() missing import %q, got %v", want, got)
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sqlgen_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := New("models", tmpDir)

	table := &schema.Table{
		Name: "users",
		Columns: []schema.Column{
			{Name: "id", DataType: "bigint", IsNullable: false},
			{Name: "name", DataType: "varchar", IsNullable: false},
		},
	}

	err = gen.Generate(table)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check file was created
	expectedFile := filepath.Join(tmpDir, "users.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("expected file %s was not created", expectedFile)
	}

	// Check file content is valid Go code
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	if !strings.Contains(string(content), "package models") {
		t.Error("generated file should contain package declaration")
	}
	if !strings.Contains(string(content), "type Users struct") {
		t.Error("generated file should contain struct declaration")
	}
}

func TestGenerateCreatesDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sqlgen_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputDir := filepath.Join(tmpDir, "nested", "models")
	gen := New("models", outputDir)

	table := &schema.Table{
		Name: "users",
		Columns: []schema.Column{
			{Name: "id", DataType: "int", IsNullable: false},
		},
	}

	err = gen.Generate(table)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("output directory was not created")
	}

	// Check file exists
	expectedFile := filepath.Join(outputDir, "users.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("expected file %s was not created", expectedFile)
	}
}

func TestGenerateSkipsWhenConfirmReturnsFalse(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sqlgen_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create existing file
	existingFile := filepath.Join(tmpDir, "users.go")
	if err := os.WriteFile(existingFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	gen := New("models", tmpDir,
		WithConfirmFunc(func(filename string) bool {
			return false // Don't overwrite
		}),
	)

	table := &schema.Table{
		Name: "users",
		Columns: []schema.Column{
			{Name: "id", DataType: "int", IsNullable: false},
		},
	}

	err = gen.Generate(table)
	if err != ErrSkipped {
		t.Errorf("Generate() error = %v, want ErrSkipped", err)
	}

	// Verify file was not modified
	content, _ := os.ReadFile(existingFile)
	if string(content) != "existing content" {
		t.Error("file should not have been modified")
	}
}

func TestGenerateOverwritesWhenForceIsTrue(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sqlgen_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create existing file
	existingFile := filepath.Join(tmpDir, "users.go")
	if err := os.WriteFile(existingFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	gen := New("models", tmpDir, WithForce(true))

	table := &schema.Table{
		Name: "users",
		Columns: []schema.Column{
			{Name: "id", DataType: "int", IsNullable: false},
		},
	}

	err = gen.Generate(table)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify file was overwritten
	content, _ := os.ReadFile(existingFile)
	if string(content) == "existing content" {
		t.Error("file should have been overwritten")
	}
	if !strings.Contains(string(content), "type Users struct") {
		t.Error("file should contain generated struct")
	}
}

func TestGenerateOverwritesWhenConfirmReturnsTrue(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sqlgen_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create existing file
	existingFile := filepath.Join(tmpDir, "users.go")
	if err := os.WriteFile(existingFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	confirmCalled := false
	gen := New("models", tmpDir,
		WithConfirmFunc(func(filename string) bool {
			confirmCalled = true
			return true // Allow overwrite
		}),
	)

	table := &schema.Table{
		Name: "users",
		Columns: []schema.Column{
			{Name: "id", DataType: "int", IsNullable: false},
		},
	}

	err = gen.Generate(table)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !confirmCalled {
		t.Error("confirmFunc should have been called")
	}

	// Verify file was overwritten
	content, _ := os.ReadFile(existingFile)
	if !strings.Contains(string(content), "type Users struct") {
		t.Error("file should contain generated struct")
	}
}

func TestErrSkipped(t *testing.T) {
	if ErrSkipped == nil {
		t.Error("ErrSkipped should not be nil")
	}
	if ErrSkipped.Error() != "skipped" {
		t.Errorf("ErrSkipped.Error() = %q, want %q", ErrSkipped.Error(), "skipped")
	}
}
