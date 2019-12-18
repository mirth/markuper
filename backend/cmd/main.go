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
	"github.com/pkg/errors"
	"github.com/recoilme/pudge"

	"backend/pkg/httpjsondecoder"
)

type ProjectID = string

type ProjectSettings struct {
}

type ProjectState struct {
}

type Project struct {
	ProjectID ProjectID       `json:"project_id"`
	Settings  ProjectSettings `json:"settings"`
	State     ProjectState    `json:"state"`
}

func openDB(samplesDBFile, markupDBFile, projectDBFile string) (*internal.DB, error) {
	storeMode := 0
	if os.Getenv("ENV") == "test" {
		storeMode = 2
		samplesDBFile = "/tmp/1"
		markupDBFile = "/tmp/2"
		projectDBFile = "/tmp/3"
	}

	cfg := &pudge.Config{StoreMode: storeMode}

	projectDB, err := pudge.Open(projectDBFile, cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	samplesDB, err := pudge.Open(samplesDBFile, cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	markupDB, err := pudge.Open(markupDBFile, cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	projectID := "project0"
	project := Project{
		ProjectID: projectID,
		Settings:  ProjectSettings{},
		State:     ProjectState{},
	}

	matches, _ := filepath.Glob("/Users/tolik/Desktop/*.png")

	if os.Getenv("ENV") == "test" {
		matches = []string{
			"img0",
			"img1",
			"img2",
		}
	}

	projectDB.Set(projectID, project)
	for i, path := range matches {
		sID := internal.SampleID{
			ProjectID: projectID,
			SampleID:  int64(i),
		}

		samplesDB.Set(sID, path)
	}

	return &internal.DB{
		Project: projectDB,
		Sample:  samplesDB,
		Markup:  markupDB,
	}, nil
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
	samplesDB := "../bin/samples"
	markupDB := "../bin/markup"
	projectDB := "../bin/project"
	db, err := openDB(samplesDB, markupDB, projectDB)

	if err != nil {
		panic(err)
	}

	s := internal.NewMarkupService(db)
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

	r := mux.NewRouter()
	r.Handle("/api/v1/next", nextHandler)
	r.Handle("/api/v1/assess", assessHandler).Methods("POST")

	port := "3889"
	err = http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		panic(err)
	}
}
