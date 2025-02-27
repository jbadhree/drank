#!/bin/sh

# Get the server's IP or hostname from environment variable or use localhost as fallback
SERVER_URL=${NEXT_PUBLIC_API_URL:-http://localhost:8080}

echo "=================================================="
echo "SIMPLE VERSION: Replacing API URL in JavaScript files"
echo "Setting API URL to: $SERVER_URL"
echo "=================================================="

# Use grep to find files containing the localhost URL
for JS_FILE in $(grep -l "http://localhost:8080" /app/.next/static/chunks/pages/*.js); do
  echo "Processing file: $JS_FILE"
  # Use simple sed command with backup
  sed -i.bak "s|http://localhost:8080|$SERVER_URL|g" "$JS_FILE"
done

echo "Done replacing API URLs"
echo "=================================================="

# Start the Next.js application
exec npm start
