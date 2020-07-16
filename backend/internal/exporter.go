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
	"sort"
	"strconv"

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

func markupKeys(j json.RawMessage) ([]string, error) {
	return markupToCSVColumns(j, false)
}

func markupValues(j json.RawMessage) ([]string, error) {
	return markupToCSVColumns(j, true)
}

func rawMessageToString(m json.RawMessage) string {
	return string(m)
}

func markupToCSVColumns(j json.RawMessage, takeValues bool) ([]string, error) {
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(j, &objmap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pairs := make([][2]string, len(objmap))
	i := 0
	for key, value := range objmap {
		pairs[i][0] = key
		pairs[i][1] = rawMessageToString(value)
		i += 1
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})

	var takeValuesI int
	if takeValues {
		takeValuesI = 1
	}

	cells := make([]string, 0)
	for _, pair := range pairs {
		cells = append(cells, pair[takeValuesI])
	}

	return cells, nil
}

func (s *ExporterServiceImpl) Export(req WithHttpRequest) (ExportResponse, error) {
	list, err := ListMarkup(s.db, req.Payload.ProjectID)
	if err != nil {
		return ExportResponse{}, err
	}

	if len(list.List) == 0 {
		return ExportResponse{}, errors.New("Markup list is empty")
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	firstSample := list.List[0]
	sampleColumns, err := markupKeys(firstSample.SampleMarkup.Markup)
	if err != nil {
		return ExportResponse{}, err
	}

	header := []string{"sample_id", "sample_uri", "assessed_at"}
	header = append(header, sampleColumns...)
	rows := [][]string{
		header,
	}

	for _, entry := range list.List {
		sampleID := strconv.FormatInt(entry.SampleID.SampleID, 10)

		row := []string{
			sampleID,
			entry.SampleURI,
			entry.SampleMarkup.CreatedAt.Format("2006-01-02T15:04:05"),
		}
		sampleValues, err := markupValues(entry.SampleMarkup.Markup)
		if err != nil {
			return ExportResponse{}, err
		}

		row = append(row, sampleValues...)
		rows = append(rows, row)
	}

	err = w.WriteAll(rows)
	if err != nil {
		return ExportResponse{}, errors.WithStack(err)
	}

	p, err := s.db.GetProject(req.Payload.ProjectID)
	if err != nil {
		return ExportResponse{}, err
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
