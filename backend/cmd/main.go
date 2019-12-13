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
)

func openDB() {
	matches, _ := filepath.Glob("/Users/tolik/Desktop/*.png")
	for i, path := range matches {
		sk := internal.SampleKey{
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

func main() {
	r := mux.NewRouter()

	openDB()

	s := &internal.TestServiceImpl{}
	urlListHandler := httptransport.NewServer(
		internal.MakeTestEndpoint(s),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	nextHandler := httptransport.NewServer(
		internal.MakeNextSampleEndpoint(s),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	r.Handle("/api/v1/images", urlListHandler)
	r.Handle("/api/v1/next", nextHandler)

	port := "3889"
	err := http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
	panic(err)
}
