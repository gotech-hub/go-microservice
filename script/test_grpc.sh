#!/bin/bash

# Script để test gRPC server
set -e

echo "Testing gRPC server..."

# Check if server is running
if ! nc -z localhost 9090; then
    echo "❌ gRPC server is not running on port 9090"
    echo "Please start the server first: make run"
    exit 1
fi

echo "✅ gRPC server is running on port 9090"

# Test with grpcurl if available
if command -v grpcurl &> /dev/null; then
    echo ""
    echo "Testing with grpcurl..."
    
    # List services
    echo "📋 Available services:"
    grpcurl -plaintext localhost:9090 list
    
    echo ""
    echo "📋 UserService methods:"
    grpcurl -plaintext localhost:9090 list user.UserService
    
    echo ""
    echo "🏥 Testing HealthCheck:"
    grpcurl -plaintext localhost:9090 user.UserService/HealthCheck
    
    echo ""
    echo "👤 Testing CreateUser:"
    grpcurl -plaintext -d '{"name": "Test User", "email": "test@example.com", "phone": "+1234567890"}' localhost:9090 user.UserService/CreateUser
    
    echo ""
    echo "📋 Testing ListUsers:"
    grpcurl -plaintext -d '{"page": 1, "limit": 5}' localhost:9090 user.UserService/ListUsers
    
    echo ""
    echo "✅ All grpcurl tests completed"
else
    echo "⚠️  grpcurl not found. Install it with:"
    echo "   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
fi

# Test with Go client if available
if [ -f "examples/grpc_client/main.go" ]; then
    echo ""
    echo "Testing with Go client..."
    
    # Build client
    go build -o bin/grpc-client examples/grpc_client/main.go
    
    # Run client
    ./bin/grpc-client
    
    echo ""
    echo "✅ Go client test completed"
else
    echo "⚠️  Go client example not found"
fi

echo ""
echo "🎉 gRPC testing completed!" 