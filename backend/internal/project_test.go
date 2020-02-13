package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestCreateProjectRequest(name string) CreateProjectRequest {
	return CreateProjectRequest{
		Template: Template{
			Task: "classification",
		},
		DataSource: DataSource{
			Type:      "local_directory",
			SourceURI: "/tmp/*.jpg",
		},
		Description: ProjectDescription{
			Name: name,
		},
	}
}

func TestCreateProject(t *testing.T) {
	db := openTestDB()
	svc := NewProjectService(db)

	req := newTestCreateProjectRequest("testproject0")

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

	req := newTestCreateProjectRequest("testproject0")

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

	req1 := newTestCreateProjectRequest("testproject0")

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

	req2 := newTestCreateProjectRequest("testproject1")
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

func TestFetchSampleList(t *testing.T) {
	db := openTestDB()

	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)
	joinTmp := func(fn string) string {
		return filepath.Join(tmpDir, fn)
	}
	imgPaths := []string{}
	for i := 0; i < 5; i++ {
		path := joinTmp(fmt.Sprintf("img%d.jpg", i))
		os.Create(path)
		imgPaths = append(imgPaths, path)
	}

	src := DataSource{
		Type:      "local_directory",
		SourceURI: joinTmp(fmt.Sprintf("*.jpg")),
	}
	proj := NewProject(Template{}, src, ProjectDescription{})
	err := fetchSampleList(db, proj)
	assert.Nil(t, err)

	{
		sIDs, _ := getAllSampleIDs(db.Sample)

		samples := [][]byte{}
		for _, id := range sIDs {
			s := []byte{}
			db.Sample.Get(id, &s)
			samples = append(samples, s)
		}

		assert.ElementsMatch(t, [][]byte{
			[]byte(fmt.Sprintf(`{"image_uri":"%s"}`, imgPaths[0])),
			[]byte(fmt.Sprintf(`{"image_uri":"%s"}`, imgPaths[1])),
			[]byte(fmt.Sprintf(`{"image_uri":"%s"}`, imgPaths[2])),
			[]byte(fmt.Sprintf(`{"image_uri":"%s"}`, imgPaths[3])),
			[]byte(fmt.Sprintf(`{"image_uri":"%s"}`, imgPaths[4])),
		}, samples)
	}
}
