.PHONY: ping run test swaggo-install swaggo 

ping:
	echo "pong"

run:
	go run ./cmd/app/main.go

test: 
	go test -v -cover -covermode=atomic ./...

swaggo-install:
	go get -u github.com/swaggo/swag/cmd/swag

	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

swaggo:
	swag init -g cmd/app/main.go --output docs