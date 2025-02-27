#!/bin/sh

# Script to replace the API URL in the built Next.js files

# Get the server's IP or hostname from environment variable or use localhost as fallback
SERVER_URL=${NEXT_PUBLIC_API_URL:-http://localhost:8080}

echo "=================================================="
echo "Replacing API URL in JavaScript files"
echo "Setting API URL to: $SERVER_URL"
echo "=================================================="

# Find and replace localhost:8080 with the server URL in all JavaScript files
# Use more portable commands
find /app/.next -type f -name "*.js" -exec sh -c "sed -i.bak 's|http://localhost:8080|'$SERVER_URL'|g' {}" \;

# Also replace it in any JSON files
find /app/.next -type f -name "*.json" -exec sh -c "sed -i.bak 's|http://localhost:8080|'$SERVER_URL'|g' {}" \;

# Clean up backup files
find /app/.next -name "*.bak" -delete

echo "Done replacing API URLs"
echo "=================================================="

# For debug purposes, let's check one of the files
echo "Sample file content check:"
head -n 20 `find /app/.next -type f -name "*.js" | head -n 1`

# Start the Next.js application
exec npm start
