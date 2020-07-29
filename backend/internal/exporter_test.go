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
		<radio group="kek" value="mark1" vizual="Mark1" />
		<radio group="kek" value="mark2" vizual="Mark2" />
		<radio group="kek" value="mark3" vizual="Mark3" />

		<radio group="aaa" value="3" vizual="aaa1" />
		<radio group="aaa" value="4" vizual="aaa2" />
		<radio group="aaa" value="5" vizual="aaa3" />

		<checkbox group="lel" value="l1" vizual="L1" />
		<checkbox group="lel" value="l2" vizual="L2" />
		<checkbox group="lel" value="l3" vizual="L3" />
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
		`{"kek":"mark0","aaa":"3","lel":["l1","l2"]}`,
	)
	AssessWithMarkup(
		t,
		m,
		newSampleIDForTest(proj.ProjectID, 1),
		`{"kek":"mark1","aaa":"4","lel":[]}`,
	)
	AssessWithMarkup(
		t,
		m,
		newSampleIDForTest(proj.ProjectID, 2),
		`{"kek":"mark2","aaa":"5","lel":["l3"]}`,
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
		assert.Equal(t, fmt.Sprintf(`sample_id,sample_uri,assessed_at,aaa,kek,lel
%s-0,sampleuri0,2015-03-07T11:06:39,"""3""","""mark0""","[""l1"",""l2""]"
%s-1,sampleuri1,2015-03-07T11:06:39,"""4""","""mark1""",[]
%s-2,sampleuri2,2015-03-07T11:06:39,"""5""","""mark2""","[""l3""]"
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
				<radio group="animal" value="cat" vizual="Cat"/>
				<radio group="animal" value="dog" vizual="Dog"/>
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
		`{"box":[{"box":{"x":0,"y":0,"width":2,"height":1},"animal":"cat"}]}`,
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
