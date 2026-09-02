.PHONY: api web test check build compose-up compose-down manifests

api:
	go run ./cmd/api

web:
	pnpm dev

test:
	go test ./...
	pnpm test

check:
	go vet ./...
	test -z "$$(gofmt -l cmd internal)"
	pnpm lint
	pnpm typecheck

build:
	go build ./cmd/api
	pnpm build

compose-up:
	docker compose up --build

compose-down:
	docker compose down

manifests:
	kubectl kustomize deploy/k8s/base
