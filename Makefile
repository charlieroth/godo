build:
	go build -o godo cmd/cli/main.go

run:
	go run cmd/cli/main.go

test:
	go test ./...

clean:
	rm -f godo
