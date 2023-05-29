package order_dish

import (
	"context"
)

type OrderDishRepository interface {
	Create(ctx context.Context, orderDish OrderDish) (int, error)
	GetById(ctx context.Context, orderDishID int) (OrderDish, error)
	Update(ctx context.Context, orderDish OrderDish) (bool, error)
	Delete(ctx context.Context, orderDishID int) (bool, error)
}

type Service struct {
	repository OrderDishRepository
}

func NewService(repo OrderDishRepository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Create(ctx context.Context, orderDish OrderDish) (int, error) {
	id, err := s.repository.Create(ctx, orderDish)
	return id, err
}

func (s *Service) GetById(ctx context.Context, orderDishID int) (OrderDish, error) {
	order, err := s.repository.GetById(ctx, orderDishID)
	if err != nil {
		return OrderDish{}, err
	}
	return order, err
}

func (s *Service) Update(ctx context.Context, orderDish OrderDish) (bool, error) {
	ok, err := s.repository.Update(ctx, orderDish)
	return ok, err
}

func (s *Service) Delete(ctx context.Context, orderDishID int) (bool, error) {
	ok, err := s.repository.Delete(ctx, orderDishID)
	return ok, err
}
