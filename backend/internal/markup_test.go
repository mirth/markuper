package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/pkg/utils"
)

func fillTestDB(db *DB) Project {
	return fillTestDBWithProj(db, "testp0")
}

func fillTestDBWithProj(db *DB, projName string) Project {
	svc := NewProjectService(db)
	p, _ := svc.CreateProject(newTestCreateProjectRequest(projName))

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
		db.Put("samples", sk, j)
	}

	return p
}

func AssessWithMarkup(
	t *testing.T,
	s MarkupServiceImpl,
	sID SampleID,
	markup string,
) {
	mkp := SampleMarkup{
		CreatedAt: utils.TestNowUTC(),
		Markup:    []byte(markup),
	}
	r := AssessRequest{
		SampleID: sID,
		SampleMarkup: SampleMarkup{
			CreatedAt: utils.TestNowUTC(),
			Markup:    json.RawMessage(markup),
		},
	}

	err := s.Assess(r)

	{
		assert.Nil(t, err)

		actual, _ := s.db.GetMarkup(sID)
		assert.Equal(t, mkp, actual)
	}
}

func TestMarkupAssess(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)

	{
		svc := MarkupServiceImpl{
			db: db,
		}

		AssessWithMarkup(t, svc, SampleID{ProjectID: "kek", SampleID: 0}, `{"kek": "kek"}`)
	}
}

func TestMarkupNext(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	proj := fillTestDB(db)

	svc := MarkupServiceImpl{
		db: db,
	}

	assertNext := func(i int64) SampleID {
		a, err := svc.GetNext(WithProjectIDRequest{
			ProjectID: proj.ProjectID,
		})
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
	defer testCloseAndReset(db)
	proj0 := fillTestDBWithProj(db, "proj0")
	proj1 := fillTestDBWithProj(db, "proj1")

	svc := MarkupServiceImpl{
		db: db,
	}

	c := testGetBucketSize(db, "markups")
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

	assessSample(SampleID{ProjectID: proj0.ProjectID, SampleID: 0})
	assessSample(SampleID{ProjectID: proj0.ProjectID, SampleID: 1})
	assessSample(SampleID{ProjectID: proj0.ProjectID, SampleID: 2})

	assessSample(SampleID{ProjectID: proj1.ProjectID, SampleID: 0})
	assessSample(SampleID{ProjectID: proj1.ProjectID, SampleID: 2})

	c = testGetBucketSize(db, "markups")
	assert.Equal(t, 5, c)

	list, err := svc.ListMarkup(WithProjectIDRequest{
		ProjectID: proj0.ProjectID,
	})
	assert.Nil(t, err)

	{
		assert.ElementsMatch(t, []MarkupListElement{
			NewMarkupListElement(
				SampleID{ProjectID: proj0.ProjectID, SampleID: 0},
				"sampleuri0",
				SampleMarkup{
					CreatedAt: utils.TestNowUTC(),
					Markup:    json.RawMessage(`{"kek":mark0}`),
				},
			),
			NewMarkupListElement(
				SampleID{ProjectID: proj0.ProjectID, SampleID: 1},
				"sampleuri1",
				SampleMarkup{
					CreatedAt: utils.TestNowUTC(),
					Markup:    json.RawMessage(`{"kek":mark1}`),
				},
			),
			NewMarkupListElement(
				SampleID{ProjectID: proj0.ProjectID, SampleID: 2},
				"sampleuri2",
				SampleMarkup{
					CreatedAt: utils.TestNowUTC(),
					Markup:    json.RawMessage(`{"kek":mark2}`),
				},
			),
		}, list.List)
	}
}
