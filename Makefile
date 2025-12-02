.PHONY: up down api next build clean migrate seed

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Run backend API
api:
	cd backend && go run cmd/api/main.go

# Run frontend
next:
	cd frontend && npm run dev

# Build backend
build-api:
	cd backend && go build -o bin/api cmd/api/main.go

# Build frontend
build-frontend:
	cd frontend && npm run build

# Run database migrations
migrate:
	cd backend && go run cmd/migrate/main.go

# Seed database
seed:
	docker-compose exec postgres psql -U myplate -d myplate -f /docker-entrypoint-initdb.d/seed.sql

# Clean build artifacts
clean:
	rm -rf backend/bin
	rm -rf frontend/.next
	rm -rf frontend/node_modules

# Install frontend dependencies
install-frontend:
	cd frontend && npm install

# Install backend dependencies
install-backend:
	cd backend && go mod download

# Run tests
test:
	cd backend && go test ./...

# Format code
fmt:
	cd backend && go fmt ./...
	cd frontend && npm run lint:fix


