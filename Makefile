db-create:
	go run ./cmd/db db-create

db-drop:
	go run ./cmd/db db-drop

migration-create:
	go run ./cmd/db migration-create $(filter-out $@,$(MAKECMDGOALS))

migration-up:
	go run ./cmd/db migration-up

migration-down:
	go run ./cmd/db migration-down

migration-status:
	go run ./cmd/db migration-status

generate-sql:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

db-reset:
	go run ./cmd/db db-drop 
	go run ./cmd/db db-create
	go run ./cmd/db migration-up 

swagger-generate:
	swag init -g ../../cmd/wasted/main.go -d ./internal/auth --parseDependency --parseInternal
swagger-format:
	swag fmt

server-start:
	go run ./cmd/wasted/main.go APP_ENV=development