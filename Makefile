VERSION := "0.1.0"

go-install-latest-tag:
	go install github.com/bariiss/transfersh@$(shell git describe --tags --abbrev=0)

go-build:
	go build -o transfersh github.com/bariiss/transfersh