package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"os"

	"path/filepath"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/recoilme/pudge"
)

type URLListResponse struct {
	Urls []string `json:"urls"`
}

type SampleResponse struct {
	SampleURI string `json:"sample_uri"`
}

type TestService interface {
	GetURLS() (URLListResponse, error)
	GetNext() (SampleResponse, error)
}

type testServiceImpl struct {
	db *pudge.Db
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

func (s *testServiceImpl) GetURLS() (URLListResponse, error) {
	rawKeys, _ := s.db.Keys(SampleKey{}, 0, 0, true)

	keys := make([]SampleKey, 0)
	for _, rawKey := range rawKeys {
		key := decodeKey(rawKey)
		keys = append(keys, key)
	}

	samples := make([]string, 0)
	for _, key := range keys {
		sample := ""
		_ = s.db.Get(key, &sample)
		samples = append(samples, sample)
	}

	return URLListResponse{
		Urls: samples,
	}, nil
}

var offset = 0

func (s *testServiceImpl) GetNext() (SampleResponse, error) {
	rawKeys, err := s.db.Keys(SampleKey{}, 1, offset, true)
	key := decodeKey(rawKeys[0])

	sampleURI := ""
	err = s.db.Get(key, &sampleURI)

	offset += 1

	return SampleResponse{
		SampleURI: sampleURI,
	}, err
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type SampleKey struct {
	ProjectID string
	SampleID  int64
}

func openDB() (*pudge.Db, error) {
	cfg := &pudge.Config{
		SyncInterval: 1}
	db, err := pudge.Open("../backend/bin/db", cfg)

	matches, _ := filepath.Glob("/Users/tolik/Desktop/*.png")
	for i, path := range matches {
		sk := SampleKey{
			ProjectID: "project0",
			SampleID:  int64(i),
		}

		db.Set(sk, path)
	}

	return db, err
}

func main() {
	r := mux.NewRouter()

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	s := &testServiceImpl{
		db: db,
	}
	urlListHandler := httptransport.NewServer(
		MakeTestEndpoint(s),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	nextHandler := httptransport.NewServer(
		MakeNextSampleEndpoint(s),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	r.Handle("/api/v1/images", urlListHandler)
	r.Handle("/api/v1/next", nextHandler)

	port := "3889"
	_ = http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
}
