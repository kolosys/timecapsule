param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  test              - Run tests" -ForegroundColor Yellow
    Write-Host "  test-coverage     - Run tests with coverage" -ForegroundColor Yellow
    Write-Host "  build             - Build the demo" -ForegroundColor Yellow
    Write-Host "  demo              - Run the demo" -ForegroundColor Yellow
    Write-Host "  clean             - Clean build artifacts" -ForegroundColor Yellow
    Write-Host "  format            - Format code" -ForegroundColor Yellow
    Write-Host "  deps              - Install dependencies" -ForegroundColor Yellow
    Write-Host "  lint              - Run linter" -ForegroundColor Yellow
    Write-Host "  security          - Run security scan" -ForegroundColor Yellow
    Write-Host "  help              - Show this help" -ForegroundColor Yellow
}

function Invoke-Test {
    Write-Host "Running tests..." -ForegroundColor Green
    go test -v ./...
}

function Invoke-TestCoverage {
    Write-Host "Running tests with coverage..." -ForegroundColor Green
    go test -v -cover ./...
}

function Invoke-Build {
    Write-Host "Building demo..." -ForegroundColor Green
    if (!(Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin"
    }
    go build -o bin/demo.exe cmd/demo/main.go
    Write-Host "Demo built: bin/demo.exe" -ForegroundColor Green
}

function Invoke-Demo {
    Write-Host "Running demo..." -ForegroundColor Green
    go run cmd/demo/main.go
}

function Invoke-Clean {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Green
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force "bin"
    }
    if (Test-Path "coverage.out") {
        Remove-Item "coverage.out"
    }
    if (Test-Path "coverage.html") {
        Remove-Item "coverage.html"
    }
    Write-Host "Clean complete" -ForegroundColor Green
}

function Invoke-Format {
    Write-Host "Formatting code..." -ForegroundColor Green
    go fmt ./...
}

function Invoke-Deps {
    Write-Host "Installing dependencies..." -ForegroundColor Green
    go mod tidy
    go mod download
}

function Invoke-Lint {
    Write-Host "Running linter..." -ForegroundColor Green
    # Check if golangci-lint is installed
    try {
        golangci-lint run
    }
    catch {
        Write-Host "golangci-lint not found. Installing..." -ForegroundColor Yellow
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        golangci-lint run
    }
}

function Invoke-Security {
    Write-Host "Running security scan..." -ForegroundColor Green
    # Check if gosec is installed
    try {
        gosec ./...
    }
    catch {
        Write-Host "gosec not found. Installing..." -ForegroundColor Yellow
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
        # Add GOPATH/bin to PATH for this session
        $env:PATH += ";$env:GOPATH\bin"
        try {
            gosec ./...
        }
        catch {
            Write-Host "Failed to run gosec after installation. Please restart your terminal." -ForegroundColor Red
        }
    }
}

# Main script logic
switch ($Command.ToLower()) {
    "test" { Invoke-Test }
    "test-coverage" { Invoke-TestCoverage }
    "build" { Invoke-Build }
    "demo" { Invoke-Demo }
    "clean" { Invoke-Clean }
    "format" { Invoke-Format }
    "deps" { Invoke-Deps }
    "lint" { Invoke-Lint }
    "security" { Invoke-Security }
    "help" { Show-Help }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Show-Help
    }
}
