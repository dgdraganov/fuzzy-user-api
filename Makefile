include dev.env
export

compose:
	docker-compose up --detach --build

tests:
	go test ./...

decompose:
	docker-compose down
