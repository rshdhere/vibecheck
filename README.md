[![Downloads](https://img.shields.io/github/downloads/rshdhere/vibecheck/total?color=cyan&label=Downloads&logo=github)](https://github.com/rshdhere/vibecheck/releases)
[![Stars](https://img.shields.io/github/stars/rshdhere/vibecheck?color=yellow&logo=github)](https://github.com/rshdhere/vibecheck/stargazers)

# vibecheck

A command-line tool for automating git commit messages using AI models like GPT-4o-mini and gpt-oss:20b (using ollama).

## Installation

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/rshdhere/vibecheck/main/install.sh | bash
```

### Windows (PowerShell)

Run PowerShell as administrator, then execute:

```powershell
iwr https://raw.githubusercontent.com/rshdhere/vibecheck/main/install.ps1 | iex
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/rshdhere/vibecheck.git
cd vibecheck

# Build (version is automatically detected from git tags)
make build

# Or install to $GOPATH/bin
make install

# Or run directly
make run ARGS="--version"
```

> **Note:** The install scripts automatically detect and remove old installations to prevent PATH conflicts.

## Upgrading

Keep vibecheck up to date with a single command:

```bash
vibecheck upgrade
```

This will:
- Check for the latest release from GitHub
- Download and install the new version automatically
- Preserve your configuration

## Usage

```bash
# Generate and commit with AI
vibecheck commit

# Use a specific provider
vibecheck commit --provider ollama

# Add custom context to the commit message
vibecheck commit --prompt "refactored authentication logic"

# Check version
vibecheck --version

# Get help
vibecheck --help
```
