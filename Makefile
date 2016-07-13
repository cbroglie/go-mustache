.PHONY: all
all: fmt vet lint test

.PHONY:
get-deps:
	go get github.com/golang/lint/golint

.PHONY: test
test:
	GO15VENDOREXPERIMENT=1 go test -cover $(go list ./... | grep -v /vendor/)

.PHONY: fmt
fmt:
	GO15VENDOREXPERIMENT=1 go fmt $(go list ./... | grep -v /vendor/)

.PHONY: vet
vet:
	GO15VENDOREXPERIMENT=1 go vet $(go list ./... | grep -v /vendor/)

.PHONY: lint
lint:
	GO15VENDOREXPERIMENT=1 golint $(go list ./... | grep -v /vendor/)

.PHONY: ci
ci: fmt vet test
