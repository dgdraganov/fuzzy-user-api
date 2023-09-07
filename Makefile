include .env
export

compose:
	docker-compose up --detach --build

test:
	go test -v ./...

decompose:
	docker-compose down
