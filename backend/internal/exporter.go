package internal

import (
	"bytes"
	"context"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

type ExporterService interface {
	Export(req WithHttpRequest) (ExportResponse, error)
}

type ExportResponse struct {
	R   *http.Request
	CSV []byte
}

type ExporterServiceImpl struct {
	db *DB
}

type WithHttpRequest struct {
	R *http.Request
}

func (s *ExporterServiceImpl) Export(req WithHttpRequest) (ExportResponse, error) {
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
		CSV: buf.Bytes(),
		R:   req.R,
	}, nil
}

func ExportEndpoint(s ExporterService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		hack := request.(WithHttpRequest)

		return s.Export(hack)
	}
}

func NewExporterService(db *DB) ExporterService {
	return &ExporterServiceImpl{
		db: db,
	}
}
