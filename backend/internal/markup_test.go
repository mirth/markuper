package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/recoilme/pudge"
	"github.com/stretchr/testify/assert"
)

func openTestDB() (string, string) {
	sampleDB, err := ioutil.TempFile("/tmp", "test123123")
	if err != nil {
		panic(err)
	}
	markupDB, err := ioutil.TempFile("/tmp", "test123123")
	if err != nil {
		panic(err)
	}

	return sampleDB.Name(), markupDB.Name()
}

func fillTestDB(sampleDB string) {
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

		pudge.Set(sampleDB, sk, path)
	}
}

func TestMarkupAssess(t *testing.T) {
	sampleDB, markupDB := openTestDB()

	{
		svc := MarkupServiceImpl{
			SamplesDB: sampleDB,
			MarkupDB:  markupDB,
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
			pudge.Get(svc.MarkupDB, sID, &actual)
			assert.Equal(t, mkp, actual)
		}
	}
}

func TestMarkupNext(t *testing.T) {
	sampleDB, markupDB := openTestDB()
	fillTestDB(sampleDB)

	svc := MarkupServiceImpl{
		SamplesDB: sampleDB,
		MarkupDB:  markupDB,
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
