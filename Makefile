.PHONY: test
test: 
	go test ./...

.PHONY: generror
generror:
	go build -o bin/generror cmd/generror/main.go

