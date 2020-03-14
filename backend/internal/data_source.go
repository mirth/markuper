package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

type Jsonable interface {
	JSON() ([]byte, error)
}
type SampleData = Jsonable

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

func NewImageGlobDataSource(sourceURI string) ImageGlobDataSource {
	return ImageGlobDataSource{
		DataSource: DataSource{
			Type:      "local_directory",
			SourceURI: sourceURI,
		},
	}
}

func removeHiddenPaths(paths []string) []string {
	filtered := make([]string, 0)

	for _, p := range paths {
		fn := filepath.Base(p)
		if !strings.HasPrefix(fn, ".") {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

func removeWithoutExtention(paths []string) []string {
	filtered := make([]string, 0)
	for _, p := range paths {
		if len(filepath.Ext(p)) > 0 {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

func (s ImageGlobDataSource) FetchSampleList() ([]SampleData, error) {
	sourceURI := s.SourceURI
	_, err := os.Stat(s.SourceURI)
	if err == nil {
		sourceURI = filepath.Join(sourceURI, "*")
	}

	matches, err := filepath.Glob(sourceURI)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	matches = removeHiddenPaths(matches)
	matches = removeWithoutExtention(matches)
	sort.Strings(matches)

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
		return NewImageGlobDataSource(src.SourceURI)
	}

	return nil
}
