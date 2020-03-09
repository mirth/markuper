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

type SampleWithMarkupResponse struct {
	SampleResponse
	SampleMarkup *SampleMarkup `json:"markup"`
}

type MarkupListElement struct {
	SampleID     SampleID     `json:"sample_id"`
	SampleURI    SampleURI    `json:"sample_uri"`
	SampleMarkup SampleMarkup `json:"sample_markup"`
}

type MarkupList struct {
	List []MarkupListElement `json:"list"`
}

type MarkupService interface {
	GetNext(WithProjectIDRequest) (SampleResponse, error)
	Assess(AssessRequest) error

	GetSample(SampleID) (SampleWithMarkupResponse, error)
	ListMarkup(WithProjectIDRequest) (MarkupList, error)
}

type WithProjectIDRequest struct {
	ProjectID ProjectID `json:"project_id"`
}

type MarkupServiceImpl struct {
	db *DB
}

func NewMarkupListElement(id SampleID, uri SampleURI, m SampleMarkup) MarkupListElement {
	return MarkupListElement{id, uri, m}
}

func NewMarkupService(db *DB) MarkupService {
	return &MarkupServiceImpl{
		db: db,
	}
}

func NextSampleEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := *request.(*WithProjectIDRequest)

		return s.GetNext(req)
	}
}

func AssessEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		r := *request.(*AssessRequest)
		err := s.Assess(r)

		return nil, err
	}
}

func getAllSampleIDsForProject(db *DB, bucket string, projectID ProjectID) ([]SampleID, error) {
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

				if sID.ProjectID == projectID {
					sIDs = append(sIDs, sID)
				}
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

func (s *MarkupServiceImpl) GetNext(req WithProjectIDRequest) (SampleResponse, error) {
	// fixme lock sample

	tmp, err := getAllSampleIDsForProject(s.db, "markups", req.ProjectID)
	if err != nil {
		return SampleResponse{}, err
	}

	doneIDs := set.New(utils.ToSliceOfInterfaces(tmp)...)
	tmp, err = getAllSampleIDsForProject(s.db, "samples", req.ProjectID)

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

func ListMarkup(db *DB, projectID ProjectID) (MarkupList, error) {
	ids, err := getAllSampleIDsForProject(db, "markups", projectID)
	if err != nil {
		return MarkupList{}, err
	}

	samples := []MarkupListElement{}
	err = db.DB.View(func(tx *bolt.Tx) error {
		m := tx.Bucket(Markups)
		s := tx.Bucket(Samples)
		for _, id := range ids {
			binID, err := encodeBin(id)
			if err != nil {
				return err
			}
			smBin := m.Get(binID)
			sm := SampleMarkup{}
			{
				err := decodeBin(smBin).Decode(&sm)
				if err != nil {
					return err
				}
			}

			sBin := s.Get(binID)
			sJson := []byte{}
			sample := ImageSample{}
			{
				err := decodeBin(sBin).Decode(&sJson)
				if err != nil {
					return err
				}
				err = json.Unmarshal(sJson, &sample)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			samples = append(samples, NewMarkupListElement(
				id,
				sample.ImageURI,
				sm,
			))
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

func (s *MarkupServiceImpl) GetSample(sID SampleID) (SampleWithMarkupResponse, error) {
	sample, err := s.db.GetSample(sID)
	if err != nil {
		return SampleWithMarkupResponse{}, err
	}

	markup, err := s.db.GetMarkup(sID)
	if err != nil {
		return SampleWithMarkupResponse{}, err
	}

	proj, err := s.db.GetProject(sID.ProjectID)
	if err != nil {
		return SampleWithMarkupResponse{}, err
	}

	return SampleWithMarkupResponse{
		SampleResponse: SampleResponse{
			SampleID: sID,
			Sample:   sample,
			Template: proj.Template,
		},
		SampleMarkup: markup,
	}, nil
}

func (s *MarkupServiceImpl) ListMarkup(req WithProjectIDRequest) (MarkupList, error) {
	return ListMarkup(s.db, req.ProjectID)
}

func ListMarkupEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := *request.(*WithProjectIDRequest)

		return s.ListMarkup(req)
	}
}
