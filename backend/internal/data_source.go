package internal

import (
	"backend/pkg/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-collections/collections/set"
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
type MediaType = string

type MediaSample struct {
	MediaURI  SampleURI `json:"media_uri"`
	MediaType MediaType `json:"media_type"`
}

func (s MediaSample) JSON() ([]byte, error) {
	return json.Marshal(s)
}

type SampleListFetcher interface {
	FetchSampleList() ([]SampleData, error)
}

type DataSource struct {
	Type      string `json:"type"`
	SourceURI string `json:"source_uri"`
}

type MediaGlobDataSource struct {
	DataSource
}

func NewMediaGlobDataSource(sourceURI string) MediaGlobDataSource {
	return MediaGlobDataSource{
		DataSource: DataSource{
			Type:      "local_directory",
			SourceURI: sourceURI,
		},
	}
}

type TestNewMediaGlobDataSource struct {
}

func (s TestNewMediaGlobDataSource) FetchSampleList() ([]SampleData, error) {
	samples := generateFiveSamples()

	return samples, nil
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

var audioFileMediaTypes = set.New(utils.ToSliceOfInterfaces([]string{
	"m4a", "flack", "mp3", "mp4", "wav", "wma", "aac", "ogg",
})...)

var imageFileMediaTypes = set.New(utils.ToSliceOfInterfaces([]string{
	"jpg", "jpeg", "png", "bmp", "gif",
})...)

var AUDIO_FILE_TYPE = "AUDIO_FILE_TYPE"
var IMAGE_FILE_TYPE = "IMAGE_FILE_TYPE"

func detectMediaDataType(sampleURI string) (string, error) {
	ext := filepath.Ext(sampleURI)
	ext = strings.ToLower(ext)[1:]

	if imageFileMediaTypes.Has(ext) {
		return IMAGE_FILE_TYPE, nil
	}

	if audioFileMediaTypes.Has(ext) {
		return AUDIO_FILE_TYPE, nil
	}

	return "", NewBusinessError(fmt.Sprintf(
		"Unsupported sample file type [%s] for %s", ext, sampleURI,
	))
}

func (s MediaGlobDataSource) FetchSampleList() ([]SampleData, error) {
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
		mediaType, err := detectMediaDataType(path)
		if err != nil {
			return nil, err
		}
		samples = append(samples, MediaSample{
			MediaURI:  path,
			MediaType: mediaType,
		})
	}

	return samples, nil
}

func GetSampleListFetcher(src DataSource) SampleListFetcher {
	switch src.Type {
	case "local_directory":
		return NewMediaGlobDataSource(src.SourceURI)
	case "fail_local_directory":
		return NewFailMediaGlobDataSource(src.SourceURI)
	case "test_local_directory":
		return TestNewMediaGlobDataSource{}
	}
	return nil
}

type FailMediaGlobDataSource struct {
	DataSource
}

func (s FailMediaGlobDataSource) FetchSampleList() ([]SampleData, error) {
	return nil, errors.New("fail")
}

func NewFailMediaGlobDataSource(sourceURI string) FailMediaGlobDataSource {
	return FailMediaGlobDataSource{
		DataSource: DataSource{
			Type:      "fail_local_directory",
			SourceURI: sourceURI,
		},
	}
}
