#!/bin/bash

# Build the compiler
echo "Building ts2go compiler..."
go build -o ts2go cmd/ts2go/main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
echo "Build successful!"
echo ""

# Test all examples
examples=("hello" "calculator" "fibonacci" "control-flow" "template-literals" "arrow-functions")

for example in "${examples[@]}"; do
    echo "Testing $example..."
    ./ts2go examples/$example.ts examples/$example.go
    if [ $? -ne 0 ]; then
        echo "Compilation failed for $example!"
        exit 1
    fi

    echo "Running $example.go:"
    go run examples/$example.go
    if [ $? -ne 0 ]; then
        echo "Execution failed for $example!"
        exit 1
    fi
    echo ""
done

echo "All tests passed!"
