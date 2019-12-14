package internal

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/recoilme/pudge"
	"github.com/stretchr/testify/assert"
)

func TestMarkupAssess(t *testing.T) {
	sampleDB, err := ioutil.TempFile("/tmp", "test123123")
	markupDB, err := ioutil.TempFile("/tmp", "test123123")

	if err != nil {
		panic(err)
	}

	{
		svc := MarkupServiceImpl{
			SamplesDB: sampleDB.Name(),
			MarkupDB:  markupDB.Name(),
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
