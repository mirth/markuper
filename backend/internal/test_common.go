package internal

import (
	"io/ioutil"
	"path"

	"github.com/recoilme/pudge"
)

func openTestDB() *DB {
	cfg := &pudge.Config{StoreMode: 2}

	tmpdir, _ := ioutil.TempDir("/tmp", "unittest")
	projectDB, _ := pudge.Open(path.Join(tmpdir, "1"), cfg)
	samplesDB, _ := pudge.Open(path.Join(tmpdir, "2"), cfg)
	markupDB, _ := pudge.Open(path.Join(tmpdir, "3"), cfg)

	return &DB{
		Project: projectDB,
		Sample:  samplesDB,
		Markup:  markupDB,
	}
}
