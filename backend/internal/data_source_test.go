package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newImageSample(tmpDir string, filename string) MediaSample {
	return MediaSample{
		MediaURI:  filepath.Join(tmpDir, filename),
		MediaType: IMAGE_FILE_TYPE,
	}
}

func TestMediaGlobDataSourceAsGlob(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	joinTmp := func(fn string) string {
		return filepath.Join(tmpDir, fn)
	}

	fillDirWithSamples(tmpDir, "jpg", 5)
	fillDirWithSamples(tmpDir, "jpeg", 5)

	src := NewMediaGlobDataSource(joinTmp(fmt.Sprintf("*.jpg")))

	list, err := src.FetchSampleList()
	assert.Nil(t, err)

	{
		actual := []MediaSample{}
		for _, iterS := range list {
			actual = append(actual, iterS.(MediaSample))
		}

		assert.ElementsMatch(t, []MediaSample{
			newImageSample(tmpDir, "img0.jpg"),
			newImageSample(tmpDir, "img1.jpg"),
			newImageSample(tmpDir, "img2.jpg"),
			newImageSample(tmpDir, "img3.jpg"),
			newImageSample(tmpDir, "img4.jpg"),
		}, actual)
	}
}

func TestMediaGlobDataSourceAsPath(t *testing.T) {
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

	var src SampleListFetcher = NewMediaGlobDataSource(tmpDir)

	list, err := src.FetchSampleList()
	assert.Nil(t, err)

	{
		actual := []MediaSample{}
		for _, iterS := range list {
			actual = append(actual, iterS.(MediaSample))
		}

		assert.ElementsMatch(t, []MediaSample{
			newImageSample(tmpDir, "img0.jpg"),
			newImageSample(tmpDir, "img1.jpg"),
			newImageSample(tmpDir, "img2.jpg"),
			newImageSample(tmpDir, "img0.jpeg"),
			newImageSample(tmpDir, "img1.jpeg"),
			newImageSample(tmpDir, "img2.jpeg"),
			newImageSample(tmpDir, "extdir.jpg"), //fixme
		}, actual)
	}
}
