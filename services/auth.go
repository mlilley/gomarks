package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mlilley/gomarks/repos"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Authorize(username string, password string) (string, error)
	HashPassword(password string) (string, error)
}

func NewAuthService(userRepo repos.UserRepo, secret []byte) AuthService {
	return &authService{userRepo: userRepo, secret: secret}
}

type authService struct {
	userRepo repos.UserRepo
	secret []byte
}

func (s *authService) Authorize(username string, password string) (string, error) {
	u, err := s.userRepo.FindByEmail(username)
	if err != nil {
		return "", err
	}

	if u == nil || !u.Active {
		return "", fmt.Errorf("user not found or not active")
	}

	isAuthenticated := CheckPassword(password, u.PasswordHash)
	if !isAuthenticated {
		return "", fmt.Errorf("password incorrect")
	}

	token, err := GenerateToken(s.secret, u.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) HashPassword(password string) (string, error) {
	return HashPassword(password)
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(signingKey []byte, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": email, "exp": time.Now().Add(time.Hour * 24).Unix()})

	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}