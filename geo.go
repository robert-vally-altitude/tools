package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type GeoInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	State   string `json:"state"`
}

var logRegex = regexp.MustCompile(`\[(.*?)\].*User (permitted|blocked).*?({.*})`)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: parse <logfile>")
		os.Exit(1)
	}

	logFile := os.Args[1]
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// CSV header
	writer.Write([]string{"timestamp", "ip", "status", "country", "state"})

	lineNum := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matches := logRegex.FindStringSubmatch(line)
		if len(matches) != 4 {
			continue
		}

		timestamp := matches[1]
		status := matches[2]
		jsonPart := matches[3]

		var geo GeoInfo
		if err := json.Unmarshal([]byte(jsonPart), &geo); err != nil {
			fmt.Fprintf(os.Stderr, "Line %d: failed to parse JSON: %v\n", lineNum)
			continue
		}

		if geo.Country == "" || geo.State == "" {
			//fmt.Fprintf(os.Stderr, "Line %d: missing country or state\n", lineNum)
			continue
		}

		writer.Write([]string{timestamp, geo.IP, status, geo.Country, geo.State})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

