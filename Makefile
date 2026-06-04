.PHONY: run migrate-up migrate-down migrate-version migrate-force

run:
	go run ./cmd/main.go

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down 1

migrate-version:
	go run ./cmd/migrate version

migrate-force:
	@test -n "$(VERSION)" || (echo "VERSION is required, e.g. make migrate-force VERSION=1" && exit 1)
	go run ./cmd/migrate force $(VERSION)
