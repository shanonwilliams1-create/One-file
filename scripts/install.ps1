# OmniConfig — PowerShell installation script for Windows
# Usage: irm https://omnicofig.sh/install.ps1 | iex

param(
    [string]$Version = "latest",
    [string]$InstallDir = "$env:ProgramFiles\OmniConfig"
)

$Binary = "omniconfig"
$Repo = "omnicofig/cli"

# --- Detect architecture ---
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64"  { return "amd64" }
        "ARM64"  { return "arm64" }
        default  { throw "Unsupported architecture: $arch" }
    }
}

# --- Main installation ---
function Install-OmniConfig {
    Write-Host "==> OmniConfig Installer" -ForegroundColor Cyan
    Write-Host "    Platform: windows-$($(Get-Architecture))"
    Write-Host "    Target:   $InstallDir"

    # Determine download URL
    $arch = Get-Architecture
    if ($Version -eq "latest") {
        $downloadUrl = "https://github.com/$Repo/releases/latest/download/$Binary-windows-$arch.exe"
    } else {
        $downloadUrl = "https://github.com/$Repo/releases/download/$Version/$Binary-windows-$arch.exe"
    }

    Write-Host "    Download: $downloadUrl"

    # Create install directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $outputPath = Join-Path $InstallDir "$Binary.exe"

    # Download binary
    try {
        [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
        Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath -UseBasicParsing
    }
    catch {
        Write-Error "Download failed: $_"
        exit 1
    }

    # Verify download
    if (-not (Test-Path $outputPath) -or (Get-Item $outputPath).Length -eq 0) {
        Write-Error "Download failed: file is empty or missing"
        exit 1
    }

    # Add to PATH if not already there
    $currentPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine)
    if ($currentPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", [EnvironmentVariableTarget]::Machine)
        Write-Host "    Added $InstallDir to system PATH" -ForegroundColor Yellow
        Write-Host "    Restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host "==> OmniConfig installed successfully!" -ForegroundColor Green
    Write-Host "    $outputPath"
    Write-Host ""
    Write-Host "    Run '$Binary --help' to get started"
    Write-Host "    Run '$Binary detect' to detect your OS and config format"
}

Install-OmniConfig