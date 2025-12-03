package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ttaatoo/sqlgen/internal/schema"
)

type Generator struct {
	packageName  string
	outputDir    string
	force        bool
	confirmFunc  func(filename string) bool
}

type Option func(*Generator)

func WithForce(force bool) Option {
	return func(g *Generator) {
		g.force = force
	}
}

func WithConfirmFunc(fn func(filename string) bool) Option {
	return func(g *Generator) {
		g.confirmFunc = fn
	}
}

func New(packageName, outputDir string, opts ...Option) *Generator {
	g := &Generator{
		packageName: packageName,
		outputDir:   outputDir,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// ErrSkipped is returned when user chooses to skip overwriting a file
var ErrSkipped = fmt.Errorf("skipped")

func (g *Generator) Generate(table *schema.Table) error {
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	code := g.generateStruct(table)
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return fmt.Errorf("failed to format code: %w\ngenerated code:\n%s", err, code)
	}

	filename := toSnakeCase(table.Name) + ".go"
	filePath := filepath.Join(g.outputDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists
		if !g.force && g.confirmFunc != nil {
			if !g.confirmFunc(filePath) {
				return ErrSkipped
			}
		}
	}

	if err := os.WriteFile(filePath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (g *Generator) generateStruct(table *schema.Table) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))

	imports := g.collectImports(table)
	if len(imports) > 0 {
		buf.WriteString("import (\n")
		for _, imp := range imports {
			buf.WriteString(fmt.Sprintf("\t%q\n", imp))
		}
		buf.WriteString(")\n\n")
	}

	structName := toCamelCase(table.Name)
	buf.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	for _, col := range table.Columns {
		fieldName := toCamelCase(col.Name)
		fieldType := mysqlTypeToGo(col.DataType, col.IsNullable, col.IsUnsigned)
		tag := fmt.Sprintf("`db:\"%s\"`", col.Name)

		buf.WriteString(fmt.Sprintf("\t%s %s %s\n", fieldName, fieldType, tag))
	}

	buf.WriteString("}\n")

	return buf.String()
}

func (g *Generator) collectImports(table *schema.Table) []string {
	imports := make(map[string]bool)

	for _, col := range table.Columns {
		switch col.DataType {
		case "datetime", "timestamp", "date", "time":
			imports["time"] = true
		}
	}

	var result []string
	for imp := range imports {
		result = append(result, imp)
	}
	return result
}

func mysqlTypeToGo(dataType string, nullable bool, unsigned bool) string {
	var goType string

	switch dataType {
	case "tinyint":
		if unsigned {
			goType = "uint8"
		} else {
			goType = "int8"
		}
	case "smallint":
		if unsigned {
			goType = "uint16"
		} else {
			goType = "int16"
		}
	case "mediumint", "int", "integer":
		if unsigned {
			goType = "uint32"
		} else {
			goType = "int32"
		}
	case "bigint":
		if unsigned {
			goType = "uint64"
		} else {
			goType = "int64"
		}
	case "float":
		goType = "float32"
	case "double", "real":
		goType = "float64"
	case "decimal", "numeric":
		goType = "float64"
	case "char", "varchar", "text", "tinytext", "mediumtext", "longtext", "enum", "set":
		goType = "string"
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		goType = "[]byte"
	case "datetime", "timestamp", "date", "time":
		goType = "time.Time"
	case "year":
		if unsigned {
			goType = "uint16"
		} else {
			goType = "int16"
		}
	case "bit":
		goType = "[]byte"
	case "json":
		goType = "string"
	default:
		goType = "string"
	}

	if nullable {
		switch goType {
		case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
			return "*" + goType
		case "float32", "float64":
			return "*" + goType
		case "string":
			return "*string"
		case "[]byte":
			return "[]byte"
		case "time.Time":
			return "*time.Time"
		}
	}

	return goType
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
