#!/bin/sh

echo "Starting application setup..."

# Wait for database to be ready
echo "Waiting for database to be ready..."
until nc -z ${DB_HOST} ${DB_PORT}; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Run migrations
echo "Running database migrations..."
migrate -path /root/migrations \
        -database "postgres://${DB_USER:-postgres}:${DB_PASSWORD:-postgres}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
        up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
else
    echo "Migration failed, but continuing anyway..."
fi

# Start the application
echo "Starting the application..."
exec "$@"