# Requires Administrator privileges
#Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

$repo = "rshdhere/vibecheck"
$bin = "vibecheck"
$installDir = "$env:ProgramFiles\$bin"

# Function to find existing installations
function Find-ExistingInstallations {
    $locations = @(
        "$env:ProgramFiles\$bin\$bin.exe",
        "$env:LOCALAPPDATA\Programs\$bin\$bin.exe",
        "$env:USERPROFILE\go\bin\$bin.exe",
        "$env:USERPROFILE\.local\bin\$bin.exe",
        "$env:GOPATH\bin\$bin.exe"
    )
    
    $found = @()
    foreach ($loc in $locations) {
        if ((Test-Path $loc) -and ($loc -ne "$installDir\$bin.exe")) {
            $found += $loc
        }
    }
    
    return $found
}

# Function to get version of a binary
function Get-BinaryVersion {
    param([string]$binaryPath)
    try {
        $version = & $binaryPath --version 2>&1 | Select-Object -First 1
        return $version
    } catch {
        return "unknown"
    }
}

# Get latest release from GitHub
Write-Host "🔍 Checking for latest release..." -ForegroundColor Blue
try {
    $release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
    $tag = $release.tag_name
} catch {
    Write-Host "❌ Failed to fetch latest release from GitHub" -ForegroundColor Red
    exit 1
}

# Check for existing installations
$existing = Find-ExistingInstallations
if ($existing.Count -gt 0) {
    Write-Host "⚠️  Found existing installation(s):" -ForegroundColor Yellow
    foreach ($loc in $existing) {
        $ver = Get-BinaryVersion $loc
        Write-Host "   $loc " -NoNewline
        Write-Host "(version: $ver)" -ForegroundColor Blue
    }
    Write-Host ""
    Write-Host "🧹 Cleaning up old installations to avoid PATH conflicts..." -ForegroundColor Yellow
    foreach ($loc in $existing) {
        try {
            Remove-Item $loc -Force -ErrorAction Stop
            Write-Host "   ✓ Removed $loc" -ForegroundColor Green
            
            # Also remove the parent directory if it's empty
            $parentDir = Split-Path $loc -Parent
            if ((Get-ChildItem $parentDir -ErrorAction SilentlyContinue).Count -eq 0) {
                Remove-Item $parentDir -Force -ErrorAction SilentlyContinue
            }
        } catch {
            Write-Host "   ⚠ Couldn't remove $loc (please remove manually)" -ForegroundColor Yellow
        }
    }
    Write-Host ""
}

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "x86_64" } else { "i386" }
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $arch = "arm64"
}

$url = "https://github.com/$repo/releases/download/$tag/${bin}_Windows_${arch}.zip"
$temp = "$env:TEMP\$bin.zip"

Write-Host "⬇️  Downloading $bin $tag for Windows/$arch..." -ForegroundColor Blue
try {
    Invoke-WebRequest -Uri $url -OutFile $temp -UseBasicParsing
} catch {
    Write-Host "❌ Failed to download $bin" -ForegroundColor Red
    Write-Host "   URL: $url" -ForegroundColor Gray
    exit 1
}

Write-Host "📦 Installing to $installDir..." -ForegroundColor Blue
# Remove old installation directory if it exists
if (Test-Path $installDir) {
    Remove-Item $installDir -Recurse -Force
}
New-Item -ItemType Directory -Path $installDir -Force | Out-Null

# Extract the zip file
try {
    Expand-Archive -Path $temp -DestinationPath $installDir -Force
    Remove-Item $temp -Force
} catch {
    Write-Host "❌ Failed to extract $bin" -ForegroundColor Red
    exit 1
}

# Add to PATH if not already present
$machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($machinePath -notlike "*$installDir*") {
    Write-Host "📌 Adding to system PATH..." -ForegroundColor Blue
    [Environment]::SetEnvironmentVariable("Path", "$machinePath;$installDir", "Machine")
    $env:Path = "$env:Path;$installDir"  # Update current session
}

Write-Host "✅ Successfully installed!" -ForegroundColor Green
Write-Host ""

# Verify installation
try {
    $installedVersion = & "$installDir\$bin.exe" --version 2>&1 | Select-Object -First 1
    Write-Host "📌 Installed version: " -NoNewline
    Write-Host "$installedVersion" -ForegroundColor Green
    Write-Host "📍 Location: " -NoNewline
    Write-Host "$installDir\$bin.exe" -ForegroundColor Blue
} catch {
    Write-Host "⚠️  Installation completed but couldn't verify version" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🚀 Run " -NoNewline
Write-Host "vibecheck --help" -ForegroundColor Green -NoNewline
Write-Host " to get started!"
Write-Host ""
Write-Host "⚠️  Note: You may need to restart your terminal or run " -NoNewline -ForegroundColor Yellow
Write-Host "refreshenv" -ForegroundColor Cyan -NoNewline
Write-Host " for PATH changes to take effect." -ForegroundColor Yellow
