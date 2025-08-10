build:
	CGO_ENABLED=0 GOOS=linux go build -a -o gommo ./...

run:
	go run ./...

init:
	go mod gommo
	go get .

test:
	go test -v -cover

buildarm:
	CGO_ENABLED=1 GOOD=linux GOARCH=arm64 go build -a -o ./bin/arm64/gommo

docker:
	sudo docker build . -t tskal.dev/gommo:latest

testReport:
	go test -covermode=count -coverpkg=./... -coverprofile cover.out -v ./...
	go tool cover -html cover.out -o cover.html
