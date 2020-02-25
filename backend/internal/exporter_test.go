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

	AssessWithMarkup(t, m, SampleID{ProjectID: proj.ProjectID, SampleID: 0}, `{"kek":mark0}`)
	AssessWithMarkup(t, m, SampleID{ProjectID: proj.ProjectID, SampleID: 1}, `{"kek":mark1}`)
	AssessWithMarkup(t, m, SampleID{ProjectID: proj.ProjectID, SampleID: 2}, `{"kek":mark2}`)

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
		assert.Equal(t, `sample_id,created_at,markup
0,2015-03-07T11:06:39,"{""kek"":mark0}"
1,2015-03-07T11:06:39,"{""kek"":mark1}"
2,2015-03-07T11:06:39,"{""kek"":mark2}"
`, string(r.CSV))
	}
}
