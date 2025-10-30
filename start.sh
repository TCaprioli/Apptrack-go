#!/bin/sh

# Run migrations
migrate -path db/migration -database "postgresql://postgres:postgres@apptrack-db:5432/apptrack?sslmode=disable" up

# Start the application
./main