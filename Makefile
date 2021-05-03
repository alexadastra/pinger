postgres:
	docker stop microservice-postgres || true
	docker run --rm --detach --name microservice-postgres \
		--env POSTGRES_USER=root \
		--env POSTGRES_PASSWORD=password\
		--env POSTGRES_DB=microservice-db \
		--publish 5432:5432 postgres

gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative url_service/url_service.proto

clean:
	rm url_service/*.go

run:
	go run server/main.go
