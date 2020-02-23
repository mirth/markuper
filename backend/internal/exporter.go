package internal

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type ExporterService interface {
	Export() (ExportResponse, error)
}

type ExportResponse struct {
	CSV string
}

type ExporterServiceImpl struct {
	db *DB
}

func (s *ExporterServiceImpl) Export() (ExportResponse, error) {
	list, err := ListMarkup(s.db)
	if err != nil {
		return ExportResponse{}, nil
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	rows := [][]string{
		{"sample_id", "created_at", "markup"},
	}
	for _, entry := range list.List {
		sampleID := strconv.FormatInt(entry.SampleID.SampleID, 10)
		fmt.Println("entry.SampleMarkup.Markup: ", entry.SampleMarkup.Markup)
		row := []string{
			sampleID,
			entry.SampleMarkup.CreatedAt.Format("2006-01-02T15:04:05"),
			string(entry.SampleMarkup.Markup),
		}

		rows = append(rows, row)
	}

	err = w.WriteAll(rows)
	if err != nil {
		return ExportResponse{}, errors.WithStack(err)
	}

	return ExportResponse{
		CSV: buf.String(),
	}, nil
}
