package services

import (
	"github.com/golang/mock/gomock"
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var validPassword = "password"
var invalidPassword = "invalid-password"
var secret = []byte("secret")
var user = app.User{ID: "1", Email: "user@test.com", PasswordHash: "$2a$14$RjtUXdY343NfD2xpp/2GguI/3wK5tlLguKh/mtRgdFwfyu7jWfgUy", Active: true}

func TestAuthServiceValidLogin(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockUserRepo := mocks.NewMockUserRepo(mockController)
	mockUserRepo.EXPECT().FindByEmail(user.Email).Return(&user, nil).Times(1)

	authService := NewAuthService(mockUserRepo, secret)
	_, err := authService.Authorize(user.Email, validPassword)
	assert.NoError(t, err)
}

func TestAuthServiceInvalidPassword(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockUserRepo := mocks.NewMockUserRepo(mockController)
	mockUserRepo.EXPECT().FindByEmail(user.Email).Return(&user, nil).Times(1)

	authService := NewAuthService(mockUserRepo, secret)
	_, err := authService.Authorize(user.Email, invalidPassword)
	assert.Error(t, err)
}

func TestAuthServiceInvalidEmail(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockUserRepo := mocks.NewMockUserRepo(mockController)
	mockUserRepo.EXPECT().FindByEmail(user.Email).Return(nil, nil).Times(1)

	authService := NewAuthService(mockUserRepo, secret)
	_, err := authService.Authorize(user.Email, validPassword)
	assert.Error(t, err)
}

func TestAuthServiceHashPassword(t *testing.T) {
	t.Parallel()
	authService := NewAuthService(nil, secret)
	//hash, err := authService.HashPassword(validPassword)
	_, err := authService.HashPassword(validPassword)
	assert.NoError(t, err)
	//fmt.Printf("PW: %s, HASH: %s, SECRET: %s\n", validPassword, hash, secret)
}
