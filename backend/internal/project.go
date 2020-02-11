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

type ProjectSettings struct {
}

type ProjectState struct {
}

type ProjectDescription struct {
	Name string `json:"name"`
}

type Project struct {
	CreatedAt   time.Time          `json:"created_at"`
	ProjectID   ProjectID          `json:"project_id"`
	Settings    ProjectSettings    `json:"settings"`
	State       ProjectState       `json:"state"`
	Description ProjectDescription `json:"description"`
}

func NewProject(desc ProjectDescription) Project {
	projectID := ProjectID(xid.New().String())
	now := utils.NowUTC()

	return Project{
		CreatedAt:   now,
		ProjectID:   projectID,
		Settings:    ProjectSettings{},
		State:       ProjectState{},
		Description: desc,
	}
}

type ProjectService interface {
	CreateProject(ProjectDescription) (Project, error)
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
		p := *request.(*ProjectDescription)
		s, _ := s.CreateProject(p)

		return s, nil
	}
}

func (s *ProjectServiceImpl) CreateProject(req ProjectDescription) (Project, error) {
	project := NewProject(req)

	err := s.db.Project.Set(project.ProjectID, project)
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
