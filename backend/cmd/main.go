package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"backend/internal"

	"github.com/go-kit/kit/endpoint"
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

func newServer(
	e endpoint.Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc) *httptransport.Server {
	return httptransport.NewServer(e, dec, enc, httptransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) {
		var e *internal.BusingessLogicError

		if errors.As(err, &e) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, e.Error())))
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func main() {
	db, err := internal.OpenDB(os.Getenv("NODE_ENV") == "test")
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	appVersion := flag.String("appversion", "", "Set electron app version")
	flag.Parse()

	ms := internal.NewMarkupService(db)
	ps := internal.NewProjectService(db)
	pts := internal.NewTemplateService()
	e := internal.NewExporterService(db)

	nextHandler := newServer(
		internal.NextSampleEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	assessHandler := newServer(
		internal.AssessEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.AssessRequest{}
		}),
		encodeResponse,
	)

	listMarkupHandler := newServer(
		internal.ListMarkupEndpoint(ms),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	createProjectHandler := newServer(
		internal.CreateProjectEndpoint(ps),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.CreateProjectRequest{}
		}),
		encodeResponse,
	)

	listProjectsHandler := newServer(
		internal.ListProjectsEndpoint(ps),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)

	getProjectEndpoint := newServer(
		internal.GetProjectEndpoint(ps),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	getProjectMetaEndpoint := newServer(
		internal.GetProjectMetaEndpoint(ps),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	listTemplatesEndpoint := newServer(
		internal.ListTemplatesEndpoint(pts),
		MakeHTTPRequestDecoder(func() interface{} {
			return &internal.WithProjectIDRequest{}
		}),
		encodeResponse,
	)

	exportToCsvEnpoint := newServer(
		internal.ExportEndpoint(e),
		withRequest,
		streamFile,
	)

	getSampleEndpoint := newServer(
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
	r.Handle("/api/v1/project/{project_id}/stats", getProjectMetaEndpoint).Methods("GET")
	r.Handle("/api/v1/project/{project_id}/assess", assessHandler).Methods("POST")
	r.Handle("/api/v1/project/{project_id}/assessed", listMarkupHandler).Methods("GET")
	r.Handle("/api/v1/project/{project_id}/export", exportToCsvEnpoint)
	r.Handle("/api/v1/project/{project_id}/samples/{sample_id}", getSampleEndpoint)
	r.Handle("/api/v1/project_templates", listTemplatesEndpoint).Methods("GET")

	r.HandleFunc("/api/v1/healz", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("VEGETALS"))
	})
	r.HandleFunc("/api/v1/version", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(*appVersion))
	})

	port := "3889"
	err = http.ListenAndServe("127.0.0.1:"+port, handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		panic(err)
	}
}
