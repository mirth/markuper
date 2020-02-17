package internal

import (
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

func openTestDB() *DB {
	db, _ := OpenDB(true)

	return db
}

func testCloseAndReset(db *DB) {
	db.DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("projects"))
		tx.DeleteBucket([]byte("samples"))
		tx.DeleteBucket([]byte("markups"))

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
