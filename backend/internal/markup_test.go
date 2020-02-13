package internal

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func fillTestDB(db *DB) {
	samples := make([]ImageSample, 0)
	for i := 0; i < 10; i++ {
		samples = append(samples, ImageSample{
			ImageURI: fmt.Sprintf("sampleuri%d", i),
		})
	}

	projectID := "testproject0"
	for i, sample := range samples {
		sk := SampleID{
			ProjectID: projectID,
			SampleID:  int64(i),
		}

		j, _ := json.Marshal(sample)
		db.Sample.Set(sk, j)
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
			SampleID: sID,
			Sample:   json.RawMessage(fmt.Sprintf(`{"image_uri":"sampleuri%d"}`, i)),
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
