package service

import (
	"context"
	"fmt"

	"storage-service/internal/entity"
	"storage-service/internal/repository"
)

type CalculationService struct {
	repo repository.CalculationRepository
}

func NewCalculationService(repo repository.CalculationRepository) *CalculationService {
	return &CalculationService{repo: repo}
}

func (s *CalculationService) ProcessCalculation(ctx context.Context, calc *entity.Calculation) error {
	if err := calc.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.Save(ctx, calc.Expression, calc.Result); err != nil {
		return fmt.Errorf("failed to persist calculation: %w", err)
	}

	return nil
}
