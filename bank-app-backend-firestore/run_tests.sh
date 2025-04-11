#!/bin/bash

# This script runs all tests for the Firebase banking application

# Check if emulators are running
if ! curl -s http://localhost:8091 > /dev/null; then
    echo "Firebase emulators do not appear to be running."
    echo "Please start them with ./start_emulators.sh in a separate terminal."
    exit 1
fi

# Function to run tests with proper output
run_tests() {
    TEST_TYPE=$1
    TEST_PATH=$2
    
    echo "Running $TEST_TYPE tests..."
    echo "=============================="
    go test $TEST_PATH -v
    
    if [ $? -eq 0 ]; then
        echo -e "\n✅ $TEST_TYPE tests passed!\n"
    else
        echo -e "\n❌ $TEST_TYPE tests failed!\n"
        FAILED=1
    fi
}

# Set initial status
FAILED=0

# Run unit tests
run_tests "Unit" "./tests/unit/..."

# Run integration tests
run_tests "Integration" "./tests/integration/..."

# Run functional tests
run_tests "Functional" "./tests/functional/..."

# Final output
if [ $FAILED -eq 0 ]; then
    echo "🎉 All tests passed!"
else
    echo "❌ Some tests failed. Please check the output above."
    exit 1
fi
