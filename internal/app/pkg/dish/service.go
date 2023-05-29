package dish

import (
	"context"
)

type DishRepository interface {
	Create(ctx context.Context, dish Dish) (int, error)
	GetById(ctx context.Context, dishID int) (Dish, error)
	GetAll(ctx context.Context) ([]*Dish, error)
	Update(ctx context.Context, dish Dish) (bool, error)
	Delete(ctx context.Context, filmID int) (bool, error)
}

type Service struct {
	repository DishRepository
}

func NewService(repo DishRepository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Create(ctx context.Context, dish Dish) (int, error) {
	id, err := s.repository.Create(ctx, dish)
	return id, err
}

func (s *Service) GetById(ctx context.Context, dishID int) (Dish, error) {
	dish, err := s.repository.GetById(ctx, dishID)
	if err != nil {
		return Dish{}, err
	}
	return dish, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*Dish, error) {
	dish, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return dish, nil
}

func (s *Service) Update(ctx context.Context, dish Dish) (bool, error) {
	ok, err := s.repository.Update(ctx, dish)
	return ok, err
}

func (s *Service) Delete(ctx context.Context, dishID int) (bool, error) {
	ok, err := s.repository.Delete(ctx, dishID)
	return ok, err
}
