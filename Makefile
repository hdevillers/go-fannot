build:
	go build -o bin/swiss-count ./cmd/swiss-count/main.go
	go build -o bin/swiss-subset ./cmd/swiss-subset/main.go
	go build -o bin/swiss-create-refdb ./cmd/swiss-create-refdb/main.go

test:
	go test -v swiss/swiss_test.go