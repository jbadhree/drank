FROM node:18-alpine AS base

# Install dependencies only when needed
FROM base AS deps
WORKDIR /app

# Copy package.json and package-lock.json (if available)
COPY package.json package-lock.json* ./

# Install dependencies
RUN npm ci

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Next.js collects anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line to disable telemetry at build time
# ENV NEXT_TELEMETRY_DISABLED 1

# Build the Next.js application
RUN npm run build

# Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
# Uncomment the following line to disable telemetry at runtime
# ENV NEXT_TELEMETRY_DISABLED 1

# Create a public directory if it doesn't exist
RUN mkdir -p /app/public

# Copy necessary files for running the application
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Copy the API URL replacement script
COPY replace-api-url.sh /app/replace-api-url.sh
RUN chmod +x /app/replace-api-url.sh

# Expose the port the app will run on
EXPOSE 3000

# Default environment variable for backend API URL
ENV NEXT_PUBLIC_API_URL=http://localhost:8080

# Start the application using our replacement script
CMD ["/app/replace-api-url.sh"]
