[![Downloads](https://img.shields.io/github/downloads/rshdhere/vibecheck/total?color=cyan&label=Downloads&logo=github)](https://github.com/rshdhere/vibecheck/releases)
[![Stars](https://img.shields.io/github/stars/rshdhere/vibecheck?color=yellow&logo=github)](https://github.com/rshdhere/vibecheck/stargazers)

# vibecheck

A cross-platform command-line ai-tool for automating git commit messages using AI models. Supports 09 providers including OpenAI, Gemini 2.5, Anthropic Claude, Groq, Grok, Kimi K2, Qwen, DeepSeek, and Ollama.

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
- Automatically request sudo privileges if needed (Linux/macOS)

> **Note:** If vibecheck is installed in a protected directory like `/usr/local/bin`, the upgrade command will automatically re-run itself with sudo to complete the installation.

## Configuration

Set up your API keys as environment variables:
> **Skip:** If you already have one of the API keys in your .env already, then it picks it up AUTOMATICALLY.

```bash
# OpenAI (GPT-4o-mini)
export OPENAI_API_KEY="your-openai-api-key"

# Google Gemini (gemini-2.5-flash)
export GEMINI_API_KEY="your-gemini-api-key"

# Anthropic Claude (claude-3.5-haiku)
export ANTHROPIC_API_KEY="your-anthropic-api-key"

# Groq (llama-3.3-70b-versatile)
export GROQ_API_KEY="your-groq-api-key"

# xAI Grok (grok-beta)
export XAI_API_KEY="your-xai-api-key"

# Moonshot AI Kimi (moonshot-v1-auto)
export MOONSHOT_API_KEY="your-moonshot-api-key"

# Alibaba Qwen (qwen-turbo)
export QWEN_API_KEY="your-qwen-api-key"

# DeepSeek (deepseek-chat)
export DEEPSEEK_API_KEY="your-deepseek-api-key"

# Ollama (local, no API key needed)
# Set OLLAMA_HOST if not using default http://localhost:11434
export OLLAMA_HOST="http://localhost:11434"
```

## Usage

```bash
# Generate and commit with AI (default: OpenAI)
vibecheck commit

# Use a specific provider
vibecheck commit --provider openai    # GPT-4o-mini
vibecheck commit --provider gemini    # Gemini 2.5 Flash
vibecheck commit --provider anthropic # Claude 3.5 Haiku
vibecheck commit --provider groq      # Llama 3.3 70B
vibecheck commit --provider grok      # Grok Beta
vibecheck commit --provider kimi      # Kimi K2 (Moonshot-v1-auto)
vibecheck commit --provider qwen      # Qwen Turbo
vibecheck commit --provider deepseek  # DeepSeek Chat
vibecheck commit --provider ollama    # gpt-oss:20b (local)

# Add custom context to the commit message
vibecheck commit --prompt "refactored authentication logic"

# Combine provider and custom context
vibecheck commit --provider gemini --prompt "fixed bug in parser"

# Check version
vibecheck --version

# Get help
vibecheck --help
```
![dashboard-cut](https://github.com/user-attachments/assets/e45d09f6-bc3a-41cf-a8aa-d26e21a04880)
![models](https://github.com/user-attachments/assets/bc496954-87e2-4487-a352-bafbb2ea70a7)


## Supported Models

All models are selected for cost-efficiency and quality comparable to GPT-4o-mini:

| Provider   | Model                    | Cost-Efficiency | Speed      |
|------------|--------------------------|-----------------|------------|
| OpenAI     | gpt-4o-mini             | High            | Fast       |
| Gemini     | gemini-2.5-flash        | Very High       | Ultra-Fast |
| Anthropic  | claude-3.5-haiku        | High            | Fast       |
| Groq       | llama-3.3-70b-versatile | Very High       | Ultra      |
| xAI        | grok-beta               | High            | Fast       |
| Kimi       | moonshot-v1-auto        | Very High       | Ultra-Fast |
| Qwen       | qwen-turbo              | Very High       | Ultra-Fast |
| DeepSeek   | deepseek-chat           | Extremely High  | Ultra-Fast |
| Ollama     | gpt-oss:20b             | Free (Local)    | Medium     |
