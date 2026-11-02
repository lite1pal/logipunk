package parser

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Level string

const (
	INFO  = "INFO"
	ERROR = "ERROR"
	WARN  = "WARN"
)

type LogRow struct {
	CreatedAt  time.Time
	Level      Level
	IP         string
	Method     string
	Endpoint   string
	StatusCode int
	Duration   time.Duration
}

func ParseLogs(contents []byte) []LogRow {
	rows := strings.Split(string(contents), "\n")
	var logRows []LogRow
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
		var level Level
		switch fields[1] {
		case "INFO":
			level = INFO
		case "ERROR":
			level = ERROR
		case "WARN":
			level = WARN
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

		logRows = append(logRows, LogRow{
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
