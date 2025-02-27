#!/bin/sh

# Create a runtime configuration file with the environment variables
cat > /app/.env.local << EOL
NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL:-http://localhost:8080}
EOL

echo "-------------------------------------"
echo "Runtime environment configuration:"
echo "NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL:-http://localhost:8080}"
echo "-------------------------------------"

# Start the Next.js application
exec npm start
