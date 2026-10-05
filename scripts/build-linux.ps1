# Build script for Linux amd64 production binary
Write-Host "Compiling ALUSI Backend for Linux amd64..." -ForegroundColor Cyan

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

go build -ldflags="-s -w" -o bin/backend-alusi cmd/server/main.go

if ($LASTEXITCODE -eq 0) {
    $size = (Get-Item "bin/backend-alusi").Length / 1MB
    Write-Host ("Build successful! Binary location: bin/backend-alusi ({0:N2} MB)" -f $size) -ForegroundColor Green
} else {
    Write-Host "Build failed with exit code $LASTEXITCODE" -ForegroundColor Red
}
