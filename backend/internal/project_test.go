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

	req := Project{
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

		actual := req.Description
		svc.db.Project.Get(p.ProjectID, &actual)
		assert.Equal(t, req.Description, actual)
	}
}
