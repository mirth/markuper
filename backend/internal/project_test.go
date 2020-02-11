package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	db := openTestDB()

	svc := ProjectServiceImpl{
		db: db,
	}

	req := ProjectDescription{
		Name: "testproject0",
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
		svc.db.Project.Get(p.ProjectID, &actual)
		assert.Equal(t, req, actual.Description)
	}

}

func TestListProjects(t *testing.T) {
	db := openTestDB()

	svc := ProjectServiceImpl{
		db: db,
	}

	req1 := ProjectDescription{
		Name: "testproject0",
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
			{Name: req1.Name},
		}, descs)
	}

	req2 := ProjectDescription{
		Name: "testproject1",
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
			{Name: req1.Name},
			{Name: req2.Name},
		}, descs)
	}
}
