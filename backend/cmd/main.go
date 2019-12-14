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

func openDB(samplesDB, markupDB, projectDB string) {
	matches, _ := filepath.Glob("/Users/tolik/Desktop/*.png")

	projectID := "project0"
	project := Project{
		ProjectID: projectID,
		Settings:  ProjectSettings{},
		State:     ProjectState{},
	}

	pudge.Set(projectDB, projectID, project)
	for i, path := range matches {
		sID := internal.SampleID{
			ProjectID: projectID,
			SampleID:  int64(i),
		}

		pudge.Set(samplesDB, sID, path)
	}
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
	openDB(samplesDB, markupDB, projectDB)

	r := mux.NewRouter()

	s := &internal.MarkupServiceImpl{
		SamplesDB: samplesDB,
		MarkupDB:  markupDB,
	}
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
