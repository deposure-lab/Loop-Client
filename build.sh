#!/bin/bash

mkdir -p build

echo "Linux (amd64, arm64)..."
GOOS=linux GOARCH=amd64 go build -o build/aggloop-linux-amd64 main.go
GOOS=linux GOARCH=arm64 go build -o build/aggloop-linux-arm64 main.go

echo "Windows (amd64, arm64)..."
GOOS=windows GOARCH=amd64 go build -o build/aggloop-windows-amd64.exe main.go
GOOS=windows GOARCH=arm64 go build -o build/aggloop-windows-arm64.exe main.go

echo "macOS (Intel, Apple Silicon)..."
GOOS=darwin GOARCH=amd64 go build -o build/aggloop-macos-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o build/aggloop-macos-arm64 main.go