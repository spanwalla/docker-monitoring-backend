include .env
export

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### remove docker volume
	docker volume rm postgres_data
.PHONY: docker-rm-volume

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations 'docker_monitoring'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:{$(POSTGRES_PORT)}/$(POSTGRES_DB)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### migration down
	echo "y" | migrate -path migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:{$(POSTGRES_PORT)}/$(POSTGRES_DB)?sslmode=disable' down
.PHONY: migrate-down

swag: ### generate swagger docs
	swag init -g 'internal/app/app.go' --parseInternal --parseDependency
.PHONY: swag

test: ### run test
	go test -v ./...
.PHONY: test

mockgen: ### generate mock
	mockgen -source='internal/service/service.go'       -destination='internal/service/mocks/mock.go'    -package=service_mocks
	mockgen -source='internal/repository/repository.go' -destination='internal/repository/mocks/mock.go' -package=repository_mocks
	mockgen -source='internal/broker/broker.go'         -destination='internal/broker/mocks/mock.go'     -package=broker_mocks
.PHONY: mockgen

bin-deps: ### install binary dependencies
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
.PHONY: bin-deps