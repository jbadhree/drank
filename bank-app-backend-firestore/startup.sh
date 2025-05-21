# 

#!/bin/sh

# Exit on error
set -e

# Wait for the database to be ready
echo "Waiting for database to be ready..."
for i in {1..30}; do
  if nc -z db 5432; then
    echo "Database is ready!"
    break
  fi
  echo "Waiting for database... ($i/30)"
  sleep 1
done

# Seed the database first if SEED_DB is set to true
if [ "$SEED_DB" = "true" ]; then
  echo "Seeding database..."
  ./drank-backend --seed
  echo "Database seeded successfully!"
fi

# Start the application
echo "Starting the application..."
exec ./drank-backend
