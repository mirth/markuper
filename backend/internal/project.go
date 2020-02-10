package internal

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type ProjectID = string

type ProjectSettings struct {
}

type ProjectState struct {
}

type ProjectDescription struct {
	Name string `json:"name"`
}

type Project struct {
	ProjectID   ProjectID          `json:"project_id"`
	Settings    ProjectSettings    `json:"settings"`
	State       ProjectState       `json:"state"`
	Description ProjectDescription `json:"description"`
}

type ProjectService interface {
	CreateProject(ProjectDescription) (Project, error)
	ListProjects() (PorjectList, error)
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
		p := *request.(*ProjectDescription)
		s, _ := s.CreateProject(p)

		return s, nil
	}
}

func (s *ProjectServiceImpl) CreateProject(req ProjectDescription) (Project, error) {
	projectID := ProjectID(xid.New().String())
	project := Project{
		ProjectID:   projectID,
		Settings:    ProjectSettings{},
		State:       ProjectState{},
		Description: req,
	}

	err := s.db.Project.Set(projectID, project)
	if err != nil {
		return Project{}, errors.WithStack(err)
	}

	return project, nil
}

type PorjectList struct {
	Projects []Project `json:"projects"`
}

func (s *ProjectServiceImpl) ListProjects() (PorjectList, error) {
	rawIDs, err := s.db.Project.Keys("", 0, 0, true)
	if err != nil {
		return PorjectList{}, errors.WithStack(err)
	}

	projects := make([]Project, 0)
	for _, rawID := range rawIDs {
		p := Project{}
		s.db.Project.Get(rawID, &p)
		projects = append(projects, p)
	}

	return PorjectList{
		Projects: projects,
	}, nil
}

func ListProjectsEndpoint(s ProjectService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return s.ListProjects()
	}
}
