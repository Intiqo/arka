test-cover:
	go test `go list ./... | grep -v cmd` -coverprofile=coverage.out
	go tool cover -html=coverage.out
