@echo off
echo === VanGo Static Site Generator Test ===
echo.

REM Test 1: Check Go files compile
echo 1. Checking Go compilation...
go build -o test-vango.exe main.go >nul 2>&1
if %errorlevel% == 0 (
    echo    ✓ Go compilation successful
    del /f test-vango.exe >nul 2>&1
) else (
    echo    ✗ Go compilation failed
    go build main.go
    exit /b 1
)

REM Test 2: Check directory structure
echo.
echo 2. Checking directory structure...
if exist "content" (echo    ✓ content directory exists) else (echo    ✗ content directory missing)
if exist "layouts" (echo    ✓ layouts directory exists) else (echo    ✗ layouts directory missing)
if exist "static" (echo    ✓ static directory exists) else (echo    ✗ static directory missing)
if exist "internal" (echo    ✓ internal directory exists) else (echo    ✗ internal directory missing)

REM Test 3: Check required files
echo.
echo 3. Checking required files...
if exist "config.toml" (echo    ✓ config.toml exists) else (echo    ✗ config.toml missing)
if exist "content\hello.md" (echo    ✓ content\hello.md exists) else (echo    ✗ content\hello.md missing)
if exist "content\about.md" (echo    ✓ content\about.md exists) else (echo    ✗ content\about.md missing)
if exist "layouts\_default\single.html" (echo    ✓ layouts\_default\single.html exists) else (echo    ✗ layouts\_default\single.html missing)
if exist "static\style.css" (echo    ✓ static\style.css exists) else (echo    ✗ static\style.css missing)

REM Test 4: Check Go modules
echo.
echo 4. Checking Go modules...
go mod tidy >nul 2>&1
if %errorlevel% == 0 (
    echo    ✓ Go modules are valid
) else (
    echo    ✗ Go modules have issues
)

REM Test 5: Try building the site
echo.
echo 5. Testing site build...
timeout /t 30 /nobreak >nul & go run main.go >nul 2>&1
if %errorlevel% == 0 (
    echo    ✓ Site build completed
    if exist "public" (
        echo    ✓ Public directory created
        dir /w public
    )
) else (
    echo    ✗ Site build failed or timed out - trying with output:
    go run main.go
)

echo.
echo === Test Summary ===
echo VanGo static site generator setup complete!
echo.
echo To use VanGo:
echo   Build site:    go run main.go
echo   Dev server:    go run main.go -mode serve
echo   Help:          go run main.go -help
echo.
pause
