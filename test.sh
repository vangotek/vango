#!/bin/bash

echo "=== VanGo Static Site Generator Test ==="
echo ""

# Test 1: Check Go files compile
echo "1. Checking Go compilation..."
if go build -o test-vango main.go 2>/dev/null; then
    echo "   ✓ Go compilation successful"
    rm -f test-vango test-vango.exe
else
    echo "   ✗ Go compilation failed"
    go build main.go
    exit 1
fi

# Test 2: Check directory structure
echo ""
echo "2. Checking directory structure..."
required_dirs=("content" "layouts" "static" "internal")
for dir in "${required_dirs[@]}"; do
    if [ -d "$dir" ]; then
        echo "   ✓ $dir directory exists"
    else
        echo "   ✗ $dir directory missing"
    fi
done

# Test 3: Check required files
echo ""
echo "3. Checking required files..."
required_files=("config.toml" "content/hello.md" "content/about.md" "layouts/_default/single.html" "static/style.css")
for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "   ✓ $file exists"
    else
        echo "   ✗ $file missing"
    fi
done

# Test 4: Check Go modules
echo ""
echo "4. Checking Go modules..."
if go mod tidy 2>/dev/null; then
    echo "   ✓ Go modules are valid"
else
    echo "   ✗ Go modules have issues"
fi

# Test 5: Try building the site
echo ""
echo "5. Testing site build..."
if timeout 30s go run main.go 2>/dev/null; then
    echo "   ✓ Site build completed"
    if [ -d "public" ]; then
        echo "   ✓ Public directory created"
        ls -la public/ | head -10
    fi
else
    echo "   ✗ Site build failed or timed out"
fi

echo ""
echo "=== Test Summary ==="
echo "VanGo static site generator setup complete!"
echo ""
echo "To use VanGo:"
echo "  Build site:    go run main.go"
echo "  Dev server:    go run main.go -mode serve"
echo "  Help:          go run main.go -help"
echo ""
