package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/mocks"
	"github.com/mlilley/gomarks/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type suiteAuthHandler struct {
	suite.Suite
	secret         []byte
	user           app.User
	mockController *gomock.Controller
	mockUserRepo   *mocks.MockUserRepo
	authService    services.AuthService
	handler        echo.HandlerFunc
}

func (s *suiteAuthHandler) SetupSuite() {
	s.secret = []byte("test-secret")
	s.user = app.User{ID: "1", Email: "user@test.com", PasswordHash: "$2a$14$NX0xoDuYKiibPruRYAaGce/3rOkFSesc2ZXeRWb.X/Wa45OKiMvzK", Active: true}
}

func (s *suiteAuthHandler) SetupTest() {
	s.mockController = gomock.NewController(s.T())
	s.mockUserRepo = mocks.NewMockUserRepo(s.mockController)
	s.authService = services.NewAuthService(s.mockUserRepo, s.secret)
	s.handler = HandleCreateToken(s.authService)
}

func (s *suiteAuthHandler) TestValidLoginJson() {
	s.mockUserRepo.EXPECT().FindByEmail(s.user.Email).Return(&s.user, nil)

	input := `{"username":"user@test.com","password":"test123"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(input))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := s.handler(ctx)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Equal(s.T(), rec.Header().Get(echo.HeaderContentType), "application/json; charset=UTF-8")
}

func (s *suiteAuthHandler) TestValidLoginForm() {
	s.mockUserRepo.EXPECT().FindByEmail(s.user.Email).Return(&s.user, nil)

	input := make(url.Values)
	input.Set("username", "user@test.com")
	input.Set("password", "test123")

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(input.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := s.handler(ctx)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Equal(s.T(), rec.Header().Get(echo.HeaderContentType), "application/json; charset=UTF-8")
}

// TODO:
//   - invalid login credentials
//   - non existent user
//   - etc

func (s *suiteAuthHandler) AfterTest() {
	s.mockController.Finish()
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(suiteAuthHandler))
}