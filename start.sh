#!/bin/sh

# Run migrations
migrate -path db/migration \
  -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_NAME}?sslmode=disable" up

# Start the application
./main
