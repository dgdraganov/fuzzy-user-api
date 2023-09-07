include .env
export

compose:
	docker-compose up --detach --build

tests:
	go test -v ./...

decompose:
	docker-compose down
