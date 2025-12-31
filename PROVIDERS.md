# AI Provider Configuration Guide

This document provides detailed information about all supported AI providers in vibecheck.

## Overview

vibecheck supports 10 AI providers, each optimized for cost-efficiency and quality:

1. **OpenAI** - GPT-4o-mini
2. **Google Gemini** - gemini-2.5-flash
3. **Anthropic** - Claude 3.5 Haiku
4. **Groq** - Llama 3.3 70B Versatile
5. **xAI** - Grok Beta
6. **Moonshot AI Kimi** - moonshot-v1-auto
7. **Alibaba Qwen** - qwen-turbo
8. **DeepSeek** - deepseek-chat
9. **Perplexity** - sonar
10. **Ollama** - gpt-oss:20b (local)

---

## Provider Setup

### 1. OpenAI (GPT-4o-mini)

**Model**: `gpt-4o-mini`  
**Cost**: $0.15/1M input tokens, $0.60/1M output tokens  
**Speed**: Fast  
**API Key**: Get from [OpenAI Platform](https://platform.openai.com/api-keys)

```bash
export OPENAI_API_KEY="sk-..."
vibecheck commit --provider openai
```

---

### 2. Google Gemini (gemini-2.5-flash)

**Model**: `gemini-2.5-flash`  
**Cost**: $0.075/1M input tokens, $0.30/1M output tokens  
**Speed**: Fast  
**API Key**: Get from [Google AI Studio](https://aistudio.google.com/app/apikey)

```bash
export GEMINI_API_KEY="..."
vibecheck commit --provider gemini
```

**Features**:
- Lowest cost per token among cloud providers
- Excellent performance for commit messages
- 1M token context window

---

### 3. Anthropic (Claude 3.5 Haiku)

**Model**: `claude-3.5-haiku-20241022`  
**Cost**: $0.80/1M input tokens, $4.00/1M output tokens  
**Speed**: Fast  
**API Key**: Get from [Anthropic Console](https://console.anthropic.com/settings/keys)

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
vibecheck commit --provider anthropic
```

**Features**:
- Most affordable Claude model
- Excellent reasoning capabilities
- Strong coding understanding

---

### 4. Groq (Llama 3.3 70B Versatile)

**Model**: `llama-3.3-70b-versatile`  
**Cost**: Free tier available  
**Speed**: Ultra-fast (fastest inference)  
**API Key**: Get from [Groq Console](https://console.groq.com/keys)

```bash
export GROQ_API_KEY="gsk_..."
vibecheck commit --provider groq
```

**Features**:
- Blazing fast inference speed
- Free tier available
- High-quality open-source model
- Best for quick iterations

---

### 5. xAI (Grok Beta)

**Model**: `grok-beta`  
**Cost**: Competitive pricing  
**Speed**: Fast  
**API Key**: Get from [xAI Console](https://console.x.ai/)

```bash
export XAI_API_KEY="xai-..."
vibecheck commit --provider grok
```

**Features**:
- Access to X's training data
- Strong technical understanding
- OpenAI-compatible API

---

### 6. Moonshot AI Kimi (moonshot-v1-auto)

**Model**: `moonshot-v1-auto`  
**Cost**: Very low (competitive with Gemini)  
**Speed**: Ultra-fast  
**API Key**: Get from [Moonshot AI Platform](https://platform.moonshot.cn/)

```bash
export MOONSHOT_API_KEY="..."
vibecheck commit --provider kimi
```

**Features**:
- 1 trillion parameter MoE architecture
- Activates only 32B parameters for efficiency
- 128K context window
- Excellent for coding tasks
- Auto-selects best model variant

---

### 7. Alibaba Qwen (qwen-turbo)

**Model**: `qwen-turbo`  
**Cost**: Very low (extremely competitive)  
**Speed**: Ultra-fast  
**API Key**: Get from [Alibaba Cloud DashScope](https://dashscope.aliyun.com/)

```bash
export QWEN_API_KEY="..."
vibecheck commit --provider qwen
```

**Features**:
- Latest Qwen 3 generation model
- MoE architecture for efficiency
- 128K+ context window
- Strong multilingual capabilities
- Excellent code understanding

---

### 8. DeepSeek (deepseek-chat)

**Model**: `deepseek-chat`  
**Cost**: Extremely low ($0.14/1M input, $0.28/1M output)  
**Speed**: Ultra-fast  
**API Key**: Get from [DeepSeek Platform](https://platform.deepseek.com/)

```bash
export DEEPSEEK_API_KEY="..."
vibecheck commit --provider deepseek
```

**Features**:
- 671B parameter model
- Best cost-to-performance ratio
- Excellent coding and reasoning
- 128K context window
- OpenAI-compatible API

---

### 9. Perplexity (sonar)

**Model**: `sonar`  
**Cost**: ~$0.60/1M input tokens, ~$0.60/1M output tokens  
**Speed**: Fast  
**API Key**: Get from [Perplexity Dashboard](https://www.perplexity.ai/settings/api)

```bash
export PERPLEXITY_API_KEY="pplx-..."
vibecheck commit --provider perplexity
```

**Features**:
- Search-grounded responses with fresh context
- 128K context window
- Balanced cost and quality
- OpenAI-compatible API surface

---

### 10. Ollama (gpt-oss:20b)

**Model**: `gpt-oss:20b`  
**Cost**: Free (runs locally)  
**Speed**: Medium (depends on hardware)  
**Setup**: Install [Ollama](https://ollama.ai) and pull the model

```bash
# Install ollama and pull the model
ollama pull gpt-oss:20b

# Optional: Set custom host
export OLLAMA_HOST="http://localhost:11434"

vibecheck commit --provider ollama
```

**Features**:
- Completely free and private
- No API key required
- Runs entirely on your machine
- No data leaves your computer

---

## Cost Comparison

For a typical commit message (≈500 input tokens, ≈150 output tokens):

| Provider   | Cost per Commit | Free Tier | Speed      |
|------------|----------------|-----------|------------|
| DeepSeek   | ~$0.00011      | Yes       | Ultra-Fast |
| Gemini     | ~$0.00006      | Yes       | Fast       |
| Kimi       | ~$0.00008      | Yes       | Ultra-Fast |
| Qwen       | ~$0.00007      | Yes       | Ultra-Fast |
| OpenAI     | ~$0.00015      | No        | Fast       |
| Anthropic  | ~$0.00100      | No        | Fast       |
| Groq       | Free*          | Yes       | Ultra      |
| xAI        | ~$0.00015      | No        | Fast       |
| Perplexity | ~$0.00045      | No        | Fast       |
| Ollama     | Free           | N/A       | Medium     |

*Groq offers a generous free tier

---

## Recommended Providers

### For Cost-Conscious Users
1. **Ollama** - Completely free, runs locally
2. **Groq** - Free tier with ultra-fast speed
3. **Gemini** - Lowest cloud pricing
4. **DeepSeek** - Extremely low cost ($0.14/M tokens)
5. **Qwen** - Very competitive pricing

### For Speed
1. **DeepSeek** - Ultra-fast inference
2. **Qwen** - Ultra-fast with MoE
3. **Kimi** - Ultra-fast with auto-selection
4. **Groq** - Fastest inference speed
5. **Gemini** - Very fast

### For Quality
1. **Anthropic** - Best reasoning
2. **OpenAI** - Excellent quality
3. **DeepSeek** - Strong coding performance
4. **Kimi** - Excellent for code
5. **Gemini** - Strong performance

### For Privacy
1. **Ollama** - Runs completely locally
2. All others require API calls

---

## Environment Variables Summary

Create a `.env` file or add to your shell profile:

```bash
# OpenAI
export OPENAI_API_KEY="sk-..."

# Google Gemini
export GEMINI_API_KEY="..."

# Anthropic Claude
export ANTHROPIC_API_KEY="sk-ant-..."

# Groq
export GROQ_API_KEY="gsk_..."

# xAI Grok
export XAI_API_KEY="xai-..."

# Moonshot AI Kimi
export MOONSHOT_API_KEY="..."

# Alibaba Qwen
export QWEN_API_KEY="..."

# DeepSeek
export DEEPSEEK_API_KEY="..."

# Perplexity
export PERPLEXITY_API_KEY="pplx-..."

# Ollama (optional)
export OLLAMA_HOST="http://localhost:11434"
```

---

## Testing Your Setup

Test each provider to ensure it's configured correctly:

```bash
# Test OpenAI
vibecheck commit --provider openai

# Test Gemini
vibecheck commit --provider gemini

# Test Anthropic
vibecheck commit --provider anthropic

# Test Groq
vibecheck commit --provider groq

# Test xAI Grok
vibecheck commit --provider grok

# Test Kimi
vibecheck commit --provider kimi

# Test Qwen
vibecheck commit --provider qwen

# Test DeepSeek
vibecheck commit --provider deepseek

# Test Perplexity
vibecheck commit --provider perplexity

# Test Ollama
vibecheck commit --provider ollama
```

---

## Troubleshooting

### "API key not set" error
Make sure you've exported the correct environment variable for your provider.

### Ollama connection error
Ensure Ollama is running: `ollama serve`

### "Model not found" (Ollama)
Pull the model first: `ollama pull gpt-oss:20b`

### Rate limiting
Some providers have rate limits on free tiers. Consider:
- Using a different provider
- Upgrading to a paid tier
- Using Ollama for unlimited local usage

---

## Support

For issues or questions:
- GitHub Issues: [vibecheck issues](https://github.com/rshdhere/vibecheck/issues)
- Documentation: [README.md](./README.md)

