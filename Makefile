build :
	go build -v

cov : all
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

check :
	golint .
	go vet -all .
	gofmt -s -l .
	goreportcard-cli -v | grep -v cyclomatic
	