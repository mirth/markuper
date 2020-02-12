package internal

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

type SampleURI = string

type SampleID struct {
	ProjectID ProjectID `json:"project_id"`
	SampleID  int64     `json:"sample_id"`
}

func (s SampleID) toString() string {
	return fmt.Sprintf("%s|%d", s.ProjectID, s.SampleID)
}

// type SampleData = json.RawMessage
type SampleData interface {
	JSON() ([]byte, error)
}

type ImageSample struct {
	ImageURI SampleURI `json:"image_uri"`
}

func (s ImageSample) JSON() ([]byte, error) {
	return json.Marshal(s)
}

type DataSource interface {
	FetchSamples() ([]SampleData, error)
}

type ImageGlobDataSource struct {
	globPattern string
}

func (s *ImageGlobDataSource) FetchSamples() ([]SampleData, error) {
	matches, err := filepath.Glob(s.globPattern)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	samples := []SampleData{}
	for _, path := range matches {
		samples = append(samples, ImageSample{
			ImageURI: path,
		})
	}

	return samples, nil
}
