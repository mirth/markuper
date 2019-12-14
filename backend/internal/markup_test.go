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

	{
		{
			a, err := svc.GetNext()
			assert.Nil(t, err)
			e := SampleResponse{
				SampleID: SampleID{
					ProjectID: "testproject0",
					SampleID:  0,
				},
				SampleURI: "sampleuri0",
			}
			assert.Equal(t, e, a)
		}
		{
			a, err := svc.GetNext()
			assert.Nil(t, err)
			e := SampleResponse{
				SampleID: SampleID{
					ProjectID: "testproject0",
					SampleID:  1,
				},
				SampleURI: "sampleuri1",
			}
			assert.Equal(t, e, a)
		}
		{
			a, err := svc.GetNext()
			assert.Nil(t, err)
			e := SampleResponse{
				SampleID: SampleID{
					ProjectID: "testproject0",
					SampleID:  2,
				},
				SampleURI: "sampleuri2",
			}
			assert.Equal(t, e, a)
		}
	}
}
