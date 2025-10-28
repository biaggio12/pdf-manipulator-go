#!/bin/bash

# Test script for PDF Manipulator API
BASE_URL="http://localhost:8080"

echo "Testing PDF Manipulator API..."

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health" | jq .

echo -e "\n2. Testing convert endpoint (single page)..."
# Note: You'll need to provide actual image files for testing
echo "To test convert endpoint, use:"
echo "curl -X POST -F \"file=@test.jpg\" -F \"multiple=false\" $BASE_URL/convert --output result.pdf"

echo -e "\n3. Testing extract endpoint..."
echo "To test extract endpoint, use:"
echo "curl -X POST -F \"file=@test.pdf\" -F \"pages=1,3,5\" $BASE_URL/extract --output extracted.pdf"

echo -e "\n4. Testing merge endpoint..."
echo "To test merge endpoint, use:"
echo "curl -X POST -F \"files=@file1.pdf\" -F \"files=@file2.pdf\" $BASE_URL/merge --output merged.pdf"

echo -e "\nAPI is ready for testing!"
