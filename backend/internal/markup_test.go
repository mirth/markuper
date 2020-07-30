package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/pkg/utils"
)

func fillTestDB(db *DB) Project {
	return fillTestDBWithProj(db, "testp0", DEFAULT_CLASSIFICATION_TEMPLATE)
}

func newSampleIDForTest(projID ProjectID, i int) string {
	return fmt.Sprintf("%s-%d", string(projID), i)
}

func fillTestDBWithProj(db *DB, projName string, template TemplateXML) Project {
	svc := NewProjectService(db)
	p, _ := svc.CreateProject(newTestCreateProjectRequest(projName, template))

	samples := make([]MediaSample, 0)
	for i := 0; i < 5; i++ {
		samples = append(samples, MediaSample{
			MediaURI:  fmt.Sprintf("sampleuri%d", i),
			MediaType: IMAGE_FILE_TYPE,
		})
	}

	for i, sample := range samples {
		sk := newSampleIDForTest(p.ProjectID, i)

		j, _ := json.Marshal(sample)
		db.PutOne("samples", sk, j)
	}

	return p
}

func AssessWithMarkup(
	t *testing.T,
	s MarkupServiceImpl,
	sID SampleID,
	markup string,
) {
	mkp := &SampleMarkup{
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

		AssessWithMarkup(t, svc, newSampleIDForTest("kek", 0), `{"kek":"kek"}`)
	}
}

func TestMarkupNext(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	proj := fillTestDB(db)

	svc := MarkupServiceImpl{
		db: db,
	}

	assertNext := func(i int) SampleID {
		a, err := svc.GetNext(WithProjectIDRequest{
			ProjectID: proj.ProjectID,
		})
		assert.Nil(t, err)
		sID := newSampleIDForTest(proj.ProjectID, i)

		e := SampleResponse{
			SampleID: sID,
			Sample:   json.RawMessage(fmt.Sprintf(`{"media_uri":"sampleuri%d","media_type":"IMAGE_FILE_TYPE"}`, i)),
			Project:  proj,
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

func TestGetSample(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)

	proj0 := fillTestDBWithProj(db, "proj0", DEFAULT_CLASSIFICATION_TEMPLATE)

	svc := MarkupServiceImpl{
		db: db,
	}

	sID := newSampleIDForTest(proj0.ProjectID, 0)
	es, _ := db.GetSample(sID)

	{
		s, err := svc.GetSample(sID)
		assert.Nil(t, err)

		assert.Equal(t, es, []byte(s.Sample))
		assert.Equal(t, sID, s.SampleID)
		assert.Equal(t, proj0, s.Project)
		assert.Nil(t, s.SampleMarkup)
	}

	assessSample(t, &svc, sID)

	{
		s, err := svc.GetSample(sID)
		assert.Nil(t, err)

		em, _ := db.GetMarkup(sID)
		assert.Equal(t, es, []byte(s.Sample))
		assert.Equal(t, sID, s.SampleID)
		assert.Equal(t, proj0, s.Project)
		assert.Equal(t, em, s.SampleMarkup)
	}
}

func assessSample(t *testing.T, svc MarkupService, sID SampleID) {
	r := AssessRequest{
		SampleID: sID,
		SampleMarkup: SampleMarkup{
			Markup: json.RawMessage(fmt.Sprintf(`{"kek":"mark%s"}`, sID)),
		},
	}
	err := svc.Assess(r)
	assert.Nil(t, err)
}

func TestListMarkup(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	proj0 := fillTestDBWithProj(db, "proj0", DEFAULT_CLASSIFICATION_TEMPLATE)
	proj1 := fillTestDBWithProj(db, "proj1", DEFAULT_CLASSIFICATION_TEMPLATE)

	svc := MarkupServiceImpl{
		db: db,
	}

	c := testGetBucketSize(db, "markups")
	assert.Zero(t, c)

	ass := func(sID SampleID) {
		assessSample(t, &svc, sID)
	}

	ass(newSampleIDForTest(proj0.ProjectID, 0))
	ass(newSampleIDForTest(proj0.ProjectID, 1))
	ass(newSampleIDForTest(proj0.ProjectID, 2))

	ass(newSampleIDForTest(proj1.ProjectID, 0))
	ass(newSampleIDForTest(proj1.ProjectID, 2))

	c = testGetBucketSize(db, "markups")
	assert.Equal(t, 5, c)

	list, err := svc.ListMarkup(WithProjectIDRequest{
		ProjectID: proj0.ProjectID,
	})
	assert.Nil(t, err)

	newSampleMarkup := func(i int) SampleMarkup {
		sID := newSampleIDForTest(proj0.ProjectID, i)
		return SampleMarkup{
			CreatedAt: utils.TestNowUTC(),
			Markup:    json.RawMessage(fmt.Sprintf(`{"kek":"mark%s"}`, sID)),
		}
	}

	{
		assert.ElementsMatch(t, []MarkupListElement{
			NewMarkupListElement(
				newSampleIDForTest(proj0.ProjectID, 0),
				"sampleuri0",
				newSampleMarkup(0),
			),
			NewMarkupListElement(
				newSampleIDForTest(proj0.ProjectID, 1),
				"sampleuri1",
				newSampleMarkup(1),
			),
			NewMarkupListElement(
				newSampleIDForTest(proj0.ProjectID, 2),
				"sampleuri2",
				newSampleMarkup(2),
			),
		}, list.List)
	}
}

func TestOutOfSamples(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	proj0 := fillTestDBWithProj(db, "proj0", DEFAULT_CLASSIFICATION_TEMPLATE)

	svc := MarkupServiceImpl{
		db: db,
	}

	ass := func(sID SampleID) {
		assessSample(t, &svc, sID)
	}
	getNext := func() SampleResponse {
		s, err := svc.GetNext(WithProjectIDRequest{
			ProjectID: proj0.ProjectID,
		})
		assert.Nil(t, err)

		return s
	}

	ass(newSampleIDForTest(proj0.ProjectID, 0))
	getNext()
	ass(newSampleIDForTest(proj0.ProjectID, 1))
	getNext()
	ass(newSampleIDForTest(proj0.ProjectID, 2))
	getNext()
	ass(newSampleIDForTest(proj0.ProjectID, 3))
	getNext()
	ass(newSampleIDForTest(proj0.ProjectID, 4))
	s := getNext()
	assert.Nil(t, s.Sample)
	assert.Equal(t, s.Project, proj0)
}
