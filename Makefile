source_files = config.go api.go main.go data.go card.go direction.go gameMap.go player.go terrain.go tile.go tileWorker.go

run:
	go run ${source_files}

init:
	go mod gommo
	go get .

test:
	go test -cover

build:
	CGO_ENABLED=0 GOOS=linux go build -a -o gommo

buildarm:
	CGO_ENABLED=1 GOOD=linux GOARCH=arm64 go build -a -o ./bin/arm64/gommo

docker:
	sudo docker build . -t tskal.dev/gommo:latest

testReport:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
