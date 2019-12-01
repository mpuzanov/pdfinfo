SOURCE=./cmd/pdfinfo
BINARY=pdfinfo
BINARY_DIR=bin

build :
	go build -v -o ${BINARY} ${SOURCE}

run : build
	./${BINARY} < files/in_file_linux.txt
	#go run ${SOURCE}

cov : all
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

check :
	~/go/bin/golint ${SOURCE}
	go vet -all ${SOURCE}
	gofmt -s -l ${SOURCE}

release:
	./scripts/build-release.sh ${BINARY_DIR}
	#GOOS=windows GOARCH=386 go build -o ${BINARY_DIR}/win/$(BINARY).exe ${SOURCE}
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_DIR}/win/${BINARY}.exe ${SOURCE}
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_DIR}/linux/${BINARY} ${SOURCE}


.PHONY: build run cov check