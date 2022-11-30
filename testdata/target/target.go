package target

import (
	"errors"
	"time"

	"github.com/kazdevl/tgen/testdata/target/repository"
	"github.com/kazdevl/tgen/testdata/target/thirdparty"
)

type SampleService struct {
	SampleRepository repository.IFSampleRepository
	SampleClient     thirdparty.IFSampleClient
}

func NewSampleService(sr *repository.SampleRepository, sc *thirdparty.SampleClient) *SampleService {
	return &SampleService{
		SampleRepository: sr,
		SampleClient:     sc,
	}
}

func (s *SampleService) GetSampleName(i int) (string, error) {
	return s.SampleRepository.GetName(i)
}

func (s *SampleService) UpdateToRandomName(i int) error {
	updateName := s.SampleClient.GenrateRandomName()
	if s.isValid(i, updateName) {
		return "", errors.New("有効ではないです")
	}

	return s.SampleRepository.Update(i, name)
}

func (s *SampleService) isValid(i int, updateName string) bool {
	if err := s.isUpdatable(i, updateName); err != nil {
		return false
	}
	name, err := s.Get(i)
	if err != nil {
		return false
	}
	if len(name) != len(updateName) {
		return false
	}
	return true
}

func (s *SampleService) isUpdatable(i int, name string) error {
	lastSaveTime, err := s.SampleRepository.GetLastSaveTime(i)
	if err != nil {
		return err
	}
	now := time.Now()
	startDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if lastSaveTime.After(startDay) {
		return errors.New("今日すでに更新済みなら更新できません")
	}
	return nil
}
