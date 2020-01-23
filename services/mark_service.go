package services

import (
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/repos"
)

type MarkService interface {
	GetMarksForUser(userId string) ([]app.Mark, error)
	GetMarkByIDForUser(markId string, userId string) (*app.Mark, error)
	CreateMarkForUser(mark *app.Mark, userId string) (*app.Mark, error)
	UpdateMarkForUser(mark *app.Mark, userId string) (*app.Mark, error)
	DeleteMarkByIDForUser(markId string, userId string) (bool, error)
	DeleteMarksForUser(userId string) error
}

func NewMarkService(markRepo repos.MarkRepo) MarkService {
	return &markService{markRepo: markRepo}
}

type markService struct {
	markRepo repos.MarkRepo
}

func (s *markService) GetMarksForUser(userId string) ([]app.Mark, error) {
	marks, err := s.markRepo.FindAllForUser(userId)
	if err != nil {
		return nil, err
	}

	return marks, nil
}

// Returned mark is nil if not found.
func (s *markService) GetMarkByIDForUser(markId string, userId string) (*app.Mark, error) {
	mark, err := s.markRepo.FindByIDForUser(markId, userId)
	if err != nil {
		return nil, err
	}

	return mark, err
}

// Ignores ID and UserID properties (UserID forced to userId).
func (s *markService) CreateMarkForUser(mark *app.Mark, userId string) (*app.Mark, error) {
	mark.ID = ""
	mark.UserID = userId
	mark, err := s.markRepo.Create(mark)
	if err != nil {
		return nil, err
	}

	return mark, nil
}

// Ignores ID and UserID properties (UserID forced to userId).
// Returned mark is nil if not found.
func (s *markService) UpdateMarkForUser(mark *app.Mark, userId string) (*app.Mark, error) {
	mark.UserID = userId
	mark, err := s.markRepo.UpdateWithResult(mark)
	if err != nil {
		return nil, err
	}

	return mark, nil
}

func (s *markService) DeleteMarkByIDForUser(markId string, userId string) (bool, error) {
	ok, err := s.markRepo.DeleteByIDForUser(markId, userId)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *markService) DeleteMarksForUser(userId string) error {
	err := s.markRepo.DeleteAllForUser(userId)
	if err != nil {
		return err
	}

	return nil
}