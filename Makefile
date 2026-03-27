.PHONY: dev build lint test docker-up docker-down clean
.PHONY: dev-ui dev-api build-ui build-api lint-ui lint-api

# ── Dev ────────────────────────────────────────────────────────────────────────
dev:
	@$(MAKE) -j2 dev-ui dev-api

dev-ui:
	cd ui && npm run dev

dev-api:
	go run ./cmd/api

# ── Build ──────────────────────────────────────────────────────────────────────
build: build-ui build-api

build-ui:
	cd ui && npm run build

build-api:
	CGO_ENABLED=1 go build -o bin/api ./cmd/api

# ── Lint / type-check ─────────────────────────────────────────────────────────
lint: lint-ui lint-api

lint-ui:
	cd ui && npm run lint

lint-api:
	go vet ./...

type-check:
	cd ui && npm run type-check

# ── Tests ──────────────────────────────────────────────────────────────────────
test:
	go test ./...

# ── Docker ────────────────────────────────────────────────────────────────────
docker-up:
	docker compose up --build

docker-down:
	docker compose down

# ── Cleanup ───────────────────────────────────────────────────────────────────
clean:
	cd ui && npm run clean
	rm -f bin/api
