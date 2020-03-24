package internal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type DB struct {
	DB *bolt.DB

	Project *bolt.Bucket
	Sample  *bolt.Bucket
	Markup  *bolt.Bucket
}

var (
	Projects []byte = []byte("projects")
	Markups  []byte = []byte("markups")
	Samples  []byte = []byte("samples")
)

func fixGob(p Project) Project {
	if p.Template.Radios == nil {
		p.Template.Radios = make([]RadioField, 0)
	}
	if p.Template.Checkboxes == nil {
		p.Template.Checkboxes = make([]CheckboxField, 0)
	}

	return p
}

func (db *DB) GetProject(pID ProjectID) (Project, error) {
	proj := Project{}
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Projects)
		pIDBin, err := encodeBin(pID)
		if err != nil {
			return err
		}
		projBin := b.Get(pIDBin)

		if projBin == nil {
			return errors.New("No project exists for [" + pID + "]")
		}

		return decodeBin(projBin).Decode(&proj)
	})

	if err != nil {
		return Project{}, err
	}

	proj = fixGob(proj)

	return proj, nil
}

func (db *DB) GetSample(sID SampleID) ([]byte, error) {
	sIDBin, err := encodeBin(sID)
	if err != nil {
		return nil, err
	}

	var sample []byte

	err = db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Samples)
		sampleBin := b.Get(sIDBin)
		if sampleBin == nil {
			return errors.New(fmt.Sprintf(
				"No sample exists for [%s, %d]",
				sID.ProjectID,
				sID.SampleID,
			))
		}

		return decodeBin(sampleBin).Decode(&sample)
	})

	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (db *DB) GetMarkup(sID SampleID) (*SampleMarkup, error) {
	sIDbin, err := encodeBin(sID)
	if err != nil {
		return nil, err
	}

	var smBin []byte
	err = db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Markups)
		smBin = b.Get(sIDbin)

		return nil
	})
	if err != nil {
		return nil, err
	}

	if smBin == nil {
		return nil, nil
	}

	sm := SampleMarkup{}
	err = decodeBin(smBin).Decode(&sm)
	if err != nil {
		return nil, err
	}

	return &sm, nil
}

func (db *DB) Put(bucket string, key, value interface{}) error {
	keyBin, err := encodeBin(key)
	if err != nil {
		return err
	}

	valueBin, err := encodeBin(value)
	if err != nil {
		return err
	}

	err = db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New(fmt.Sprintf("No such bucket [%s]", bucket))
		}

		err := b.Put(keyBin, valueBin)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func encodeBin(x interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(x)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buf.Bytes(), nil
}

func decodeBin(bin []byte) *gob.Decoder {
	buf := bytes.NewBuffer(bin)
	return gob.NewDecoder(buf)
}

func OpenDB(test bool) (*DB, error) {
	opts := &bolt.Options{
		NoSync:  test,
		Timeout: 1,
	}

	usr, err := user.Current()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	appDataDir := filepath.Join(usr.HomeDir, "Library/Application Support", "com.levchik.markuper")
	if runtime.GOOS == "windows" {
		appDataDir = filepath.Join(usr.HomeDir, os.Getenv("APPDATA"), "com.levchik.markuper")
	}

	dbFilename := filepath.Join(appDataDir, "db.db")

	if test {
		tmpdir, _ := ioutil.TempDir("", "unittest")
		dbFilename = path.Join(tmpdir, "testmarkuper.db")
	} else {
		if _, err := os.Stat(appDataDir); os.IsNotExist(err) {
			err = os.Mkdir(appDataDir, 0755)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}

	blt, err := bolt.Open(dbFilename, 0600, opts)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	db := &DB{}
	err = blt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(Projects)
		if err != nil {
			return err
		}
		db.Project = b

		b, err = tx.CreateBucketIfNotExists(Samples)
		if err != nil {
			return err
		}
		db.Sample = b

		b, err = tx.CreateBucketIfNotExists(Markups)
		if err != nil {
			return err
		}
		db.Markup = b

		return nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	db.DB = blt

	return db, nil
}
