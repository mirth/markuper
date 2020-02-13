package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"backend/internal"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/recoilme/pudge"

	"backend/pkg/httpjsondecoder"
)

func openDB(samplesDBFile, markupDBFile, projectDBFile string) (*internal.DB, error) {
	storeMode := 0

	// TODO: it e2e!
	if os.Getenv("ENV") == "test" {
		storeMode = 2
		samplesDBFile = "/tmp/test_1"
		markupDBFile = "/tmp/test_2"
		projectDBFile = "/tmp/test_3"
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
	samplesDB := "/tmp/1"
	markupDB := "/tmp/2"
	projectDB := "/tmp/3"
	db, err := openDB(samplesDB, markupDB, projectDB)

	if err != nil {
		panic(err)
	}

	ms := internal.NewMarkupService(db)
	ps := internal.NewProjectService(db)
	pts := internal.NewTemplateService()

	nextHandler := httptransport.NewServer(
		internal.NextSampleEndpoint(ms),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	assessHandler := httptransport.NewServer(
		internal.AssessEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.AssessRequest{}
		}),
		encodeResponse,
	)

	createProjectHandler := httptransport.NewServer(
		internal.CreateProjectEndpoint(ps),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.CreateProjectRequest{}
		}),
		encodeResponse,
	)

	listProjectsHandler := httptransport.NewServer(
		internal.ListProjectsEndpoint(ps),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	getProjectEndpoint := httptransport.NewServer(
		internal.GetProjectEndpoint(ps),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.GetProjectRequest{}
		}),
		encodeResponse,
	)

	listTemplatesEndpoint := httptransport.NewServer(
		internal.ListTemplatesEndpoint(pts),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Handle("/api/v1/next", nextHandler)
	r.Handle("/api/v1/assess", assessHandler).Methods("POST")
	r.Handle("/api/v1/project", createProjectHandler).Methods("POST")
	r.Handle("/api/v1/projects", listProjectsHandler).Methods("GET")
	r.Handle("/api/v1/project/{project_id}", getProjectEndpoint).Methods("GET")
	r.Handle("/api/v1/project_templates", listTemplatesEndpoint).Methods("GET")

	r.HandleFunc("/api/v1/healz", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("VEGETALS"))
	})

	port := "3889"
	err = http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		panic(err)
	}
}
