SOURCE=./cmd/pdfinfo
APP=pdfinfo
BINARY_DIR=bin
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

GO_SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)
GO_TEST_DIRS := $(shell \
	find . -name "*_test.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)	

.DEFAULT_GOAL = build 

build : lint
	go build -v -o ${APP} ${SOURCE}

run :
	go run ${SOURCE} < files/in_file_linux.txt

runwin :
	go run ${SOURCE} < files/in_file_win.txt

clean:
	rm -rf ${APP}

cov : all
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

lint :
	goimports -w ${GO_SRC_DIRS}
	golangci-lint run
	@golint -set_exit_status ${PKG_LIST}
	@go vet ${PKG_LIST}	
	@#golint ${SOURCE}
	@#go vet -all ${SOURCE}	
	@#gofmt -s -w $(GO_SRC_DIRS)

release:
	./scripts/build-release.sh ${BINARY_DIR}
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY_DIR}/win/${APP}.exe ${SOURCE}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_DIR}/linux/${APP} ${SOURCE}

dep: ## Get the dependencies
	@go get -v -d ./...
	@go get -u golang.org/x/lint

.PHONY: build run cov check clean dep