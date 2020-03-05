package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportCSV(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	proj := fillTestDB(db)

	m := MarkupServiceImpl{
		db: db,
	}

	AssessWithMarkup(t, m, SampleID{ProjectID: proj.ProjectID, SampleID: 0}, `{"kek":"mark0","aaa": 3}`)
	AssessWithMarkup(t, m, SampleID{ProjectID: proj.ProjectID, SampleID: 1}, `{"kek":"mark1","aaa": 4}`)
	AssessWithMarkup(t, m, SampleID{ProjectID: proj.ProjectID, SampleID: 2}, `{"kek":"mark2","aaa": 5}`)

	s := ExporterServiceImpl{
		db: db,
	}

	r, err := s.Export(WithHttpRequest{
		R: nil,
		Payload: WithProjectIDRequest{
			ProjectID: proj.ProjectID,
		},
	})
	assert.Nil(t, err)

	{
		assert.Equal(t, `sample_id,assessed_at,aaa,kek
0,2015-03-07T11:06:39,3,"""mark0"""
1,2015-03-07T11:06:39,4,"""mark1"""
2,2015-03-07T11:06:39,5,"""mark2"""
`, string(r.CSV))
	}
}
