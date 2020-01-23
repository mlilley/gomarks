package services

import (
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/repos"
)

type UserService interface {
	GetUserByID(id string) (*app.User, error)
}

func NewUserService(userRepo repos.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

type userService struct {
	userRepo repos.UserRepo
}

func (s *userService) GetUserByID(id string) (*app.User, error) {
	return s.userRepo.FindByID(id)
}
