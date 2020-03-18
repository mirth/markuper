package internal

import (
	"backend/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	bolt "go.etcd.io/bbolt"
)

type ProjectID = string

type ProjectDescription struct {
	Name string `json:"name"`
}

type Project struct {
	CreatedAt   time.Time    `json:"created_at"`
	ProjectID   ProjectID    `json:"project_id"`
	Template    Template     `json:"template"`
	DataSources []DataSource `json:"data_sources"`

	Description ProjectDescription `json:"description"`
}

type CreateProjectRequest struct {
	Template    TemplateXML        `json:"template"`
	DataSources []DataSource       `json:"data_sources"`
	Description ProjectDescription `json:"description"`
}

func NewProject(
	template TemplateXML,
	dataSrc []DataSource,
	desc ProjectDescription,
) (Project, error) {
	projectID := ProjectID(xid.New().String())
	now := utils.NowUTC()

	t, err := XMLToTemplate(template.XML)
	if err != nil {
		return Project{}, err
	}

	return Project{
		CreatedAt:   now,
		ProjectID:   projectID,
		Template:    t,
		DataSources: dataSrc,
		Description: desc,
	}, nil
}

type ProjectService interface {
	CreateProject(CreateProjectRequest) (Project, error)
	ListProjects() (ProjectList, error)
	GetProject(WithProjectIDRequest) (Project, error)
}

type ProjectServiceImpl struct {
	db *DB
}

func NewProjectService(db *DB) ProjectService {
	return &ProjectServiceImpl{
		db: db,
	}
}

func CreateProjectEndpoint(s ProjectService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := *request.(*CreateProjectRequest)
		s, err := s.CreateProject(p)

		return s, err
	}
}

func fetchSampleList(db *DB, proj Project, src DataSource) ([]Jsonable, error) {
	fetcher := GetSampleListFetcher(src)

	list, err := fetcher.FetchSampleList()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func putSamples(db *DB, projectID ProjectID, list []Jsonable) error {
	for i, sample := range list {
		sID := SampleID{
			ProjectID: projectID,
			SampleID:  int64(i), //fixme -> uuid
		}

		j, err := sample.JSON()
		if err != nil {
			return errors.WithStack(err)
		}

		err = db.Put("samples", sID, j)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *ProjectServiceImpl) CreateProject(req CreateProjectRequest) (Project, error) {
	project, err := NewProject(req.Template, req.DataSources, req.Description)
	if err != nil {
		return Project{}, err
	}

	allSamples := make([]Jsonable, 0)
	for _, src := range project.DataSources {
		list, err := fetchSampleList(s.db, project, src)
		if err != nil {
			return Project{}, errors.Wrapf(err, "Failed to fetchSampleList for [%s]", src.SourceURI)
		}

		allSamples = append(allSamples, list...)
	}

	err = putSamples(s.db, project.ProjectID, allSamples)
	if err != nil {
		return Project{}, err
	}

	err = s.db.Put("projects", project.ProjectID, project)
	if err != nil {
		return Project{}, errors.WithStack(err)
	}

	return project, nil
}

type ProjectList struct {
	Projects []Project `json:"projects"`
}

func (s *ProjectServiceImpl) ListProjects() (ProjectList, error) {
	projects := make([]Project, 0)

	err := s.db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Projects)
		err := b.ForEach(func(_k, v []byte) error {
			project := Project{}

			{
				buf := bytes.NewBuffer(v)
				dec := gob.NewDecoder(buf)
				err := dec.Decode(&project)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			projects = append(projects, project)

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return ProjectList{}, err
	}

	return ProjectList{
		Projects: projects,
	}, nil
}

func ListProjectsEndpoint(s ProjectService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return s.ListProjects()
	}
}

func (s *ProjectServiceImpl) GetProject(req WithProjectIDRequest) (Project, error) {
	return s.db.GetProject(req.ProjectID)
}

func GetProjectEndpoint(s ProjectService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := *request.(*WithProjectIDRequest)
		return s.GetProject(req)
	}
}
