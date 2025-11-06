$repo = "rshdhere/vibecheck"
$bin = "vibecheck"
$tag = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
$url = "https://github.com/$repo/releases/download/$tag/${bin}_Windows_x86_64.zip"
$temp = "$env:TEMP\$bin.zip"
$dest = "$env:ProgramFiles\$bin"

Write-Host "⬇️  Downloading $bin $tag..."
Invoke-WebRequest -Uri $url -OutFile $temp

Write-Host "📦 Installing to $dest..."
Expand-Archive -Path $temp -DestinationPath $dest -Force

$path = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($path -notlike "*$dest*") {
  setx PATH "$path;$dest" > $null
}

Write-Host "✅ Installed! Run '$bin commit' in a new terminal."
