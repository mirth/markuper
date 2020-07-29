package internal

import (
	"backend/pkg/utils"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

type ExporterService interface {
	Export(req WithHttpRequest) (ExportResponse, error)
}

type ExportResponse struct {
	R        *http.Request
	CSV      []byte
	Filename string
}

type ExporterServiceImpl struct {
	db *DB
}

type WithHttpRequest struct {
	R       *http.Request
	Payload WithProjectIDRequest
}

func takeFirstLevelGroups(j json.RawMessage, sampleColumns []string) ([]string, error) {
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(j, &objmap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	values := []string{}
	for _, c := range sampleColumns {
		value, ok := objmap[c]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Column [%s] doesn't exist in markup", c))
		}

		values = append(values, string(value))
	}

	return values, nil
}

func (s *ExporterServiceImpl) Export(req WithHttpRequest) (ExportResponse, error) {
	list, err := ListMarkup(s.db, req.Payload.ProjectID)
	if err != nil {
		return ExportResponse{}, err
	}

	p, err := s.db.GetProject(req.Payload.ProjectID)
	if err != nil {
		return ExportResponse{}, err
	}

	sampleColumns := append([]string(nil), p.Template.FieldsOrder...)

	header := []string{"sample_id", "sample_uri", "assessed_at"}
	header = append(header, sampleColumns...)
	rows := [][]string{
		header,
	}

	for _, entry := range list.List {
		row := []string{
			entry.SampleID,
			entry.SampleURI,
			entry.SampleMarkup.CreatedAt.Format("2006-01-02T15:04:05"),
		}
		sampleValues, err := takeFirstLevelGroups(entry.SampleMarkup.Markup, sampleColumns)
		if err != nil {
			return ExportResponse{}, err
		}

		row = append(row, sampleValues...)
		rows = append(rows, row)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	err = w.WriteAll(rows)
	if err != nil {
		return ExportResponse{}, errors.WithStack(err)
	}

	now := utils.NowUTC()
	projName := url.QueryEscape(p.Description.Name)
	filename := fmt.Sprintf("%s_%s.csv", projName, now.Format("2006-01-02T15:04:05"))

	return ExportResponse{
		CSV:      buf.Bytes(),
		R:        req.R,
		Filename: filename,
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
