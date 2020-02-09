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
	CreateProject(Project) (Project, error)
}

type ProjectServiceImpl struct {
	db *DB
}

func NewProject(db *DB) ProjectService {
	return &ProjectServiceImpl{
		db: db,
	}
}

func CreateProjectEndpoint(s ProjectService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := *request.(*Project)
		s, _ := s.CreateProject(p)

		return s, nil
	}
}

func (s *ProjectServiceImpl) CreateProject(req Project) (Project, error) {
	projectID := ProjectID(xid.New().String())
	project := Project{
		ProjectID:   projectID,
		Settings:    req.Settings,
		State:       ProjectState{},
		Description: req.Description,
	}

	err := s.db.Project.Set(projectID, project)
	if err != nil {
		return Project{}, errors.WithStack(err)
	}

	return project, nil
}
