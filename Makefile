#!make

start:
	cd cmd/coinbase && go run .

tests:
	cd cmd/coinbase && go test -v -cover
	cd internal/indicator && go test -v -cover
