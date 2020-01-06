package services

import (
	"github.com/golang/mock/gomock"
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/mocks"
	"testing"
)

func TestMarkService(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testUser := app.User{ID: "1", Email: "user1@test.com", PasswordHash: "", Active: true}
	testMarks := []app.Mark{{ID: "1", Title: "Mark1", URL: "http://go.to/mark1"}}

	mockMarkRepo := mocks.NewMockMarkRepo(mockController)
	mockMarkRepo.EXPECT().FindAllByUserID(testUser.ID).Return(testMarks, nil).Times(1)

	markService := NewMarkService(mockMarkRepo)

	marks, err := markService.GetMarksForUser(&testUser)
	if err != nil {
		t.Fail()
	}

	if !(len(marks) == len(testMarks) && marks[0] == testMarks[0]) {
		t.Fail()
	}
}
