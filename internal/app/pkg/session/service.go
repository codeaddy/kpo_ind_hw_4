package session

import (
	"context"
	"github.com/golang-jwt/jwt"
	"time"
)

type SessionRepository interface {
	Create(ctx context.Context, session Session) (int, error)
	GetById(ctx context.Context, sessionID int) (Session, error)
	GetByUserId(ctx context.Context, userID int) (Session, error)
	GetByToken(ctx context.Context, token string) (Session, error)
	Update(ctx context.Context, session Session) (bool, error)
	Delete(ctx context.Context, sessionID int) (bool, error)
}

type Service struct {
	repository SessionRepository
}

func NewService(repo SessionRepository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Create(ctx context.Context, session Session) (int, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Now().Add(time.Hour),
	})
	tokenString, err := token.SignedString([]byte("megasnyus"))

	if err != nil {
		return 0, err
	}
	session.SessionToken = tokenString

	id, err := s.repository.Create(ctx, session)
	return id, err
}

func (s *Service) GetById(ctx context.Context, sessionID int) (Session, error) {
	order, err := s.repository.GetById(ctx, sessionID)
	if err != nil {
		return Session{}, err
	}
	return order, nil
}

func (s *Service) GetByUserId(ctx context.Context, userID int) (Session, error) {
	order, err := s.repository.GetByUserId(ctx, userID)
	if err != nil {
		return Session{}, err
	}
	return order, nil
}

func (s *Service) GetByToken(ctx context.Context, token string) (Session, error) {
	order, err := s.repository.GetByToken(ctx, token)
	if err != nil {
		return Session{}, err
	}
	return order, nil
}

func (s *Service) Update(ctx context.Context, session Session) (bool, error) {
	ok, err := s.repository.Update(ctx, session)
	return ok, err
}

func (s *Service) Delete(ctx context.Context, sessionID int) (bool, error) {
	ok, err := s.repository.Delete(ctx, sessionID)
	return ok, err
}
