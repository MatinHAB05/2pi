package service

import (
	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/pkg/random"
)

type randomService struct {
	engine *random.RandomGeneratorEngine
}

func NewRandomService() service_contract.RandomService {
	return &randomService{
		engine: random.NewRandomGeneratorEngine(),
	}
}

func (s *randomService) GeneratePositiveIntRandomNumber(digits int) int64 {
	return s.engine.GeneratePositiveIntRandomNumber(digits)
}

func (s *randomService) GeneratePositiveIntRandomNumberInRange(min, max int64) int64 {
	return s.engine.GeneratePositiveIntRandomNumberInRange(min, max)
}
