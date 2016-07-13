.PHONY: all
all: fmt vet lint test

.PHONY:
get-deps:
	go get github.com/golang/lint/golint

.PHONY: test
test:
	go test -cover $(go list ./... | grep -v /vendor/)

.PHONY: fmt
fmt:
	go fmt $(go list ./... | grep -v /vendor/)

.PHONY: vet
vet:
	go vet $(go list ./... | grep -v /vendor/)

.PHONY: lint
lint:
	golint $(go list ./... | grep -v /vendor/)

.PHONY: ci
ci: fmt vet test
