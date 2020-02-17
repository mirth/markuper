package internal

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-collections/collections/set"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"

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

func getAllSampleIDs(db *DB, bucket string) ([]SampleID, error) {
	sIDs := make([]SampleID, 0)

	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.ForEach(func(k, _v []byte) error {
			sID := SampleID{}

			{
				err := decodeBin(k).Decode(&sID)
				if err != nil {
					return errors.WithStack(err)
				}

				sIDs = append(sIDs, sID)
			}

			return nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return sIDs, nil
}

func (s *MarkupServiceImpl) GetNext() (SampleResponse, error) {
	// fixme lock sample

	tmp, err := getAllSampleIDs(s.db, "markups")
	if err != nil {
		return SampleResponse{}, err
	}

	doneIDs := set.New(utils.ToSliceOfInterfaces(tmp)...)
	tmp, err = getAllSampleIDs(s.db, "samples")

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
	sample, err := s.db.GetSample(sID)

	if err != nil {
		return SampleResponse{}, err
	}

	proj, err := s.db.GetProject(sID.ProjectID)
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
	err := s.db.Put("markups", r.SampleID, r.SampleMarkup)

	return err
}

func (svc *MarkupServiceImpl) ListMarkup() (MarkupList, error) {
	ids, err := getAllSampleIDs(svc.db, "markups")
	if err != nil {
		return MarkupList{}, err
	}

	samples := []MarkupListElement{}
	err = svc.db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Markups)
		for _, id := range ids {
			binID, err := encodeBin(id)
			if err != nil {
				return err
			}
			smBin := b.Get(binID)
			sm := SampleMarkup{}
			err = decodeBin(smBin).Decode(&sm)
			if err != nil {
				return err
			}

			samples = append(samples, MarkupListElement{
				SampleID:     id,
				SampleMarkup: sm,
			})
		}

		return nil
	})

	if err != nil {
		return MarkupList{}, err
	}

	return MarkupList{
		List: samples,
	}, nil
}

func ListMarkupEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return s.ListMarkup()
	}
}
