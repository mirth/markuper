package internal

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-collections/collections/set"
	"github.com/pkg/errors"
	"github.com/recoilme/pudge"

	"backend/pkg/utils"
)

type SampleMarkup struct {
	CreatedAt time.Time `json:"created_at"`

	Markup json.RawMessage `json:"markup"`
}

type AssessRequest struct {
	SampleID     SampleID     `json:"sample_id"`
	SampleMarkup SampleMarkup `json:"sample_markup"`
}

type SampleResponse struct {
	SampleID SampleID        `json:"sample_id"`
	Sample   json.RawMessage `json:"sample"`
	Template Template        `json:"template"`
}

type MarkupListElement = AssessRequest
type MarkupList struct {
	List []MarkupListElement `json:"list"`
}

type MarkupService interface {
	GetNext() (SampleResponse, error)
	Assess(AssessRequest) error

	ListMarkup() (MarkupList, error)
}

type MarkupServiceImpl struct {
	db *DB
}

func NewMarkupService(db *DB) MarkupService {
	return &MarkupServiceImpl{
		db: db,
	}
}

func NextSampleEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return s.GetNext()
	}
}

func AssessEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r := *request.(*AssessRequest)
		err := s.Assess(r)

		return nil, err
	}
}

func getAllSampleIDs(db *pudge.Db) ([]SampleID, error) {
	rawIDs, err := db.Keys(SampleID{}, 0, 0, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sIDs := make([]SampleID, 0)
	for _, rawKey := range rawIDs {
		key := *decodeBinary(rawKey, func() interface{} {
			return &SampleID{}
		}).(*SampleID)
		sIDs = append(sIDs, key)
	}

	return sIDs, nil
}

func (s *MarkupServiceImpl) GetNext() (SampleResponse, error) {
	tmp, err := getAllSampleIDs(s.db.Markup)
	if err != nil {
		return SampleResponse{}, err
	}
	doneIDs := set.New(utils.ToSliceOfInterfaces(tmp)...)
	tmp, err = getAllSampleIDs(s.db.Sample)
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

	// FIXME empty toAssess
	sID := toAssess[0]
	sample := []byte{}
	err = s.db.Sample.Get(sID, &sample)
	if err != nil {
		return SampleResponse{}, errors.WithStack(err)
	}

	proj := Project{}
	err = s.db.Project.Get(sID.ProjectID, &proj)
	if err != nil {
		return SampleResponse{}, errors.WithStack(err)
	}

	return SampleResponse{
		SampleID: sID,
		Sample:   sample,
		Template: proj.Template,
	}, err
}

func (s *MarkupServiceImpl) Assess(r AssessRequest) error {
	r.SampleMarkup.CreatedAt = utils.NowUTC()
	err := s.db.Markup.Set(r.SampleID, r.SampleMarkup)

	return err
}

func (svc *MarkupServiceImpl) ListMarkup() (MarkupList, error) {
	ids, err := getAllSampleIDs(svc.db.Markup)
	if err != nil {
		return MarkupList{}, err
	}

	samples := []MarkupListElement{}
	for _, id := range ids {
		s := SampleMarkup{}
		err := svc.db.Markup.Get(id, &s)
		if err != nil {
			return MarkupList{}, errors.WithStack(err)
		}

		samples = append(samples, MarkupListElement{
			SampleID:     id,
			SampleMarkup: s,
		})
	}

	return MarkupList{
		List: samples,
	}, nil
}
