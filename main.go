package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ttaatoo/sqlgen/internal/generator"
	"github.com/ttaatoo/sqlgen/internal/schema"
)

func main() {
	var (
		host     = flag.String("host", "localhost", "MySQL host")
		port     = flag.Int("port", 3306, "MySQL port")
		user     = flag.String("user", "root", "MySQL user")
		password = flag.String("password", "", "MySQL password")
		database = flag.String("database", "", "MySQL database name (required)")
		table    = flag.String("table", "", "Table name (optional, generates all tables if empty)")
		output   = flag.String("output", "./models", "Output directory")
		pkg      = flag.String("package", "models", "Package name for generated files")
	)

	flag.Parse()

	if *database == "" {
		fmt.Fprintln(os.Stderr, "Error: -database is required")
		flag.Usage()
		os.Exit(1)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		*user, *password, *host, *port, *database)

	reader, err := schema.NewReader(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	var tables []string
	if *table != "" {
		tables = []string{*table}
	} else {
		tables, err = reader.GetTables(*database)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting tables: %v\n", err)
			os.Exit(1)
		}
	}

	gen := generator.New(*pkg, *output)

	for _, tableName := range tables {
		tableSchema, err := reader.GetTableSchema(*database, tableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting schema for table %s: %v\n", tableName, err)
			continue
		}

		if err := gen.Generate(tableSchema); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating struct for table %s: %v\n", tableName, err)
			continue
		}

		fmt.Printf("Generated: %s\n", tableName)
	}

	fmt.Println("Done!")
}
