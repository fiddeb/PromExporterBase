version = 0.0.1
BRANCH_NAME = $(shell git branch --show-current)

build:
	go build -o bin/main cmd/PromExporterBase/peb.go

run:
	go run cmd/PromExporterBase/peb.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -tags static_all -ldflags="-X 'github.com/fiddeb/PromExporterBase/pkg/build.Version=$(version)'-X 'github.com/fiddeb/PromExporterBase/pkg/build.Branch=$(BRANCH_NAME)'" -o bin/test-amd64-linux cmd/PromExporterBase/peb.go
	GOOS=darwin GOARCH=amd64 go build -tags static_all -ldflags="-X 'github.com/fiddeb/PromExporterBase/pkg/build.Version=$(version)'-X 'github.com/fiddeb/PromExporterBase/pkg/build.Branch=$(BRANCH_NAME)'" -o bin/test-amd64-darwin cmd/PromExporterBase/peb.go

all: compile