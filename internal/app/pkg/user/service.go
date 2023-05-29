package user

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (int, error)
	GetById(ctx context.Context, userID int) (User, error)
	GetByEmailPassword(ctx context.Context, email string, passwordHash string) (User, error)
	Update(ctx context.Context, user User) (bool, error)
	Delete(ctx context.Context, userID int) (bool, error)
}

type Service struct {
	repository UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Create(ctx context.Context, user User) (int, error) {
	id, err := s.repository.Create(ctx, user)
	return id, err
}

func (s *Service) GetById(ctx context.Context, userID int) (User, error) {
	user, err := s.repository.GetById(ctx, userID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) GetByEmailPassword(ctx context.Context, email string, passwordHash string) (User, error) {
	user, err := s.repository.GetByEmailPassword(ctx, email, passwordHash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) Update(ctx context.Context, user User) (bool, error) {
	ok, err := s.repository.Update(ctx, user)
	return ok, err
}

func (s *Service) Delete(ctx context.Context, userID int) (bool, error) {
	ok, err := s.repository.Delete(ctx, userID)
	return ok, err
}
