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

type SampleData interface {
	JSON() ([]byte, error)
}

type ImageSample struct {
	ImageURI SampleURI `json:"image_uri"`
}

func (s ImageSample) JSON() ([]byte, error) {
	return json.Marshal(s)
}

type SampleListFetcher interface {
	FetchSampleList() ([]SampleData, error)
}

type DataSource struct {
	Type      string `json:"type"`
	SourceURI string `json:"source_uri"`
}

type ImageGlobDataSource struct {
	DataSource
}

func (s ImageGlobDataSource) FetchSampleList() ([]SampleData, error) {
	matches, err := filepath.Glob(s.SourceURI)
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

func GetSampleListFetcher(src DataSource) SampleListFetcher {
	switch src.Type {
	case "local_directory":
		return ImageGlobDataSource{
			DataSource: src,
		}
	}

	return nil
}
