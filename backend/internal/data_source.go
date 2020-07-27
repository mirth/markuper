package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type SampleURI = string

type SampleID = string

func GetProjectIDFromSampleID(id SampleID) ProjectID {
	pair := strings.Split(string(id), "-")

	return ProjectID(pair[0])
}

func NewSampleIDFor(projID ProjectID) string {
	return fmt.Sprintf("%s-%s", string(projID), xid.New().String())
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
	// fixme should check if glob
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
	case "fail_local_directory":
		return NewFailImageGlobDataSource(src.SourceURI)
	}

	return nil
}

type FailImageGlobDataSource struct {
	DataSource
}

func (s FailImageGlobDataSource) FetchSampleList() ([]SampleData, error) {
	return nil, errors.New("fail")
}

func NewFailImageGlobDataSource(sourceURI string) FailImageGlobDataSource {
	return FailImageGlobDataSource{
		DataSource: DataSource{
			Type:      "fail_local_directory",
			SourceURI: sourceURI,
		},
	}
}
