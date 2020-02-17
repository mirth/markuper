package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/pkg/utils"
)

func fillTestDB(db *DB) Project {
	svc := NewProjectService(db)
	p, _ := svc.CreateProject(newTestCreateProjectRequest("testp0"))

	samples := make([]ImageSample, 0)
	for i := 0; i < 10; i++ {
		samples = append(samples, ImageSample{
			ImageURI: fmt.Sprintf("sampleuri%d", i),
		})
	}

	for i, sample := range samples {
		sk := SampleID{
			ProjectID: p.ProjectID,
			SampleID:  int64(i),
		}

		j, _ := json.Marshal(sample)
		db.Sample.Set(sk, j)
	}

	return p
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
			CreatedAt: utils.TestNowUTC(),
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
	proj := fillTestDB(db)

	svc := MarkupServiceImpl{
		db: db,
	}

	assertNext := func(i int64) SampleID {
		a, err := svc.GetNext()
		assert.Nil(t, err)
		sID := SampleID{
			ProjectID: proj.ProjectID,
			SampleID:  i,
		}
		e := SampleResponse{
			SampleID: sID,
			Sample:   json.RawMessage(fmt.Sprintf(`{"image_uri":"sampleuri%d"}`, i)),
			Template: proj.Template,
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

func TestListMarkup(t *testing.T) {
	db := openTestDB()
	proj := fillTestDB(db)

	svc := MarkupServiceImpl{
		db: db,
	}

	c, _ := db.Markup.Count()
	assert.Zero(t, c)

	assessSample := func(sID SampleID) {
		r := AssessRequest{
			SampleID: sID,
			SampleMarkup: SampleMarkup{
				Markup: json.RawMessage(fmt.Sprintf(`{"kek":mark%d}`, sID.SampleID)),
			},
		}
		err := svc.Assess(r)
		assert.Nil(t, err)
	}

	assessSample(SampleID{ProjectID: proj.ProjectID, SampleID: 0})
	assessSample(SampleID{ProjectID: proj.ProjectID, SampleID: 1})
	assessSample(SampleID{ProjectID: proj.ProjectID, SampleID: 2})

	c, _ = db.Markup.Count()
	assert.Equal(t, 3, c)

	list, err := svc.ListMarkup()
	assert.Nil(t, err)

	{
		assert.ElementsMatch(t, []MarkupListElement{
			{
				SampleID: SampleID{ProjectID: proj.ProjectID, SampleID: 0},
				SampleMarkup: SampleMarkup{
					CreatedAt: utils.TestNowUTC(),
					Markup:    json.RawMessage(`{"kek":mark0}`),
				},
			},
			{
				SampleID: SampleID{ProjectID: proj.ProjectID, SampleID: 1},
				SampleMarkup: SampleMarkup{
					CreatedAt: utils.TestNowUTC(),
					Markup:    json.RawMessage(`{"kek":mark1}`),
				},
			},
			{
				SampleID: SampleID{ProjectID: proj.ProjectID, SampleID: 2},
				SampleMarkup: SampleMarkup{
					CreatedAt: utils.TestNowUTC(),
					Markup:    json.RawMessage(`{"kek":mark2}`),
				},
			},
		}, list.List)
	}
}
