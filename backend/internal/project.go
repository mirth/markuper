package internal

import (
	"backend/pkg/utils"
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type ProjectID = string

type ProjectDescription struct {
	Name string `json:"name"`
}

type Project struct {
	CreatedAt  time.Time       `json:"created_at"`
	ProjectID  ProjectID       `json:"project_id"`
	Template   ProjectTemplate `json:"template"`
	DataSource DataSource      `json:"data_source"`

	Description ProjectDescription `json:"description"`
}

type CreateProjectRequest struct {
	Template    ProjectTemplate    `json:"template"`
	DataSource  DataSource         `json:"data_source"`
	Description ProjectDescription `json:"description"`
}

func NewProject(
	template ProjectTemplate,
	dataSrc DataSource,
	desc ProjectDescription,
) Project {
	projectID := ProjectID(xid.New().String())
	now := utils.NowUTC()

	return Project{
		CreatedAt:   now,
		ProjectID:   projectID,
		Template:    template,
		DataSource:  dataSrc,
		Description: desc,
	}
}

type ProjectService interface {
	CreateProject(CreateProjectRequest) (Project, error)
	ListProjects() (ProjectList, error)
	GetProject(GetProjectRequest) (Project, error)
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

func fetchSampleList(db *DB, proj Project) error {
	fetcher := GetSampleListFetcher(proj.DataSource)

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

		err = db.Sample.Set(sID, j)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *ProjectServiceImpl) CreateProject(req CreateProjectRequest) (Project, error) {
	project := NewProject(req.Template, req.DataSource, req.Description)

	err := fetchSampleList(s.db, project)
	if err != nil {
		return Project{}, err
	}

	err = s.db.Project.Set(project.ProjectID, project)
	if err != nil {
		return Project{}, errors.WithStack(err)
	}

	return project, nil
}

type ProjectList struct {
	Projects []Project `json:"projects"`
}

func (s *ProjectServiceImpl) ListProjects() (ProjectList, error) {
	rawIDs, err := s.db.Project.Keys("", 0, 0, true)
	if err != nil {
		return ProjectList{}, errors.WithStack(err)
	}

	projects := make([]Project, 0)
	for _, rawID := range rawIDs {
		p := Project{}
		err = s.db.Project.Get(rawID, &p)
		if err != nil {
			return ProjectList{}, errors.WithStack(err)
		}
		projects = append(projects, p)
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

type GetProjectRequest struct {
	ProjectID ProjectID `json:"project_id"`
}

func (s *ProjectServiceImpl) GetProject(req GetProjectRequest) (Project, error) {
	p := Project{}
	err := s.db.Project.Get(req.ProjectID, &p)
	if err != nil {
		return Project{}, errors.WithStack(err)
	}

	return p, err
}

func GetProjectEndpoint(s ProjectService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := *request.(*GetProjectRequest)
		return s.GetProject(req)
	}
}
