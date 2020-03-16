package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageGlobDataSourceAsGlob(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	joinTmp := func(fn string) string {
		return filepath.Join(tmpDir, fn)
	}

	fillDirWithSamples(tmpDir, "jpg", 5)
	fillDirWithSamples(tmpDir, "jpeg", 5)

	src := NewImageGlobDataSource(joinTmp(fmt.Sprintf("*.jpg")))

	list, err := src.FetchSampleList()
	assert.Nil(t, err)

	{
		actual := []ImageSample{}
		for _, iterS := range list {
			actual = append(actual, iterS.(ImageSample))
		}

		assert.ElementsMatch(t, []ImageSample{
			{ImageURI: joinTmp("img0.jpg")},
			{ImageURI: joinTmp("img1.jpg")},
			{ImageURI: joinTmp("img2.jpg")},
			{ImageURI: joinTmp("img3.jpg")},
			{ImageURI: joinTmp("img4.jpg")},
		}, actual)
	}
}

func TestImageGlobDataSourceAsPath(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	joinTmp := func(fn string) string {
		return filepath.Join(tmpDir, fn)
	}

	fillDirWithSamples(tmpDir, "jpg", 3)
	fillDirWithSamples(tmpDir, "jpeg", 3)

	os.Create(joinTmp(".hidden.jpg"))
	os.Create(joinTmp("noextfile"))
	os.Mkdir(joinTmp("noextdir"), 0755)
	os.Mkdir(joinTmp("extdir.jpg"), 0755)

	var src SampleListFetcher = NewImageGlobDataSource(tmpDir)

	list, err := src.FetchSampleList()
	assert.Nil(t, err)

	{
		actual := []ImageSample{}
		for _, iterS := range list {
			actual = append(actual, iterS.(ImageSample))
		}

		assert.ElementsMatch(t, []ImageSample{
			{ImageURI: joinTmp("img0.jpg")},
			{ImageURI: joinTmp("img1.jpg")},
			{ImageURI: joinTmp("img2.jpg")},
			{ImageURI: joinTmp("img0.jpeg")},
			{ImageURI: joinTmp("img1.jpeg")},
			{ImageURI: joinTmp("img2.jpeg")},
			{ImageURI: joinTmp("extdir.jpg")}, //fixme
		}, actual)
	}
}
