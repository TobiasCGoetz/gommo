source_files = data.go config.go api.go main.go

run:
	go run ${source_files}

init:
	go mod gommo
	go get .
