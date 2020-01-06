.PHONY: test mocks
PKGS := $(shell go list ./... | grep -v /vendor)

start:
	go run main.go
test:
	go test $(PKGS)

mocks:
	rm -rf mocks/*.go
	mockgen -destination=mocks/mock_user_repo.go -package=mocks github.com/mlilley/gomarks/repos UserRepo
	mockgen -destination=mocks/mock_mark_repo.go -package=mocks github.com/mlilley/gomarks/repos MarkRepo

