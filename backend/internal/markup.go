package internal

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/go-kit/kit/endpoint"
	"github.com/recoilme/pudge"
)

var SamplesDB string = "../bin/samples"

type SampleKey struct {
	ProjectID string
	SampleID  int64
}

type URLListResponse struct {
	Urls []string `json:"urls"`
}

type SampleResponse struct {
	SampleURI string `json:"sample_uri"`
}

type TestService interface {
	GetURLS() (URLListResponse, error)
	GetNext() (SampleResponse, error)
	// Assess(SampleKey) error
}

type TestServiceImpl struct {
}

func MakeNextSampleEndpoint(s TestService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		s, _ := s.GetNext()

		return s, nil
	}
}

func MakeTestEndpoint(s TestService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		urlList, _ := s.GetURLS()

		return urlList, nil
	}
}

func decodeKey(raw []byte) SampleKey {
	buf := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buf)
	key := SampleKey{}
	_ = dec.Decode(&key)

	return key
}

func (s *TestServiceImpl) GetURLS() (URLListResponse, error) {
	rawKeys, err := pudge.Keys(SamplesDB, SampleKey{}, 0, 0, true)

	keys := make([]SampleKey, 0)
	for _, rawKey := range rawKeys {
		key := decodeKey(rawKey)
		keys = append(keys, key)
	}

	samples := make([]string, 0)
	for _, key := range keys {
		sample := ""
		_ = pudge.Get(SamplesDB, key, &sample)
		samples = append(samples, sample)
	}

	return URLListResponse{
		Urls: samples,
	}, err
}

var offset = 0

func (s *TestServiceImpl) GetNext() (SampleResponse, error) {
	rawKeys, err := pudge.Keys(SamplesDB, SampleKey{}, 1, offset, true)
	key := decodeKey(rawKeys[0])

	sampleURI := ""
	err = pudge.Get(SamplesDB, key, &sampleURI)

	offset += 1

	return SampleResponse{
		SampleURI: sampleURI,
	}, err
}
