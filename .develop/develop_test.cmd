@echo off

REM Function to check if a command exists
:command_exists
    where %1 >nul 2>nul
    if errorlevel 1 (
        echo Error: %1 is not installed.
        exit /b 1
    )
    exit /b 0

REM Check if required commands exist
set REQUIRED_COMMANDS=go golangci-lint

for %%c in (%REQUIRED_COMMANDS%) do (
    call :command_exists %%c
    if errorlevel 1 exit /b 1
)

REM Navigate to parent directory
cd ..

REM Run go mod tidy
go mod tidy

REM Run golangci-lint
golangci-lint run -v --timeout=10m --fix

REM Run go tests with race detector and coverage
go test -race -v .\... -coverprofile .\coverage.txt

REM Generate HTML coverage report
go tool cover -html=.\coverage.txt