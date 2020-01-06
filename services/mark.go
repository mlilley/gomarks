package services

import (
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/repos"
)

type MarkService interface {
	GetMarksForUser(user *app.User) ([]app.Mark, error)
}

func NewMarkService(markRepo repos.MarkRepo) MarkService {
	return &markService{markRepo: markRepo}
}

type markService struct {
	markRepo repos.MarkRepo
}

func (s *markService) GetMarksForUser(user *app.User) ([]app.Mark, error) {
	marks, err := s.markRepo.FindAllByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return marks, nil
}
