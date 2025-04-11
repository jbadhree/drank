#!/bin/bash

# This script builds and runs the application

# Ensure environment variables are set
if [ ! -f .env ]; then
    echo "Creating default .env file..."
    echo "FIREBASE_PROJECT_ID=drank-firebase" > .env
    echo "FIREBASE_AUTH_EMULATOR_HOST=localhost:9099" >> .env
    echo "FIRESTORE_EMULATOR_HOST=localhost:8091" >> .env
    echo "PORT=8080" >> .env
    echo "JWT_SECRET=your-very-secret-jwt-key-change-in-production" >> .env
    echo "Created default .env file. Please review and update as needed."
fi

# Check if emulators are running
if ! curl -s http://localhost:8091 > /dev/null; then
    echo "Firebase emulators do not appear to be running."
    echo "Please start them with ./start_emulators.sh in a separate terminal."
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Build the application
echo "Building the application..."
go build -o bank-app-backend-firestore

# Run the application
echo "Starting the application..."
./bank-app-backend-firestore
