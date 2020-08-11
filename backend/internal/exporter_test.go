package internal

import (
	"backend/pkg/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportCSV(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	template := TemplateXML{
		Task: "some_task",
		XML: `
	<content>
		<radio group="kek" value="mark1" visual="Mark1" />
		<radio group="kek" value="mark2" visual="Mark2" />
		<radio group="kek" value="mark3" visual="Mark3" />

		<radio group="aaa" value="3" visual="aaa1" />
		<radio group="aaa" value="4" visual="aaa2" />
		<radio group="aaa" value="5" visual="aaa3" />

		<checkbox group="lel" value="l1" visual="L1" />
		<checkbox group="lel" value="l2" visual="L2" />
		<checkbox group="lel" value="l3" visual="L3" />
	</content>`,
	}

	proj := fillTestDBWithProj(db, "testp0", template)

	m := MarkupServiceImpl{
		db: db,
	}

	AssessWithMarkup(
		t,
		m,
		newSampleIDForTest(proj.ProjectID, 0),
		[]byte(`{"kek":"mark0","aaa":"3","lel":["l1","l2"]}`),
	)
	AssessWithMarkup(
		t,
		m,
		newSampleIDForTest(proj.ProjectID, 1),
		[]byte(`{"kek":"mark1","aaa":"4","lel":[]}`),
	)
	AssessWithMarkup(
		t,
		m,
		newSampleIDForTest(proj.ProjectID, 2),
		[]byte(`{"kek":"mark2","aaa":"5","lel":["l3"]}`),
	)

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
		pID := proj.ProjectID
		assert.Equal(t, fmt.Sprintf(`sample_id,sample_uri,assessed_at,kek,aaa,lel
%s-0,sampleuri0,2015-03-07T11:06:39,mark0,3,"[""l1"",""l2""]"
%s-1,sampleuri1,2015-03-07T11:06:39,mark1,4,[]
%s-2,sampleuri2,2015-03-07T11:06:39,mark2,5,"[""l3""]"
`, pID, pID, pID), string(r.CSV))

		assert.Equal(
			t,
			fmt.Sprintf("%s_%s.csv", proj.Description.Name,
				utils.TestNowUTC().Format("2006-01-02T15:04:05")),
			r.Filename,
		)
	}
}

func TestExportCSVForBBox(t *testing.T) {
	db := openTestDB()
	defer testCloseAndReset(db)
	template := TemplateXML{
		Task: "some_task",
		XML: `
		<content>
			<bounding_box group="box">
				<radio group="animal" value="cat" visual="Cat"/>
				<radio group="animal" value="dog" visual="Dog"/>
			</bounding_box>
		</content>`,
	}

	proj := fillTestDBWithProj(db, "testp0", template)
	m := MarkupServiceImpl{
		db: db,
	}

	AssessWithMarkup(
		t,
		m,
		newSampleIDForTest(proj.ProjectID, 0),
		[]byte(`{"box":[{"box":{"x":0,"y":0,"width":2,"height":1},"animal":"cat"}]}`),
	)

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
		pID := proj.ProjectID
		assert.Equal(t, fmt.Sprintf(`sample_id,sample_uri,assessed_at,box
%s-0,sampleuri0,2015-03-07T11:06:39,"[{""box"":{""x"":0,""y"":0,""width"":2,""height"":1},""animal"":""cat""}]"
`,
			pID), string(r.CSV))
	}
}
