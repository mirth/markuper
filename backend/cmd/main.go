package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type URLListResponse struct {
	Urls []string `json:"urls"`
}

type TestService interface {
	GetURLS() (URLListResponse, error)
}

type testServiceImpl struct {
}

func MakeTestEndpoint(s TestService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		urlList, _ := s.GetURLS()

		return urlList, nil
	}
}

func (testServiceImpl) GetURLS() (URLListResponse, error) {
	return URLListResponse{
		Urls: []string{
			"/Users/tolik/Desktop/mk4.png",
		},
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()

	s := testServiceImpl{}
	testHandler := httptransport.NewServer(
		MakeTestEndpoint(s),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	r.Handle("/api/v1/images", testHandler)

	port := "3889"
	_ = http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
}
