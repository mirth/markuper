package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
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

type ProjectMeta struct {
	TotalNumberOfSamples    int `json:"total_number_of_samples"`
	AssessedNumberOfSamples int `json:"assessed_number_of_samples"`
}

type AssessRequest struct {
	SampleID     SampleID     `json:"sample_id"`
	SampleMarkup SampleMarkup `json:"sample_markup"`
}

type SampleRequest struct {
	SampleID SampleID `json:"sample_id"`
}

type SampleResponse struct {
	SampleID SampleID        `json:"sample_id"`
	Sample   json.RawMessage `json:"sample"`
	Project  Project         `json:"project"`
}

func NewSampleResponse(
	sampleID SampleID,
	sample json.RawMessage,
	project Project,
) SampleResponse {
	return SampleResponse{
		SampleID: sampleID,
		Sample:   sample,
		Project:  project,
	}
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
		c := tx.Bucket([]byte(bucket)).Cursor()

		prefix := []byte(projectID)
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			sIDs = append(sIDs, string(k))
		}

		return nil
	})

	return sIDs, errors.WithStack(err)
}

func getAllSamplesForProject(db *DB, projectID ProjectID) ([]json.RawMessage, error) {
	samples := make([]json.RawMessage, 0)

	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("samples"))
		err := b.ForEach(func(k, v []byte) error {
			sID := string(k)
			s := json.RawMessage{}

			{
				err := json.Unmarshal(v, &s)
				if err != nil {
					return errors.WithStack(err)
				}

				if GetProjectIDFromSampleID(sID) == projectID {
					samples = append(samples, s)
				}
			}

			return nil
		})

		return err
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return samples, nil
}

type SampleWithOrder struct {
	SampleID string
	Index    uint
}

func indexSamples(ids []string) map[SampleID]int {
	m := map[SampleID]int{}
	for i, s := range ids {
		m[s] = i
	}

	return m
}

func findSampleWithMinIndex(toAssess []string, allIDsIndex map[SampleID]int) SampleID {
	minIndex := math.MaxInt32
	var sampleID SampleID
	for _, sID := range toAssess {
		index, ok := allIDsIndex[sID]
		if ok && index < minIndex {
			minIndex = index
			sampleID = sID
		}
	}

	return sampleID
}

func (s *MarkupServiceImpl) GetNext(req WithProjectIDRequest) (SampleResponse, error) {
	tmp, err := getAllSampleIDsForProject(s.db, "markups", req.ProjectID)
	if err != nil {
		return SampleResponse{}, err
	}
	doneIDs := set.New(utils.ToSliceOfInterfaces(tmp)...)

	tmp, err = getAllSampleIDsForProject(s.db, "samples", req.ProjectID)
	if err != nil {
		return SampleResponse{}, err
	}

	allIDsIndex := indexSamples(tmp)
	allIDs := set.New(utils.ToSliceOfInterfaces(tmp)...)

	toAssess := make([]SampleID, 0)
	allIDs.Difference(doneIDs).Do(func(sID interface{}) {
		toAssess = append(toAssess, sID.(SampleID))
	})

	proj, err := s.db.GetProject(req.ProjectID)
	if err != nil {
		return SampleResponse{}, errors.WithStack(err)
	}

	if len(toAssess) == 0 {
		return SampleResponse{
			Project: proj,
		}, nil
	}

	sID := findSampleWithMinIndex(toAssess, allIDsIndex)

	sample, err := s.db.GetSample(sID)

	if err != nil {
		return SampleResponse{}, err
	}

	return NewSampleResponse(sID, sample, proj), err
}

func (s *MarkupServiceImpl) Assess(r AssessRequest) error {
	r.SampleMarkup.CreatedAt = utils.NowUTC()
	m, err := json.Marshal(r.SampleMarkup)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.db.PutOne("markups", r.SampleID, m)

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
		s := tx.Bucket([]byte(Samples))
		for _, id := range ids {
			binID := []byte(id)
			smJson := m.Get(binID)
			sm := SampleMarkup{}
			{
				err := json.Unmarshal(smJson, &sm)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			sJson := s.Get(binID)
			sample := MediaSample{}
			{
				err = json.Unmarshal(sJson, &sample)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			samples = append(samples, NewMarkupListElement(
				id,
				sample.MediaURI,
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

	proj, err := s.db.GetProject(GetProjectIDFromSampleID(sID))
	if err != nil {
		return SampleWithMarkupResponse{}, err
	}

	return SampleWithMarkupResponse{
		SampleResponse: NewSampleResponse(sID, sample, proj),
		SampleMarkup:   markup,
	}, nil
}

func GetSampleEndpoint(s MarkupService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := *request.(*SampleRequest)

		return s.GetSample(req.SampleID)
	}
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
