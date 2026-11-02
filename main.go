// a tool to process log files and report useful analytics
package main

import (
	"fmt"
	"github.com/lite1pal/logipunk/parser"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Print("Starting Logipunk...")

	// support only one log file for now
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("provide a path to log file")
	}

	logFilePath := args[0]

	log.Printf("Processing log file: %s", logFilePath)

	contents, err := os.ReadFile(logFilePath)

	if err != nil {
		log.Fatalf("read file: %v", err)
	}

	logRows := parseLogs(contents)

	levelCounts := make(map[parser.Level]int)

	for _, logRow := range logRows {
		levelCounts[logRow.Level]++
	}

	fmt.Println("LOGIPUNK ANALYSIS")
	fmt.Println("=================")
	fmt.Printf("File: %s\n", logFilePath)
	fmt.Printf("Total rows: %d\n", len(logRows))
	fmt.Println()
	displayLevels(logRows)
	fmt.Println()
	displayHttpStatuses(logRows)
	fmt.Println()
	displayTopEndpoints(logRows)
}

func parseLogs(contents []byte) []parser.LogRow {
	rows := strings.Split(string(contents), "\n")
	var logRows []parser.LogRow
	for _, row := range rows {
		if row == "" {
			continue
		}
		formattedRow := strings.ReplaceAll(row, "  ", " ")

		fields := strings.Split(formattedRow, " ")

		// extract createdAt field
		createdAt, err := time.Parse(time.RFC3339, fields[0])
		if err != nil {
			log.Fatalf("parse createdAt field: row: %s: %v", row, err)
		}

		// extract level field
		var level parser.Level
		switch fields[1] {
		case "INFO":
			level = parser.INFO
		case "ERROR":
			level = parser.ERROR
		case "WARN":
			level = parser.WARN
		default:
			log.Fatalf("invalid level field: %s: row: %s", level, row)
		}

		// extract ip field
		ip := fields[2]

		// extract method field
		method := fields[3]

		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
		if !slices.Contains(methods, method) {
			log.Fatalf("invalid method %s in row: %s", method, row)
		}

		// extract endpoint field
		endpoint := fields[4]

		// extract status code field
		statusCode, err := strconv.Atoi(fields[5])
		if err != nil {
			log.Fatalf("convert status code to int: row: %s: %v", row, err)
		}

		// extract duration field
		duration, err := time.ParseDuration(fields[6])
		if err != nil {
			log.Fatalf("parse duration field: row: %s: %v", row, err)
		}

		logRows = append(logRows, parser.LogRow{
			CreatedAt:  createdAt,
			Level:      level,
			IP:         ip,
			Method:     method,
			Endpoint:   endpoint,
			StatusCode: statusCode,
			Duration:   duration,
		})
	}

	return logRows
}

func displayLevels(logRows []parser.LogRow) {
	levelCounts := calculateLevels(logRows)

	if len(levelCounts) == 0 {
		return
	}

	fmt.Println("LEVELS")
	fmt.Println("======")

	for key, val := range levelCounts {
		fmt.Printf("%s: %d\n", key, val)
	}
}

func calculateLevels(logRows []parser.LogRow) map[parser.Level]int {
	levelCounts := make(map[parser.Level]int)

	for _, logRow := range logRows {
		levelCounts[logRow.Level]++
	}

	return levelCounts
}

func displayHttpStatuses(logRows []parser.LogRow) {
	httpStatusCounts := calculateHttpStatuses(logRows)

	fmt.Println("HTTP STATUS")
	fmt.Println("===========")

	for key, val := range httpStatusCounts {
		fmt.Printf("%d: %d\n", key, val)
	}
}

func calculateHttpStatuses(logRows []parser.LogRow) map[int]int {
	httpStatusCounts := make(map[int]int)

	for _, logRow := range logRows {
		httpStatusCounts[logRow.StatusCode]++
	}

	return httpStatusCounts
}

func displayTopEndpoints(logRows []parser.LogRow) {
	topEndpointCounts := calculateTopEndpoints(logRows)

	if len(topEndpointCounts) == 0 {
		return
	}

	fmt.Println("TOP ENDPOINTS")
	fmt.Println("======")

	for key, val := range topEndpointCounts {
		fmt.Printf("%s: %d\n", key, val)
	}
}

func calculateTopEndpoints(logRows []parser.LogRow) map[string]int {
	topEndpointCounts := make(map[string]int)

	for _, logRow := range logRows {
		topEndpointCounts[logRow.Endpoint]++
	}

	return topEndpointCounts
}
