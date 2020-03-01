package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestCreateProjectRequest(name string) CreateProjectRequest {
	return CreateProjectRequest{
		Template: DEFAULT_CLASSIFICATION_TEMPLATE,
		DataSources: []DataSource{{
			Type:      "local_directory",
			SourceURI: "/tmp/*.jpg",
		}},
		Description: ProjectDescription{
			Name: name,
		},
	}
}

func TestCreateProject(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	req := newTestCreateProjectRequest("testproject0")

	c := testGetBucketSize(db, "projects")
	assert.Zero(t, c)

	p, err := svc.CreateProject(req)
	assert.Nil(t, err)
	{
		c = testGetBucketSize(db, "projects")
		assert.Equal(t, 1, c)

		actual, err := db.GetProject(p.ProjectID)
		assert.Nil(t, err)

		assert.Equal(t, req.Template, actual.Template)
		assert.Equal(t, req.Description, actual.Description)
	}
}

func TestGetProject(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	req := newTestCreateProjectRequest("testproject0")

	c := testGetBucketSize(db, "projects")
	assert.Zero(t, c)

	p, _ := svc.CreateProject(req)

	{
		actual, err := svc.GetProject(WithProjectIDRequest{
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
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	req1 := newTestCreateProjectRequest("testproject0")

	c := testGetBucketSize(db, "projects")
	assert.Zero(t, c)

	_, _ = svc.CreateProject(req1)
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
	_, _ = svc.CreateProject(req2)

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
	defer testCloseAndReset(db)

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
	proj := NewProject(Template{}, []DataSource{src}, ProjectDescription{})
	err := fetchSampleList(db, proj, src)
	assert.Nil(t, err)

	{
		sIDs, _ := getAllSampleIDsForProject(db, "samples", proj.ProjectID)

		samples := [][]byte{}
		for _, id := range sIDs {
			s, _ := db.GetSample(id)
			samples = append(samples, s)
		}

		j := func(imgPath string) []byte {
			return []byte(fmt.Sprintf(`{"image_uri":%s}`, strconv.Quote(imgPath)))
		}

		assert.ElementsMatch(t, [][]byte{
			j(imgPaths[0]),
			j(imgPaths[1]),
			j(imgPaths[2]),
			j(imgPaths[3]),
			j(imgPaths[4]),
		}, samples)
	}
}
