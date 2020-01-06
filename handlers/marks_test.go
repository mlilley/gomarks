package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/mocks"
	"github.com/mlilley/gomarks/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type suiteMarksHandler struct {
	suite.Suite
	secret         []byte
	user           app.User
	marks          []app.Mark
	mockController *gomock.Controller
	mockMarkRepo   *mocks.MockMarkRepo
	markService    services.MarkService
	handler        echo.HandlerFunc
	req            *http.Request
	rec            *httptest.ResponseRecorder
	ctx            echo.Context
}

func (s *suiteMarksHandler) SetupSuite() {
	s.secret = []byte("test-secret")
	s.user = app.User{ID: "1", Email: "user@test.com", PasswordHash: "$2a$14$NX0xoDuYKiibPruRYAaGce/3rOkFSesc2ZXeRWb.X/Wa45OKiMvzK", Active: true}
	s.marks = []app.Mark{
		{ID: "1", Title: "Google", URL: "http://google.com"},
		{ID: "2", Title: "Amazon", URL: "http://amazon.com"},
	}
}

func (s *suiteMarksHandler) SetupTest() {
	s.mockController = gomock.NewController(s.T())
	s.mockMarkRepo = mocks.NewMockMarkRepo(s.mockController)
	s.markService = services.NewMarkService(s.mockMarkRepo)
	s.handler = HandleGetMarks(s.markService)

	e := echo.New()
	s.req = httptest.NewRequest(http.MethodPost, "/marks", nil)
	s.rec = httptest.NewRecorder()
	s.ctx = e.NewContext(s.req, s.rec)
}

func (s *suiteMarksHandler) TestGetMarks() {
	s.mockMarkRepo.EXPECT().FindAllByUserID(s.user.ID).Return(s.marks, nil)

	err := s.handler(s.ctx)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, s.rec.Code)
	assert.Equal(s.T(), s.rec.Header().Get(echo.HeaderContentType), "application/json; charset=UTF-8")

	expected, err := json.Marshal(map[string][]app.Mark{"marks": s.marks})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), string(expected), strings.TrimSpace(s.rec.Body.String()))
}

func (s *suiteMarksHandler) TestGetMarksNone() {
	// return an empty array
	s.mockMarkRepo.EXPECT().FindAllByUserID(s.user.ID).Return([]app.Mark{}, nil)

	err := s.handler(s.ctx)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, s.rec.Code)
	assert.Equal(s.T(), s.rec.Header().Get(echo.HeaderContentType), "application/json; charset=UTF-8")

	expected, err := json.Marshal(map[string][]app.Mark{"marks": {}})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), string(expected), strings.TrimSpace(s.rec.Body.String()))
}

func (s *suiteMarksHandler) TestGetMarksError() {
	// return an empty array
	s.mockMarkRepo.EXPECT().FindAllByUserID(s.user.ID).Return(nil, fmt.Errorf("whoops"))

	err := s.handler(s.ctx)
	assert.Error(s.T(), err)
}

func (s *suiteMarksHandler) AfterTest() {
	s.mockController.Finish()
}

func TestMarksHandlerSuite(t *testing.T) {
	suite.Run(t, new(suiteMarksHandler))
}