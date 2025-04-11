#!/bin/bash

# This script checks if the codebase can be built successfully

cd "$(dirname "$0")"
go build -v
if [ $? -eq 0 ]; then
    echo "Build successful!"
else
    echo "Build failed!"
fi
