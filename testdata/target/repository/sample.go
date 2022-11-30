package repository

import "time"

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

type IFSampleRepository interface {
	GetName(i int) (string, error)
	GetLastSaveTime(i int) (time.Time, error)
	Update(i int, name string) error
}

type SampleRepository struct {
}

func NewSampleRepository() *SampleRepository {
	return &SampleRepository{}
}

func (r *SampleRepository) GetName(i int) (string, error) {
	return "Sample", nil
}

func (r *SampleRepository) GetLastSaveTime(i int) (time.Time, error) {
	return time.Time{}, nil
}

func (r *SampleRepository) Update(i int, name string) error {
	return nil
}
