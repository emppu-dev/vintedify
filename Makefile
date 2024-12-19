build:
	@go build -o bin/vintedify main.go

run: build
	@./bin/vintedify