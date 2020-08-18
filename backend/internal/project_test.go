package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestCreateProjectRequest(name string, template TemplateXML) CreateProjectRequest {
	return CreateProjectRequest{
		Template: template,
		DataSources: []DataSource{{
			Type:      "local_directory",
			SourceURI: "/tmp/*.jpg",
		}},
		Description: ProjectDescription{
			Name: name,
		},

		ShuffleSamples: false,
	}
}

func TestCreateProject(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	req := newTestCreateProjectRequest("testproject0", DEFAULT_CLASSIFICATION_TEMPLATE)

	c := testGetBucketSize(db, "projects")
	assert.Zero(t, c)

	p, err := svc.CreateProject(req)
	assert.Nil(t, err)
	{
		c = testGetBucketSize(db, "projects")
		assert.Equal(t, 1, c)

		actual, err := db.GetProject(p.ProjectID)
		assert.Nil(t, err)

		tt, _ := XMLToTemplate(req.Template.XML)
		assert.Equal(t, tt, actual.Template)
		assert.Equal(t, req.Description, actual.Description)
	}
}

// fixme create project with empty source
func TestGetProject(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	req := newTestCreateProjectRequest("testproject0", DEFAULT_CLASSIFICATION_TEMPLATE)

	c := testGetBucketSize(db, "projects")
	assert.Zero(t, c)

	p, _ := svc.CreateProject(req)

	{
		actual, err := svc.GetProject(WithProjectIDRequest{
			ProjectID: p.ProjectID,
		})
		assert.Nil(t, err)

		tt, _ := XMLToTemplate(req.Template.XML)
		assert.Equal(t, p.ProjectID, actual.ProjectID)
		assert.Equal(t, p.Template, actual.Template)
		assert.Equal(t, tt, p.Template)
		assert.Equal(t, p.Description, actual.Description)
		assert.Equal(t, req.Description, actual.Description)
		assert.Equal(t, req.DataSources, actual.DataSources)
	}
}

func TestListProjects(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	req1 := newTestCreateProjectRequest("testproject0", DEFAULT_CLASSIFICATION_TEMPLATE)

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

	req2 := newTestCreateProjectRequest("testproject1", DEFAULT_CLASSIFICATION_TEMPLATE)
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

// fixmet test putSamples
func TestFetchSampleList(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)

	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)
	imgPaths := fillDirWithSamples(tmpDir, "jpg", 5)

	joinTmp := func(fn string) string {
		return filepath.Join(tmpDir, fn)
	}
	src := DataSource{
		Type:      "local_directory",
		SourceURI: joinTmp(fmt.Sprintf("*.jpg")),
	}

	shuffleSamples := false
	proj, _ := NewProject(
		DEFAULT_CLASSIFICATION_TEMPLATE,
		[]DataSource{src},
		ProjectDescription{},
		shuffleSamples,
	)
	list, err := fetchSampleList(db, proj, src)
	assert.Nil(t, err)

	{
		j := func(imgPath string) MediaSample {
			return MediaSample{
				MediaURI:  (imgPath),
				MediaType: IMAGE_FILE_TYPE,
			}
		}

		assert.ElementsMatch(t, []MediaSample{
			j(imgPaths[0]),
			j(imgPaths[1]),
			j(imgPaths[2]),
			j(imgPaths[3]),
			j(imgPaths[4]),
		}, list)
	}
}

func TestCreateProjectWithMultipleDataSources(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	tmpDir0, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir0)
	tmpDir1, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir1)
	tmpDir2, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir2)

	imgs0 := fillDirWithSamples(tmpDir0, "jpg", 3)
	imgs1 := fillDirWithSamples(tmpDir1, "png", 2)
	fillDirWithSamples(tmpDir2, "tiff", 2)

	req := newTestCreateProjectRequest("testproject0", DEFAULT_CLASSIFICATION_TEMPLATE)
	req.DataSources = append(
		req.DataSources,
		NewMediaGlobDataSource(filepath.Join(tmpDir0, "*.jpg")).DataSource,
		NewMediaGlobDataSource(filepath.Join(tmpDir1, "*.png")).DataSource,
		NewMediaGlobDataSource(filepath.Join(tmpDir2, "*.png")).DataSource,
	)

	c := testGetBucketSize(db, "samples")
	assert.Zero(t, c)

	p, err := svc.CreateProject(req)
	assert.Nil(t, err)

	{
		c := testGetBucketSize(db, "samples")
		assert.Equal(t, 5, c)

		samples, err := getAllSamplesForProject(db, p.ProjectID)
		assert.Nil(t, err)

		uris := make([]SampleURI, 0)
		for _, s := range samples {
			var objmap map[string]json.RawMessage
			json.Unmarshal(s, &objmap)
			uri := ""
			json.Unmarshal(objmap["media_uri"], &uri)
			uris = append(uris, uri)
		}

		assert.ElementsMatch(
			t,
			append(imgs0, imgs1...),
			uris,
		)
	}
}

func TestCreateProjectWithMultipleDataSourcesWhenSourceFail(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	svc := NewProjectService(db)

	tmpDir0, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir0)
	tmpDir1, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir1)
	tmpDir2, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir2)

	fillDirWithSamples(tmpDir0, "jpg", 3)
	fillDirWithSamples(tmpDir1, "png", 2)
	fillDirWithSamples(tmpDir2, "tiff", 2)

	req := newTestCreateProjectRequest("testproject0", DEFAULT_CLASSIFICATION_TEMPLATE)
	req.DataSources = append(
		req.DataSources,
		NewMediaGlobDataSource(filepath.Join(tmpDir0, "*.jpg")).DataSource,
		NewFailMediaGlobDataSource(filepath.Join(tmpDir1, "*.png")).DataSource,
		NewMediaGlobDataSource(filepath.Join(tmpDir2, "*.png")).DataSource,
	)

	c := testGetBucketSize(db, "projects")
	assert.Zero(t, c)

	c = testGetBucketSize(db, "samples")
	assert.Zero(t, c)

	_, err := svc.CreateProject(req)
	assert.NotNil(t, err)

	{
		c := testGetBucketSize(db, "projects")
		assert.Zero(t, c)

		c = testGetBucketSize(db, "samples")
		assert.Zero(t, c)
	}
}
