package internal

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-collections/collections/set"
	"github.com/pkg/errors"
	"github.com/recoilme/pudge"

	"backend/pkg/utils"
)

type SampleID struct {
	ProjectID string `json:"project_id"`
	SampleID  int64  `json:"sample_id"`
}

func (s SampleID) toString() string {
	return fmt.Sprintf("%s|%d", s.ProjectID, s.SampleID)
}

type SampleMarkup struct {
	CreatedAt time.Time `json:"created_at"`

	Markup json.RawMessage `json:"markup"`
}

type SampleResponse struct {
	SampleID  SampleID `json:"sample_id"`
	SampleURI string   `json:"sample_uri"`
}

type AssessRequest struct {
	SampleID     SampleID     `json:"sample_id"`
	SampleMarkup SampleMarkup `json:"sample_markup"`
}

type MarkupService interface {
	GetNext() (SampleResponse, error)
	Assess(AssessRequest) error
}

type MarkupServiceImpl struct {
	SamplesDB string
	MarkupDB  string
}

func NextSampleEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		s, _ := s.GetNext()

		return s, nil
	}
}

func AssessEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r := *request.(*AssessRequest)
		err := s.Assess(r)

		return nil, err
	}
}

func decodeKey(raw []byte) SampleID {
	buf := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buf)
	key := SampleID{}
	_ = dec.Decode(&key)

	return key
}

var offset = 0

func getAllSampleIDs(db string) ([]SampleID, error) {
	rawIDs, err := pudge.Keys(db, SampleID{}, 0, 0, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sIDs := make([]SampleID, 0)
	for _, rawKey := range rawIDs {
		key := decodeKey(rawKey)
		sIDs = append(sIDs, key)
	}

	return sIDs, nil
}

func (s *MarkupServiceImpl) GetNext() (SampleResponse, error) {
	tmp, err := getAllSampleIDs(s.MarkupDB)
	if err != nil {
		return SampleResponse{}, err
	}
	doneIDs := set.New(utils.ToSliceOfInterfaces(tmp)...)

	tmp, err = getAllSampleIDs(s.SamplesDB)
	if err != nil {
		return SampleResponse{}, err
	}
	allIDs := set.New(utils.ToSliceOfInterfaces(tmp)...)

	toAssess := make([]SampleID, 0)
	allIDs.Difference(doneIDs).Do(func(sID interface{}) {
		toAssess = append(toAssess, sID.(SampleID))
	})

	sort.SliceStable(toAssess, func(i, j int) bool {
		return toAssess[i].SampleID < toAssess[j].SampleID
	})

	sID := toAssess[0]
	sampleURI := ""
	err = pudge.Get(s.SamplesDB, sID, &sampleURI)
	if err != nil {
		return SampleResponse{}, errors.WithStack(err)
	}

	return SampleResponse{
		SampleID:  sID,
		SampleURI: sampleURI,
	}, err
}

func (s *MarkupServiceImpl) Assess(r AssessRequest) error {
	err := pudge.Set(s.MarkupDB, r.SampleID, r.SampleMarkup)

	return err
}
