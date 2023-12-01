## help: Print this help message.
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## migrate/create: Create a migration file. Pass the filename as name=<filename>. Example - make migrate/create name=create_user_table
migrate/create:
	goose -dir="./migrations" create $(name) sql

## fix/fieldalignment: Rearrange the struct fields to optimize memory. Automatically fix possible optimizations.
fix/fieldalignment:
	fieldalignment -fix ./...

## tests: Run all test cases with race detector and coverage enabled.
tests:
	go clean -testcache && go test -race -cover -p 1 ./...

## run/client: Run client with race detector enabled.
run/client:
	go run -race ./cmd/client

## run/server: Run server with race detector enabled.
run/server:
	go run -race ./cmd/server