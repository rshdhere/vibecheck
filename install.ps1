# Requires Administrator privileges
#Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

# --- TRACK INSTALL ---
iwr https://install.raashed.xyz/track/windows -UseBasicParsing | Out-Null

$repo = "rshdhere/vibecheck"
$bin = "vibecheck"
$installDir = "$env:ProgramFiles\$bin"

function Show-Ascii {
@"
██╗   ██╗██╗██████╗ ███████╗ ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗
██║   ██║██║██╔══██╗██╔════╝██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝
██║   ██║██║██████╔╝█████╗  ██║     ███████║█████╗  ██║     █████╔╝ 
╚██╗ ██╔╝██║██╔══██╗██╔══╝  ██║     ██╔══██║██╔══╝  ██║     ██╔═██╗ 
 ╚████╔╝ ██║██████╔╝███████╗╚██████╗██║  ██║███████╗╚██████╗██║  ██╗
  ╚═══╝  ╚═╝╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝
"@
}

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

function Get-BinaryVersion {
    param([string]$binaryPath)
    try {
        $version = & $binaryPath --version 2>&1 | Select-Object -First 1
        return $version
    } catch {
        return "unknown"
    }
}

Write-Host "🔍 Checking for latest release..." -ForegroundColor Blue
try {
    $release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
    $tag = $release.tag_name
} catch {
    Write-Host "❌ Failed to fetch latest release from GitHub" -ForegroundColor Red
    exit 1
}

$existing = Find-ExistingInstallations
if ($existing.Count -gt 0) {
    Write-Host "⚠️  Found existing installation(s):" -ForegroundColor Yellow
    foreach ($loc in $existing) {
        $ver = Get-BinaryVersion $loc
        Write-Host "   $loc " -NoNewline
        Write-Host "(version: $ver)" -ForegroundColor Blue
    }
    Write-Host ""
    Write-Host "🧹 Cleaning up old installations..." -ForegroundColor Yellow
    foreach ($loc in $existing) {
        try {
            Remove-Item $loc -Force
            Write-Host "   ✓ Removed $loc" -ForegroundColor Green
        } catch {
            Write-Host "   ⚠ Couldn't remove $loc" -ForegroundColor Yellow
        }
    }
    Write-Host ""
}

$arch = if ([Environment]::Is64BitOperatingSystem) { "x86_64" } else { "i386" }
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $arch = "arm64"
}

$url = "https://github.com/$repo/releases/download/$tag/${bin}_Windows_${arch}.zip"
$temp = "$env:TEMP\$bin.zip"

Write-Host "⬇️  Downloading $bin $tag..." -ForegroundColor Blue
Invoke-WebRequest -Uri $url -OutFile $temp -UseBasicParsing

Write-Host "📦 Installing to $installDir..." -ForegroundColor Blue
if (Test-Path $installDir) { Remove-Item $installDir -Recurse -Force }
New-Item -ItemType Directory -Path $installDir -Force | Out-Null

Expand-Archive -Path $temp -DestinationPath $installDir -Force
Remove-Item $temp -Force

$machinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($machinePath -notlike "*$installDir*") {
    Write-Host "📌 Adding to PATH..." -ForegroundColor Blue
    [Environment]::SetEnvironmentVariable("Path", "$machinePath;$installDir", "Machine")
}

Write-Host "✅ Successfully installed!" -ForegroundColor Green
try {
    $installedVersion = & "$installDir\$bin.exe" --version 2>&1 | Select-Object -First 1
    Write-Host "📌 Installed version: $installedVersion" -ForegroundColor Green
} catch {}

Write-Host "🚀 Run vibecheck --help to get started!" -ForegroundColor Green
Write-Host ""
Show-Ascii | ForEach-Object { Write-Host $_ }
Write-Host ""
Write-Host "vibecheck is a lightweight, cross-platform command line AI-tool which`nautomatically generates meaningful and consistent Git Commit Messages`nby analyzing your code changes — ship faster with vibecheck"
