package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttaatoo/sqlgen/internal/generator"
	"github.com/ttaatoo/sqlgen/internal/schema"
)

func confirmOverwrite(filename string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("File %s already exists. Overwrite? [y/N]: ", filename)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: sqlgen [options]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  sqlgen -U root -p secret -db myapp -o ./models\n")
	fmt.Fprintf(os.Stderr, "  sqlgen -U root -p secret -db myapp -table users -o ./models\n")
	fmt.Fprintf(os.Stderr, "  sqlgen -H 192.168.1.100 -P 3306 -U admin -p pass -db myapp -o ./models -f\n")
}

func main() {
	var (
		host     string
		port     int
		user     string
		password string
		database string
		table    string
		output   string
		force    bool
	)

	flag.StringVar(&host, "H", "localhost", "MySQL host")
	flag.IntVar(&port, "P", 3306, "MySQL port")
	flag.StringVar(&user, "U", "root", "MySQL user")
	flag.StringVar(&password, "p", "", "MySQL password")
	flag.StringVar(&database, "db", "", "MySQL database name (required)")
	flag.StringVar(&table, "table", "", "Table name (optional, generates all tables if empty)")
	flag.StringVar(&output, "o", "", "Output directory (required)")
	flag.BoolVar(&force, "f", false, "Force overwrite existing files without confirmation")

	flag.Usage = printUsage
	flag.Parse()

	if database == "" {
		fmt.Fprintln(os.Stderr, "Error: -db is required")
		flag.Usage()
		os.Exit(1)
	}

	if output == "" {
		fmt.Fprintln(os.Stderr, "Error: -o is required")
		flag.Usage()
		os.Exit(1)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user, password, host, port, database)

	reader, err := schema.NewReader(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	var tables []string
	if table != "" {
		tables = []string{table}
	} else {
		tables, err = reader.GetTables(database)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting tables: %v\n", err)
			os.Exit(1)
		}
	}

	pkg := filepath.Base(output)
	gen := generator.New(pkg, output,
		generator.WithForce(force),
		generator.WithConfirmFunc(confirmOverwrite),
	)

	for _, tableName := range tables {
		tableSchema, err := reader.GetTableSchema(database, tableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting schema for table %s: %v\n", tableName, err)
			continue
		}

		if err := gen.Generate(tableSchema); err != nil {
			if errors.Is(err, generator.ErrSkipped) {
				fmt.Printf("Skipped: %s\n", tableName)
				continue
			}
			fmt.Fprintf(os.Stderr, "Error generating struct for table %s: %v\n", tableName, err)
			continue
		}

		fmt.Printf("Generated: %s\n", tableName)
	}

	fmt.Println("Done!")
}
