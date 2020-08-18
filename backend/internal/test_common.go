package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

func openTestDB() *DB {
	db, err := OpenDB(true)
	if err != nil {
		panic(err)
	}

	return db
}

func testCloseAndReset(db *DB) {
	db.DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(Projects)
		tx.DeleteBucket(Samples)
		tx.DeleteBucket(Markups)

		return nil
	})
	db.DB.Close()
}

func testGetBucketSize(db *DB, bucket string) int {
	size := 0
	db.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		size = b.Stats().KeyN

		return nil
	})

	return size
}

func fillDirWithSamples(tmpDir, sampleExt string, nSamples int) []string {
	joinTmp := func(fn string) string {
		return filepath.Join(tmpDir, fn)
	}
	imgPaths := []string{}
	for i := 0; i < nSamples; i++ {
		path := joinTmp(fmt.Sprintf("img%d.%s", i, sampleExt))
		os.Create(path)
		imgPaths = append(imgPaths, path)
	}

	return imgPaths
}

func generateFiveSamples() []SampleData {
	samples := make([]SampleData, 0)
	for i := 0; i < 5; i++ {
		samples = append(samples, MediaSample{
			MediaURI:  fmt.Sprintf("sampleuri%d", i),
			MediaType: IMAGE_FILE_TYPE,
		})
	}

	return samples
}
