source_files = config.go api.go main.go data.go card.go direction.go gameMap.go player.go terrain.go

run:
	go run ${source_files}

init:
	go mod gommo
	go get .

test:
	go test -cover

build:
	CGO_ENABLED=0 GOOS=linux go build -a -o gommo

docker:
	sudo docker build . -t tskal.dev/gommo:latest

testReport:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
