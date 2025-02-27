#!/bin/sh
# Simple startup script that replaces API URL in built files

echo "==== Running API URL replacement ===="
API_URL="${NEXT_PUBLIC_API_URL}"
echo "Setting API URL to: $API_URL"

# Find JS files containing the localhost URL and replace
# Use a more compatible approach that works on all platforms
find /app/.next -type f -name "*.js" | while read file; do
  if grep -q "http://localhost:8080" "$file"; then
    echo "Replacing in: $file"
    # Use POSIX-compatible sed syntax that works on all platforms
    sed "s|http://localhost:8080|$API_URL|g" "$file" > "$file.tmp" && mv "$file.tmp" "$file"
  fi
done

echo "==== Replacement complete ===="

# Start the Next.js app
exec npm start
