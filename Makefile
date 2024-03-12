GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOTIDY=$(GOMOD) tidy
GOTOOL=$(GOCMD) tool

.PHONY: tidy
tidy:
	$(GOTIDY)

.PHONY: test
test: tidy
	$(GOTEST) -coverprofile cover.out -v ./...

.PHONY: cover
cover: test
	$(GOTOOL) cover -html=cover.out