package internal

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/recoilme/pudge"
	"github.com/stretchr/testify/assert"
)

func openTestDB() *DB {
	cfg := &pudge.Config{StoreMode: 2}

	projectDB, _ := pudge.Open("/tmp/1", cfg)
	samplesDB, _ := pudge.Open("/tmp/2", cfg)
	markupDB, _ := pudge.Open("/tmp/3", cfg)

	return &DB{
		Project: projectDB,
		Sample:  samplesDB,
		Markup:  markupDB,
	}
}

func fillTestDB(db *DB) {
	matches := make([]string, 0)
	for i := 0; i < 10; i++ {
		matches = append(matches, fmt.Sprintf("sampleuri%d", i))
	}

	projectID := "testproject0"
	for i, path := range matches {
		sk := SampleID{
			ProjectID: projectID,
			SampleID:  int64(i),
		}

		db.Sample.Set(sk, path)
	}
}

func TestMarkupAssess(t *testing.T) {
	db := openTestDB()

	{
		svc := MarkupServiceImpl{
			db: db,
		}

		sID := SampleID{
			ProjectID: "project0",
			SampleID:  0,
		}
		mkp := SampleMarkup{
			CreatedAt: time.Now().UTC(),
			Markup:    json.RawMessage(`{"kek": "kek"}`),
		}
		r := AssessRequest{
			SampleID:     sID,
			SampleMarkup: mkp,
		}

		err := svc.Assess(r)

		{
			assert.Nil(t, err)

			actual := SampleMarkup{}
			svc.db.Markup.Get(sID, &actual)
			assert.Equal(t, mkp, actual)
		}
	}
}

func TestMarkupNext(t *testing.T) {
	db := openTestDB()
	fillTestDB(db)

	svc := MarkupServiceImpl{
		db: db,
	}

	assertNext := func(i int64) SampleID {
		a, err := svc.GetNext()
		assert.Nil(t, err)
		sID := SampleID{
			ProjectID: "testproject0",
			SampleID:  i,
		}
		e := SampleResponse{
			SampleID:  sID,
			SampleURI: fmt.Sprintf("sampleuri%d", i),
		}
		assert.Equal(t, e, a)

		return sID
	}

	assessSample := func(sID SampleID) {
		r := AssessRequest{SampleID: sID}
		err := svc.Assess(r)
		assert.Nil(t, err)
	}

	{
		{
			sID := assertNext(0)
			assessSample(sID)
		}
		{
			sID := assertNext(1)
			assessSample(sID)
		}
		{
			sID := assertNext(2)
			assessSample(sID)
		}
	}
}
