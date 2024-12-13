package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
)

const (
	ExitOK = iota
	ExitFailReadingHeader
	ExitFailReadingRow
	ExitFailEncodingRecords
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
}

func main() {
	os.Exit(readMain(os.Args))
}

func readMain(args []string) int {
	reader := csv.NewReader(os.Stdin)
	// TODO: Set options for no headers case
	headers, err := reader.Read()
	if err != nil {
		slog.Error("failed to read header", slog.String("error", err.Error()))
		return ExitFailReadingHeader
	}
	var records []map[string]string
	for {
		row, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			slog.Error("failed to read Row", slog.String("error", err.Error()))
			return ExitFailReadingRow
		}
		record := make(map[string]string)
		for i, header := range headers {
			record[header] = row[i]
		}
		records = append(records, record)
	}

	if err := json.NewEncoder(os.Stdout).Encode(records); err != nil {
		slog.Error("failed to encode records", slog.String("error", err.Error()))
		return ExitFailEncodingRecords
	}

	return ExitOK
}
