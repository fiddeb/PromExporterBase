version = 0.0.1
GIT_CURRENT_BRANCH_NAME = $(shell git branch --show-current)


build:
	go build -o bin/main cmd/PromExporterBase/peb.go

run:
	go run cmd/PromExporterBase/peb.go

branch:
	echo $(GIT_CURRENT_BRANCH_NAME)

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -tags static_all -ldflags="-X 'github.com/fiddeb/PromExporterBase/pkg/build.Version=$(version)'" -o bin/test-amd64-linux cmd/PromExporterBase/peb.go
	GOOS=darwin GOARCH=amd64 go build -tags static_all -ldflags="-X 'github.com/fiddeb/PromExporterBase/pkg/build.Version=$(version)'" -o bin/test-amd64-darwin cmd/PromExporterBase/peb.go

all: hello build