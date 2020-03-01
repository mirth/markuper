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
	Template    Template           `json:"template"`
	DataSources []DataSource       `json:"data_sources"`
	Description ProjectDescription `json:"description"`
}

func NewProject(
	template Template,
	dataSrc []DataSource,
	desc ProjectDescription,
) Project {
	projectID := ProjectID(xid.New().String())
	now := utils.NowUTC()

	return Project{
		CreatedAt:   now,
		ProjectID:   projectID,
		Template:    template,
		DataSources: dataSrc,
		Description: desc,
	}
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
		s, _ := s.CreateProject(p)

		return s, nil
	}
}

func fetchSampleList(db *DB, proj Project, src DataSource) error {
	fetcher := GetSampleListFetcher(src)

	list, err := fetcher.FetchSampleList()
	if err != nil {
		return err
	}

	for i, sample := range list {
		sID := SampleID{
			ProjectID: proj.ProjectID,
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
	project := NewProject(req.Template, req.DataSources, req.Description)

	for _, src := range project.DataSources {
		err := fetchSampleList(s.db, project, src)
		if err != nil {
			return Project{}, errors.Wrapf(err, "Failed to fetchSampleList for [%s]", src.SourceURI)
		}

		err = s.db.Put("projects", project.ProjectID, project)
		if err != nil {
			return Project{}, errors.WithStack(err)
		}
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
