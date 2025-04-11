#!/bin/bash

# This script starts the Firebase emulators for local development

# Check if Firebase CLI is installed
if ! command -v firebase &> /dev/null; then
    echo "Firebase CLI is not installed. Please install it with:"
    echo "npm install -g firebase-tools"
    exit 1
fi

# Go to the parent directory where firebase.json is located
cd ..

# Start the emulators
echo "Starting Firebase emulators..."
firebase emulators:start
