.PHONY: test
test: 
	go test -race ./...

.PHONY: generror
generror:
	go build -o bin/generror cmd/generror/main.go

