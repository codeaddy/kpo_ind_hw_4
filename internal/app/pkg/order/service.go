package order

import (
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order Order) (int, error)
	GetById(ctx context.Context, orderID int) (Order, error)
	Update(ctx context.Context, order Order) (bool, error)
	Delete(ctx context.Context, orderID int) (bool, error)
}

type Service struct {
	repository OrderRepository
}

func NewService(repo OrderRepository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Create(ctx context.Context, order Order) (int, error) {
	id, err := s.repository.Create(ctx, order)
	return id, err
}

func (s *Service) GetById(ctx context.Context, orderID int) (Order, error) {
	order, err := s.repository.GetById(ctx, orderID)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (s *Service) Update(ctx context.Context, order Order) (bool, error) {
	ok, err := s.repository.Update(ctx, order)
	return ok, err
}

func (s *Service) Delete(ctx context.Context, orderID int) (bool, error) {
	ok, err := s.repository.Delete(ctx, orderID)
	return ok, err
}
