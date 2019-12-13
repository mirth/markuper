package main

import (
	"net/http"
	"os"

	"path/filepath"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/recoilme/pudge"
)

func openDB() (*pudge.Db, error) {
	cfg := &pudge.Config{
		SyncInterval: 1}
	db, err := pudge.Open("../bin/db", cfg)

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
