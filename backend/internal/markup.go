package internal

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/recoilme/pudge"
)

var SamplesDB string = "../bin/samples"
var MarkupDB string = "../bin/markup"

type SampleID struct {
	ProjectID string `json:"project_id"`
	SampleID  int64  `json:"sample_id"`
}

type SampleMarkup struct {
	CreatedAt time.Time `json:"created_at"`

	Markup json.RawMessage `json:"markup"`
}

type SampleResponse struct {
	SampleURI string `json:"sample_uri"`
}

type AssessRequest struct {
	SampleID     SampleID     `json:"sample_id"`
	SampleMarkup SampleMarkup `json:"sample_markup"`
}

type MarkupService interface {
	GetNext() (SampleResponse, error)
	Assess(AssessRequest) error
}

type MarkupServiceImpl struct {
}

func NextSampleEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		s, _ := s.GetNext()

		return s, nil
	}
}

func AssessEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r := request.(*AssessRequest)
		err := s.Assess(*r)

		s := SampleMarkup{}
		err = pudge.Get(MarkupDB, r.SampleID, &s)
		r.SampleMarkup = s

		return r, err
	}
}

func decodeKey(raw []byte) SampleID {
	buf := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buf)
	key := SampleID{}
	_ = dec.Decode(&key)

	return key
}

var offset = 0

func (s *MarkupServiceImpl) GetNext() (SampleResponse, error) {
	rawKeys, err := pudge.Keys(SamplesDB, SampleID{}, 1, offset, true)
	key := decodeKey(rawKeys[0])

	sampleURI := ""
	err = pudge.Get(SamplesDB, key, &sampleURI)

	offset += 1

	return SampleResponse{
		SampleURI: sampleURI,
	}, err
}

func (s *MarkupServiceImpl) Assess(r AssessRequest) error {
	err := pudge.Set(MarkupDB, r.SampleID, r.SampleMarkup)

	return err
}
