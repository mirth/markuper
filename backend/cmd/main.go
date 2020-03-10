package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"backend/internal"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"backend/pkg/httpjsondecoder"
	"backend/pkg/utils"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func MakeHTTPRequestDecoder(payloadMaker func() interface{}) httptransport.DecodeRequestFunc {
	decoder := httpjsondecoder.NewDecoder()

	return func(_ context.Context, req *http.Request) (request interface{}, err error) {
		payload := payloadMaker()
		err = decoder.Decode(req, payload)

		return payload, errors.WithStack(err)
	}
}

func streamFile(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	now := utils.NowUTC()
	resp := response.(internal.ExportResponse)
	content := resp.CSV
	reader := bytes.NewReader(content)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resp.Filename))
	http.ServeContent(w, resp.R, ".csv", now, reader)

	return nil
}

var withProjectIDRequestDecoder = MakeHTTPRequestDecoder(func() interface{} {
	return &internal.WithProjectIDRequest{}
})

func withRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	payload, err := withProjectIDRequestDecoder(ctx, req)
	if err != nil {
		return nil, err
	}

	return internal.WithHttpRequest{
		R:       req,
		Payload: *payload.(*internal.WithProjectIDRequest),
	}, nil
}

func main() {
	db, err := internal.OpenDB(os.Getenv("ENV") == "test")
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	ms := internal.NewMarkupService(db)
	ps := internal.NewProjectService(db)
	pts := internal.NewTemplateService()
	e := internal.NewExporterService(db)

	nextHandler := httptransport.NewServer(
		internal.NextSampleEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	assessHandler := httptransport.NewServer(
		internal.AssessEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.AssessRequest{}
		}),
		encodeResponse,
	)

	listMarkupHandler := httptransport.NewServer(
		internal.ListMarkupEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
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
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	listTemplatesEndpoint := httptransport.NewServer(
		internal.ListTemplatesEndpoint(pts),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	exportToCsvEnpoint := httptransport.NewServer(
		internal.ExportEndpoint(e),
		withRequest,
		streamFile,
	)

	getSampleEndpoint := httptransport.NewServer(
		internal.GetSampleEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.SampleID{}
		}),
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Handle("/api/v1/project", createProjectHandler).Methods("POST")
	r.Handle("/api/v1/projects", listProjectsHandler).Methods("GET")
	r.Handle("/api/v1/project/{project_id}/next", nextHandler)
	r.Handle("/api/v1/project/{project_id}", getProjectEndpoint).Methods("GET")
	r.Handle("/api/v1/project/{project_id}/assess", assessHandler).Methods("POST")
	r.Handle("/api/v1/project/{project_id}/assessed", listMarkupHandler).Methods("GET")
	r.Handle("/api/v1/project/{project_id}/export", exportToCsvEnpoint)
	r.Handle("/api/v1/project/{project_id}/samples/{sample_id}", getSampleEndpoint)
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
