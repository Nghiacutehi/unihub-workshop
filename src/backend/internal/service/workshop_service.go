package service

import (
	"context"
	"fmt"
	"time"

	"unihub-workshop/internal/model"
	"unihub-workshop/internal/repository"
)

type WorkshopService struct {
	repo *repository.WorkshopRepo
}

func NewWorkshopService(repo *repository.WorkshopRepo) *WorkshopService {
	return &WorkshopService{repo: repo}
}

func (s *WorkshopService) ListAll(ctx context.Context) ([]model.Workshop, error) {
	return s.repo.FindAll(ctx)
}

func (s *WorkshopService) GetByID(ctx context.Context, id string) (*model.Workshop, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *WorkshopService) Create(ctx context.Context, req *model.CreateWorkshopRequest) (*model.Workshop, error) {
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %w", err)
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format: %w", err)
	}

	w := &model.Workshop{
		Title:          req.Title,
		Description:    req.Description,
		Speaker:        req.Speaker,
		Room:           req.Room,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       req.Capacity,
		AvailableSeats: req.Capacity,
		Price:          req.Price,
		Status:         model.WorkshopPublished,
	}

	if err := s.repo.Create(ctx, w); err != nil {
		return nil, fmt.Errorf("failed to create workshop: %w", err)
	}
	return w, nil
}

func (s *WorkshopService) Update(ctx context.Context, id string, req *model.UpdateWorkshopRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *WorkshopService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
