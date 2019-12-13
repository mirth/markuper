package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"path/filepath"

	"backend/internal"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/recoilme/pudge"

	"backend/pkg/httpjsondecoder"
)

func openDB() {
	matches, _ := filepath.Glob("/Users/tolik/Desktop/*.png")
	for i, path := range matches {
		sk := internal.SampleID{
			ProjectID: "project0",
			SampleID:  int64(i),
		}

		pudge.Set(internal.SamplesDB, sk, path)
	}

	// return db, err
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func MakeHTTPRequestDecoder(payloadMaker func() interface{}) httptransport.DecodeRequestFunc {
	decoder := httpjsondecoder.NewDecoder()

	return func(_ context.Context, req *http.Request) (request interface{}, err error) {
		payload := payloadMaker()
		err = decoder.Decode(req, payload)

		return payload, err
	}
}

func main() {
	r := mux.NewRouter()

	openDB()

	s := &internal.MarkupServiceImpl{}
	nextHandler := httptransport.NewServer(
		internal.NextSampleEndpoint(s),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	assessHandler := httptransport.NewServer(
		internal.AssessEndpoint(s),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.AssessRequest{}
		}),
		encodeResponse,
	)

	r.Handle("/api/v1/next", nextHandler)
	r.Handle("/api/v1/assess", assessHandler).Methods("POST")

	port := "3889"
	err := http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
	panic(err)
}
