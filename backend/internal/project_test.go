package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	db := openTestDB()
	svc := NewProjectService(db)

	req := CreateProjectRequest{
		Template: ProjectTemplate{
			Type: "classification",
		},
		Description: ProjectDescription{
			Name: "testproject0",
		},
	}

	c, err := db.Project.Count()
	assert.Nil(t, err)
	assert.Zero(t, c)

	p, err := svc.CreateProject(req)
	{
		assert.Nil(t, err)

		c, err = db.Project.Count()
		assert.Nil(t, err)
		assert.Equal(t, 1, c)

		actual := Project{}
		db.Project.Get(p.ProjectID, &actual)
		assert.Equal(t, req.Template, actual.Template)
		assert.Equal(t, req.Description, actual.Description)
	}
}

func TestGetProject(t *testing.T) {
	db := openTestDB()
	svc := NewProjectService(db)

	req := CreateProjectRequest{
		Template: ProjectTemplate{
			Type: "classification",
		},
		Description: ProjectDescription{
			Name: "testproject0",
		},
	}

	c, err := db.Project.Count()
	assert.Nil(t, err)
	assert.Zero(t, c)

	p, err := svc.CreateProject(req)

	{
		actual, err := svc.GetProject(GetProjectRequest{
			ProjectID: p.ProjectID,
		})
		assert.Nil(t, err)

		assert.Equal(t, p.ProjectID, actual.ProjectID)
		assert.Equal(t, p.Template, actual.Template)
		assert.Equal(t, req.Template, actual.Template)
		assert.Equal(t, p.Description, actual.Description)
		assert.Equal(t, req.Description, actual.Description)
	}
}

func TestListProjects(t *testing.T) {
	db := openTestDB()
	svc := NewProjectService(db)

	req1 := CreateProjectRequest{
		Template: ProjectTemplate{
			Type: "classification",
		},
		Description: ProjectDescription{
			Name: "testproject0",
		},
	}
	c, err := db.Project.Count()
	assert.Nil(t, err)
	assert.Zero(t, c)

	_, err = svc.CreateProject(req1)
	{
		list, err := svc.ListProjects()
		assert.Nil(t, err)

		descs := []ProjectDescription{}
		for _, p := range list.Projects {
			descs = append(descs, p.Description)
		}
		assert.ElementsMatch(t, []ProjectDescription{
			{Name: req1.Description.Name},
		}, descs)
	}

	req2 := CreateProjectRequest{
		Template: ProjectTemplate{
			Type: "classification",
		},
		Description: ProjectDescription{
			Name: "testproject1",
		},
	}
	_, err = svc.CreateProject(req2)

	{
		list, err := svc.ListProjects()
		assert.Nil(t, err)

		descs := []ProjectDescription{}
		for _, p := range list.Projects {
			descs = append(descs, p.Description)
		}
		assert.ElementsMatch(t, []ProjectDescription{
			{Name: req1.Description.Name},
			{Name: req2.Description.Name},
		}, descs)
	}
}
