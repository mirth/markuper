package internal

import (
	"encoding/json"
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

func (db *DB) GetProject(pID ProjectID) (Project, error) {
	proj := Project{}
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Projects)
		pIDBin := []byte(pID)
		projJson := b.Get(pIDBin)

		if projJson == nil {
			return errors.New("No project exists for [" + pID + "]")
		}

		err := json.Unmarshal(projJson, &proj)

		return errors.WithStack(err)
	})

	if err != nil {
		return Project{}, err
	}

	return proj, nil
}

func (db *DB) GetSample(sID SampleID) ([]byte, error) {
	sIDBin := []byte(sID)

	var sample []byte

	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Samples)
		sample = b.Get(sIDBin)
		if sample == nil {
			return errors.New(fmt.Sprintf(
				"No sample exists for [%s]", sID),
			)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (db *DB) GetMarkup(sID SampleID) (*SampleMarkup, error) {
	sIDbin := []byte(sID)

	var smBin []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
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
	err = json.Unmarshal(smBin, &sm)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &sm, nil
}

func (db *DB) PutOne(bucket, key string, value json.RawMessage) error {
	keyBin := []byte(key)

	err := db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New(fmt.Sprintf("No such bucket [%s]", bucket))
		}

		err := b.Put(keyBin, value)

		return errors.WithStack(err)
	})

	if err != nil {
		return err
	}

	return nil
}

type KeyValue struct {
	Key   string
	Value json.RawMessage
}

func (db *DB) PutMany(bucket string, pairs []KeyValue) error {
	err := db.DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New(fmt.Sprintf("No such bucket [%s]", bucket))
		}

		for _, pair := range pairs {
			keyBin := []byte(pair.Key)

			err := b.Put(keyBin, pair.Value)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})

	return err
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

func inefficientProjectMeta(db *DB, projID ProjectID) (ProjectMeta, error) {
	doneIDs, err := getAllSampleIDsForProject(db, "markups", projID)
	if err != nil {
		return ProjectMeta{}, err
	}

	allIDs, err := getAllSampleIDsForProject(db, "samples", projID)
	if err != nil {
		return ProjectMeta{}, err
	}

	projMeta := ProjectMeta{
		TotalNumberOfSamples:    len(allIDs),
		AssessedNumberOfSamples: len(doneIDs),
	}

	return projMeta, nil
}
