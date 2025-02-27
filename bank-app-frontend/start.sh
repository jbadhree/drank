#!/bin/sh
# Simple startup script that replaces API URL in built files

echo "==== Running API URL replacement ===="
API_URL="${NEXT_PUBLIC_API_URL}"
echo "Setting API URL to: $API_URL"

# Find JS files containing the localhost URL and replace
find /app/.next -type f -name "*.js" | xargs grep -l "http://localhost:8080" | 
while read file; do
  echo "Replacing in: $file"
  sed -i.bak "s|http://localhost:8080|$API_URL|g" "$file"
done

echo "==== Replacement complete ===="

# Start the Next.js app
exec npm start
