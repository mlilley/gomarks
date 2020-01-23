package services

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/repos"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const AccessTokenExpiry = time.Minute * 30
const RefreshTokenExpiry = time.Minute * 60 * 24 * 365 // 1 year

type AuthService interface {
	Login(email string, password string, deviceId string) (*app.User, string, string, string, error)
	Authorize(accessToken string) (*app.User, string, error)
	Refresh(refreshToken string) (*app.User, string, string, error)
	HashPassword(password string) (string, error)
}

func NewAuthService(userRepo repos.UserRepo, deviceRepo repos.DeviceRepo, secret []byte) AuthService {
	return &authService{userRepo: userRepo, deviceRepo: deviceRepo, secret: secret}
}

type authService struct {
	userRepo   repos.UserRepo
	deviceRepo repos.DeviceRepo
	secret     []byte
}

type accessTokenClaims struct {
	DeviceId string `json:"device_id"`
	jwt.StandardClaims
}

type refreshTokenClaims struct {
	DeviceId string `json:"device_id"`
	jwt.StandardClaims
}

func (s *authService) Login(email string, password string, deviceId string) (*app.User, string, string, string, error) {
	// load the specified user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("error finding user: %w", err)
	}

	// ensure the user exists and is active
	if user == nil || !user.Active {
		return nil, "", "", "", fmt.Errorf("user not found or not active")
	}

	// verify the provided password matches the stored hash
	isAuthenticated := CheckPassword(password, user.PasswordHash)
	if !isAuthenticated {
		return nil, "", "", "", fmt.Errorf("password incorrect")
	}

	// issue both access and refresh tokens
	accessToken, err := generateAccessToken(s.secret, user.ID, deviceId)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("error generating access token: %w", err)
	}
	refreshToken, err := generateRefreshToken(s.secret, user.ID, deviceId)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("error generating refresh token: %w", err)
	}

	// persist the device (with the hash of the refresh token)
	device := app.Device{
		UserId:    user.ID,
		DeviceId:  deviceId,
		TokenHash: hashToken(refreshToken),
	}
	err = s.deviceRepo.Upsert(&device)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("failed to upsert device: %w", err)
	}

	return nil, deviceId, accessToken, refreshToken, nil
}

// If accessToken is valid and refers to a user that exists and is active,
// then returns the user and deviceId, else error.
func (s *authService) Authorize(accessToken string) (*app.User, string, error) {
	claims := accessTokenClaims{}
	err := parseToken(accessToken, s.secret, &claims)
	if err != nil {
		return nil, "", err
	}

	userId := claims.Subject
	user, err := s.validateUser(userId)
	if err != nil {
		return nil, "", err
	}

	deviceId := claims.DeviceId
	return user, deviceId, nil
}

// If refresh token is valid, refers to a user that exists and is active,
// refers to a device that exists for the user,
// and has not been revoked, then issues a new access token, then stores a hash of returns the hydrated user and
func (s *authService) Refresh(refreshToken string) (*app.User, string, string, error) {
	claims := refreshTokenClaims{}
	err := parseToken(refreshToken, s.secret, &claims)
	if err != nil {
		return nil, "", "", err
	}

	userId := claims.Subject
	user, err := s.validateUser(userId)
	if err != nil {
		return nil, "", "", err
	}

	deviceId := claims.DeviceId
	device, err := s.deviceRepo.FindById(deviceId, userId)
	if err != nil {
		return nil, "", "", fmt.Errorf("error loading device: %w", err)
	}

	if device == nil {
		return nil, "", "", fmt.Errorf("device not found: %w", err)
	}

	// get hash of current refresh token and ensure it matches the stored hash
	tokenHash := hashToken(refreshToken)
	if tokenHash != device.TokenHash {
		return nil, "", "", fmt.Errorf("refresh token differs: %w", err)
	}

	// issue a new access token
	accessToken, err := generateAccessToken(s.secret, user.ID, deviceId)
	if err != nil {
		return nil, "", "", fmt.Errorf("error generating access token: %w", err)
	}

	return user, deviceId, accessToken, nil
}

func (s *authService) validateUser(userId string) (*app.User, error) {
	user, err := s.userRepo.FindByID(userId)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user '%s' not found", userId)
	}
	if !user.Active {
		return nil, fmt.Errorf("user '%s' not active", userId)
	}
	return user, nil
}

func (s *authService) HashPassword(password string) (string, error) {
	return HashPassword(password)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashToken(token string) [32]byte {
	return sha256.Sum256([]byte(token))
}

// If token is valid, returns the claims contained therein.
func parseToken(tokenStr string, secretKey []byte, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token alg '%v'", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func generateAccessToken(secretKey []byte, userId string, deviceId string) (string, error) {
	claims := accessTokenClaims{
		DeviceId:       deviceId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpiry).Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   userId,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
}

func generateRefreshToken(secretKey []byte, userId string, deviceId string) (string, error) {
	claims := refreshTokenClaims{
		DeviceId: deviceId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshTokenExpiry).Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   userId,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
}
