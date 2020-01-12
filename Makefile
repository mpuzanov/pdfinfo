SOURCE=.
APP=pdfinfo
BINARY_DIR=bin

build :
	go build -v -o ${APP} ${SOURCE}

run :
	go run ${SOURCE} < files/in_file_linux.txt

clean:
	rm -rf ${APP}

cov : all
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

check :
	~/go/bin/golint ${SOURCE}
	go vet -all ${SOURCE}
	gofmt -s -l ${SOURCE}

release:
	./scripts/build-release.sh ${BINARY_DIR}
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_DIR}/win/${APP}.exe ${SOURCE}
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_DIR}/linux/${APP} ${SOURCE}


.PHONY: build run cov check clean