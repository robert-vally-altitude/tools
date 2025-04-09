package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   string
	Port     string
	Username string
	Password string
	DBName   string
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", 
		c.Username, c.Password, c.Server, c.Port, c.DBName)
}

func (c *Config) Validate() error {
	var missing []string
	if c.Server == "" {
		missing = append(missing, "server")
	}
	if c.Port == "" {
		missing = append(missing, "port")
	}
	if c.Username == "" {
		missing = append(missing, "username")
	}
	if c.Password == "" {
		missing = append(missing, "password")
	}
	if c.DBName == "" {
		missing = append(missing, "database name")
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("missing required parameters: %s", strings.Join(missing, ", "))
	}
	return nil
}

func loadConfig(envFile string) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil && envFile != ".env" {
		return nil, fmt.Errorf("could not load env file %s: %v", envFile, err)
	}

	config := &Config{
		Server:   os.Getenv("DB_SERVER"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	return config, nil
}

func loadQueryFromFile(queryName string) (string, error) {
	filename := filepath.Join("queries", strings.ReplaceAll(queryName, "-", "_")+".sql")
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read query file %s: %v", filename, err)
	}
	return string(content), nil
}

func getAvailableQueries() ([]string, error) {
	files, err := os.ReadDir("queries")
	if err != nil {
		return nil, fmt.Errorf("failed to read queries directory: %v", err)
	}

	var queries []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			queryName := strings.TrimSuffix(file.Name(), ".sql")
			queryName = strings.ReplaceAll(queryName, "_", "-")
			queries = append(queries, queryName)
		}
	}
	return queries, nil
}

func printUsage() {
	queries, err := getAvailableQueries()
	if err != nil {
		log.Printf("Warning: Failed to list available queries: %v", err)
	}

	fmt.Println("Usage: go run queries.go <query> [-f=<env_file>] [-server=<host> -port=<port> -username=<username> -password=<password> -dbname=<dbname>] [param1=value1 param2=value2 ...]")
	fmt.Println("")
	fmt.Println("Either provide a valid .env file using -f flag (defaults to .env)")
	fmt.Println("or specify all connection parameters via command line flags.")
	fmt.Println("")
	fmt.Println("Parameters:")
	fmt.Println("  Query parameters should be passed as key=value pairs.")
	fmt.Println("  For example:")
	fmt.Println("    go run queries.go -f .env.local is-verified email=user@example.com")
	fmt.Println("    go run queries.go -f .env.prod new-user-entries date=2024-03-20")
	fmt.Println("")
	fmt.Println("Available queries:")
	for _, query := range queries {
		fmt.Printf("  %s\n", query)
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	envFile := ".env" 
	for i, field := range os.Args {
		if field == "-f" {
			envFile = os.Args[i+1]
			os.Args = append(os.Args[:i], os.Args[i+2:]...)
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	queryName := os.Args[1]
	sqlQuery, err := loadQueryFromFile(queryName)
	if err != nil {
		log.Fatalf("Error: %v", err)
		printUsage()
		os.Exit(1)
	}

	server := flag.String("server", "", "Database server address")
	port := flag.String("port", "", "Database port")
	username := flag.String("username", "", "Database username")
	password := flag.String("password", "", "Database password")
	dbname := flag.String("dbname", "", "Database name")
	
	flag.CommandLine.Parse(os.Args[2:])

	params := make(map[string]string)
	for _, arg := range flag.Args() {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}

	config, err := loadConfig(envFile)
	if err != nil {
		log.Printf("Warning: %v", err)
	}

	if *server != "" {
		config.Server = *server
	}
	if *port != "" {
		config.Port = *port
	}
	if *username != "" {
		config.Username = *username
	}
	if *password != "" {
		config.Password = *password
	}
	if *dbname != "" {
		config.DBName = *dbname
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %v\n\nPlease either:\n1. Create a .env file\n2. Specify an environment file with -f\n3. Provide all connection parameters via command line flags", err)
	}

	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	for paramName, _ := range params {
		placeholder := fmt.Sprintf("?{%s}", paramName)
		sqlQuery = strings.ReplaceAll(sqlQuery, placeholder, "?")
	}

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatalf("Failed to prepare query: %v", err)
	}
	defer stmt.Close()

	paramValues := make([]interface{}, len(params))
	i := 0
	for _, value := range params {
		paramValues[i] = value
		i++
	}

	rows, err := stmt.Query(paramValues...)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	num_cols := len(columns)
	for i, col := range columns {
		fmt.Print(col)
		if i < num_cols-1 {
			fmt.Printf(",")
		}
	}
	fmt.Println("")

	values := make([]interface{}, num_cols)
	valuePtrs := make([]interface{}, num_cols)
	for i := range num_cols {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		for i, val := range values {
			switch val.(type) {
			case nil:
				fmt.Printf("")
			case []byte:
				fmt.Printf("%s", val)
			default:
				fmt.Printf("%v", val)
			}

			if i < num_cols-1 {
				fmt.Printf(",")
			}
		}
		fmt.Println("")
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error during row iteration: %v", err)
	}
}

// vim: ts=4 sts=4 sw=4 et
