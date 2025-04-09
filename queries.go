package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var queries = map[string]string{
	"entries": `
        SELECT 
          DATE(created_at) as date, 
          count(distinct playfab_id) as unique_entries,
          count(*) as total_entries
        FROM promo_entries
        WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
        GROUP BY DATE(created_at)
        ORDER BY date
    `,
	"entries-per-promo": `
        SELECT 
          DATE(created_at) as date, 
          count(distinct playfab_id) as unique_entries,
          count(*) as total_entries,
          promo_id 
        FROM promo_entries
        WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
        GROUP BY DATE(created_at), promo_id
        ORDER BY date
    `,
	"total-verified": `
        SELECT count(*) FROM email_verifications WHERE created_at >= "2024-09-11"
    `,
	"verified": `
        SELECT 
            DATE(a.created_at) as created_at,
            COUNT(*) as num_created,
            SUM(CASE WHEN email_verified_at IS NOT NULL THEN 1 ELSE 0 END) as num_verified
        FROM (
            SELECT playfab_id, DATE(created_at) as created_at, email_verified_at
            FROM email_verifications
            WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
            GROUP BY playfab_id, DATE(created_at)
        ) as a
        GROUP BY DATE(created_at)
        ORDER BY created_at
    `,
	"all-sweepz": `
        SELECT DATE(pe.created_at) as date, count(pe.id) as entries, p.sweepz_entry as sweepz_per_entry, count(pe.id) * p.sweepz_entry as total_sweepz, pe.promo_id, pe.playfab_id
        FROM promo_entries pe
        LEFT JOIN promos p on pe.promo_id = p.id
        WHERE pe.created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 14 DAY)
        GROUP BY DATE(pe.created_at), pe.promo_id, pe.playfab_id
        ORDER BY date, pe.playfab_id;
    `,
	"promos": `
        SELECT p.id, p.title, p.status, p.sweepz_entry as sweepz_per_entry, count(pe.id) as entries, count(distinct pe.playfab_id) as unique_entries, count(pe.id) * p.sweepz_entry as total_sweepz, p.start_at, p.end_at
        FROM promo_entries pe
        LEFT JOIN promos p on pe.promo_id = p.id
        WHERE p.start_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
        GROUP BY pe.promo_id
        ORDER BY p.id
    `,
	"new-user-entries": `
        SELECT DATE(ev.email_verified_at), ev.email, ev.playfab_id, pe.promo_id
        FROM email_verifications ev
        INNER JOIN promo_entries pe ON pe.playfab_id = ev.playfab_id
        WHERE DATE(ev.email_verified_at) = '2025-01-16';
    `,
	"orders": `
        SELECT playfab_id, pspReference, country_code, order_number, amount, game_id, playfab_item_id, display_name, created_at, updated_at
        FROM orders 
        WHERE status = 'paid'
    `,
	"is-verified": `
        SELECT playfab_id, email, email_verified_at, created_at
        FROM email_verifications
        WHERE email = 'rowlandsamantha065@gmail.com'
    `,
}

func printUsage() {
	fmt.Println("Usage: go run queries.go <query> [-f=<env_file>] [-server=<host> -port=<port> -username=<username> -password=<password> -dbname=<dbname>]")
	fmt.Println("")
	fmt.Println("Either provide a valid .env file using -f flag (defaults to .env)")
	fmt.Println("or specify all connection parameters via command line flags.")
	fmt.Println("")
	fmt.Println("Available queries:")
	for k, _ := range queries {
		fmt.Printf("  %s\n", k)
	}
}

func hasRequiredFlags(server, port, username, password, dbname string) bool {
	return server != "" && port != "" && username != "" && password != "" && dbname != ""
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	queryName := os.Args[1]
	sqlQuery, exists := queries[queryName]
	if !exists {
		log.Fatalf("Unknown query: %v", queryName)
		printUsage()
		os.Exit(1)
	}

	envFile := flag.String("f", ".env", "Path to the environment file")
	server := flag.String("server", "", "Database server address")
	port := flag.String("port", "", "Database port")
	username := flag.String("username", "", "Database username")
	password := flag.String("password", "", "Database password")
	dbname := flag.String("dbname", "", "Database name")

	flag.CommandLine.Parse(os.Args[2:])

	envFileLoaded := false
	if err := godotenv.Load(*envFile); err == nil {
		envFileLoaded = true
	} else {
		if *envFile != ".env" {
			log.Printf("Warning: Could not load specified env file %s: %v", *envFile, err)
		}
	}

	serverAddr := *server
	if serverAddr == "" {
		serverAddr = os.Getenv("DB_SERVER")
	}

	dbPort := *port
	if dbPort == "" {
		dbPort = os.Getenv("DB_PORT")
	}

	user := *username
	if user == "" {
		user = os.Getenv("DB_USERNAME")
	}

	pass := *password
	if pass == "" {
		pass = os.Getenv("DB_PASSWORD")
	}

	dbName := *dbname
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}

	if !hasRequiredFlags(serverAddr, dbPort, user, pass, dbName) {
		if !envFileLoaded {
			log.Fatalf("Error: No environment file found at %s and not all connection parameters provided via flags.\n\nPlease either:\n1. Create a .env file\n2. Specify an environment file with -f\n3. Provide all connection parameters via command line flags", *envFile)
		}
		log.Fatalf("Error: Missing required database connection parameters. Please check your environment file or provide all parameters via command line flags")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, serverAddr, dbPort, dbName)

	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	rows, err := dbConn.Query(sqlQuery)
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
