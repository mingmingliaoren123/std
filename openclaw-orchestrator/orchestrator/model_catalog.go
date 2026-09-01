package orchestrator

import "strings"

const ModelCatalogVersion = "openclaw-2026.7.1-2"

// fixedModelCatalog is generated from `openclaw models list --all --json` for OpenClaw 2026.7.1-2.
// Update this catalog together with the pinned OpenClaw version.
var fixedModelCatalog = []Model{
	{Key: "deepseek/deepseek-chat", Name: "DeepSeek Chat", Input: "text", ContextWindow: 131072},
	{Key: "deepseek/deepseek-reasoner", Name: "DeepSeek Reasoner", Input: "text", ContextWindow: 131072},
	{Key: "deepseek/deepseek-v4-flash", Name: "DeepSeek V4 Flash", Input: "text", ContextWindow: 1000000, Tags: []string{"alias:DeepSeek"}},
	{Key: "deepseek/deepseek-v4-pro", Name: "DeepSeek V4 Pro", Input: "text", ContextWindow: 1000000},
	{Key: "minimax/MiniMax-M2.7", Name: "MiniMax M2.7", Input: "text", ContextWindow: 204800},
	{Key: "minimax/MiniMax-M2.7-highspeed", Name: "MiniMax M2.7 Highspeed", Input: "text", ContextWindow: 204800},
	{Key: "minimax/MiniMax-M3", Name: "MiniMax M3", Input: "text+image", ContextWindow: 1000000},
	{Key: "anthropic/claude-fable-5", Name: "Claude Fable 5", Input: "text+image", ContextWindow: 1000000},
	{Key: "anthropic/claude-haiku-4-5", Name: "Claude Haiku 4.5", Input: "text+image", ContextWindow: 200000},
	{Key: "anthropic/claude-haiku-4-5-20251001", Name: "Claude Haiku 4.5", Input: "text+image", ContextWindow: 200000},
	{Key: "anthropic/claude-mythos-5", Name: "Claude Mythos 5", Input: "text+image", ContextWindow: 1000000},
	{Key: "anthropic/claude-opus-4-6", Name: "Claude Opus 4.6", Input: "text+image", ContextWindow: 200000},
	{Key: "anthropic/claude-opus-4-7", Name: "Claude Opus 4.7", Input: "text+image", ContextWindow: 200000},
	{Key: "anthropic/claude-opus-4-8", Name: "Claude Opus 4.8", Input: "text+image", ContextWindow: 1048576},
	{Key: "anthropic/claude-sonnet-4-6", Name: "Claude Sonnet 4.6", Input: "text+image", ContextWindow: 200000},
	{Key: "anthropic/claude-sonnet-5", Name: "Claude Sonnet 5", Input: "text+image", ContextWindow: 1000000},
	{Key: "byteplus/glm-4-7-251222", Name: "GLM 4.7", Input: "text+image", ContextWindow: 200000},
	{Key: "byteplus/kimi-k2-5-260127", Name: "Kimi K2.5", Input: "text+image", ContextWindow: 256000},
	{Key: "byteplus/seed-1-8-251228", Name: "Seed 1.8", Input: "text+image", ContextWindow: 256000},
	{Key: "byteplus-plan/ark-code-latest", Name: "Ark Coding Plan", Input: "text", ContextWindow: 256000},
	{Key: "byteplus-plan/doubao-seed-code", Name: "Doubao Seed Code", Input: "text", ContextWindow: 256000},
	{Key: "byteplus-plan/glm-4.7", Name: "GLM 4.7 Coding", Input: "text", ContextWindow: 200000},
	{Key: "byteplus-plan/kimi-k2-thinking", Name: "Kimi K2 Thinking", Input: "text", ContextWindow: 256000},
	{Key: "byteplus-plan/kimi-k2.5", Name: "Kimi K2.5 Coding", Input: "text", ContextWindow: 256000},
	{Key: "claude-cli/claude-opus-4-6", Name: "Claude Opus 4.6 (Claude CLI)", Input: "text+image", ContextWindow: 200000},
	{Key: "claude-cli/claude-opus-4-7", Name: "Claude Opus 4.7 (Claude CLI)", Input: "text+image", ContextWindow: 200000},
	{Key: "claude-cli/claude-opus-4-8", Name: "Claude Opus 4.8 (Claude CLI)", Input: "text+image", ContextWindow: 1048576},
	{Key: "claude-cli/claude-sonnet-4-6", Name: "Claude Sonnet 4.6 (Claude CLI)", Input: "text+image", ContextWindow: 200000},
	{Key: "claude-cli/claude-sonnet-5", Name: "Claude Sonnet 5 (Claude CLI)", Input: "text+image", ContextWindow: 1000000},
	{Key: "cohere/command-a-03-2025", Name: "Command A", Input: "text", ContextWindow: 256000},
	{Key: "github-copilot/claude-opus-4.6", Name: "Claude Opus 4.6", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/claude-opus-4.7", Name: "Claude Opus 4.7", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/claude-opus-4.8", Name: "Claude Opus 4.8", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/claude-sonnet-4.6", Name: "Claude Sonnet 4.6", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/gemini-2.5-pro", Name: "Gemini 2.5 Pro", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/gemini-3-flash", Name: "Gemini 3 Flash", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/gemini-3.1-pro", Name: "Gemini 3.1 Pro", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/goldeneye", Name: "Goldeneye", Input: "text", ContextWindow: 128000},
	{Key: "github-copilot/gpt-5.3-codex", Name: "GPT-5.3-Codex", Input: "text", ContextWindow: 128000},
	{Key: "github-copilot/gpt-5.4", Name: "GPT-5.4", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/gpt-5.4-mini", Name: "GPT-5.4 mini", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/gpt-5.4-nano", Name: "GPT-5.4 nano", Input: "text+image", ContextWindow: 128000},
	{Key: "github-copilot/gpt-5.5", Name: "GPT-5.5", Input: "text+image", ContextWindow: 400000},
	{Key: "github-copilot/raptor-mini", Name: "Raptor mini", Input: "text", ContextWindow: 128000},
	{Key: "meta/muse-spark-1.1", Name: "Muse Spark 1.1", Input: "text+image", ContextWindow: 1048576},
	{Key: "mistral/codestral-latest", Name: "Codestral (latest)", Input: "text", ContextWindow: 256000},
	{Key: "mistral/devstral-medium-latest", Name: "Devstral 2 (latest)", Input: "text", ContextWindow: 262144},
	{Key: "mistral/magistral-small", Name: "Magistral Small", Input: "text", ContextWindow: 128000},
	{Key: "mistral/mistral-large-latest", Name: "Mistral Large (latest)", Input: "text+image", ContextWindow: 262144},
	{Key: "mistral/mistral-medium-2508", Name: "Mistral Medium 3.1", Input: "text+image", ContextWindow: 262144},
	{Key: "mistral/mistral-medium-3-5", Name: "Mistral Medium 3.5", Input: "text+image", ContextWindow: 262144},
	{Key: "mistral/mistral-small-latest", Name: "Mistral Small (latest)", Input: "text+image", ContextWindow: 128000},
	{Key: "mistral/pixtral-large-latest", Name: "Pixtral Large (latest)", Input: "text+image", ContextWindow: 128000},
	{Key: "novita/deepseek/deepseek-r1-0528", Name: "DeepSeek R1 0528", Input: "text", ContextWindow: 163840},
	{Key: "novita/deepseek/deepseek-v3-0324", Name: "DeepSeek V3 0324", Input: "text", ContextWindow: 163840},
	{Key: "novita/minimax/minimax-m2.7", Name: "MiniMax M2.7", Input: "text", ContextWindow: 1000000},
	{Key: "novita/moonshotai/kimi-k2.5", Name: "Kimi K2.5", Input: "text+image", ContextWindow: 262144},
	{Key: "novita/qwen/qwen3-235b-a22b-fp8", Name: "Qwen3 235B A22B FP8", Input: "text", ContextWindow: 262144},
	{Key: "novita/zai-org/glm-5", Name: "GLM-5", Input: "text", ContextWindow: 202752},
	{Key: "nvidia/minimaxai/minimax-m2.5", Name: "MiniMax M2.5", Input: "text", ContextWindow: 196608},
	{Key: "nvidia/minimaxai/minimax-m2.7", Name: "Minimax M2.7", Input: "text", ContextWindow: 196608},
	{Key: "nvidia/moonshotai/kimi-k2.5", Name: "Kimi K2.5", Input: "text", ContextWindow: 262144},
	{Key: "nvidia/nemotron-3-super-120b-a12b", Name: "NVIDIA Nemotron 3 Super 120B", Input: "text", ContextWindow: 1048576},
	{Key: "nvidia/nemotron-3-ultra-550b-a55b", Name: "NVIDIA Nemotron 3 Ultra 550B", Input: "text", ContextWindow: 1000000},
	{Key: "nvidia/z-ai/glm-5.1", Name: "GLM 5.1", Input: "text", ContextWindow: 202752},
	{Key: "nvidia/z-ai/glm5", Name: "GLM-5", Input: "text", ContextWindow: 202752},
	{Key: "ollama-cloud/glm-5.1:cloud", Name: "glm-5.1:cloud", Input: "text", ContextWindow: 128000},
	{Key: "ollama-cloud/glm-5.2:cloud", Name: "glm-5.2:cloud", Input: "text", ContextWindow: 1000000},
	{Key: "ollama-cloud/kimi-k2.5:cloud", Name: "kimi-k2.5:cloud", Input: "text", ContextWindow: 128000},
	{Key: "ollama-cloud/minimax-m2.7:cloud", Name: "minimax-m2.7:cloud", Input: "text", ContextWindow: 128000},
	{Key: "together/deepseek-ai/DeepSeek-V4-Pro", Name: "DeepSeek V4 Pro", Input: "text", ContextWindow: 512000},
	{Key: "together/meta-llama/Llama-3.3-70B-Instruct-Turbo", Name: "Llama 3.3 70B Instruct Turbo", Input: "text", ContextWindow: 131072},
	{Key: "together/moonshotai/Kimi-K2.6", Name: "Kimi K2.6 FP4", Input: "text+image", ContextWindow: 262144},
	{Key: "together/Qwen/Qwen2.5-7B-Instruct-Turbo", Name: "Qwen2.5 7B Instruct Turbo", Input: "text", ContextWindow: 32768},
	{Key: "together/zai-org/GLM-5.1", Name: "GLM 5.1 FP4", Input: "text", ContextWindow: 202752},
	{Key: "volcengine/deepseek-v3-2-251201", Name: "DeepSeek V3.2", Input: "text+image", ContextWindow: 128000},
	{Key: "volcengine/doubao-seed-1-8-251228", Name: "Doubao Seed 1.8", Input: "text+image", ContextWindow: 256000},
	{Key: "volcengine/doubao-seed-code-preview-251028", Name: "doubao-seed-code-preview-251028", Input: "text+image", ContextWindow: 256000},
	{Key: "volcengine/glm-4-7-251222", Name: "GLM 4.7", Input: "text+image", ContextWindow: 200000},
	{Key: "volcengine/kimi-k2-5-260127", Name: "Kimi K2.5", Input: "text+image", ContextWindow: 256000},
	{Key: "volcengine-plan/ark-code-latest", Name: "Ark Coding Plan", Input: "text", ContextWindow: 256000},
	{Key: "volcengine-plan/doubao-seed-code", Name: "Doubao Seed Code", Input: "text", ContextWindow: 256000},
	{Key: "volcengine-plan/doubao-seed-code-preview-251028", Name: "Doubao Seed Code Preview", Input: "text", ContextWindow: 256000},
	{Key: "volcengine-plan/glm-4.7", Name: "GLM 4.7 Coding", Input: "text", ContextWindow: 200000},
	{Key: "volcengine-plan/kimi-k2-thinking", Name: "Kimi K2 Thinking", Input: "text", ContextWindow: 256000},
	{Key: "volcengine-plan/kimi-k2.5", Name: "Kimi K2.5 Coding", Input: "text", ContextWindow: 256000},
	{Key: "xiaomi/mimo-v2-flash", Name: "Xiaomi MiMo V2 Flash", Input: "text", ContextWindow: 262144},
	{Key: "xiaomi/mimo-v2-omni", Name: "Xiaomi MiMo V2 Omni", Input: "text+image", ContextWindow: 262144},
	{Key: "xiaomi/mimo-v2-pro", Name: "Xiaomi MiMo V2 Pro", Input: "text", ContextWindow: 1048576},
	{Key: "moonshot/kimi-k2.6", Name: "Kimi K2.6", Input: "text+image", ContextWindow: 262144},
	{Key: "moonshot/kimi-k2.7-code", Name: "Kimi K2.7 Code", Input: "text+image", ContextWindow: 262144},
}

func fixedModels() []Model {
	return append([]Model(nil), fixedModelCatalog...)
}

func fixedModelProvider(provider string) bool {
	for _, model := range fixedModelCatalog {
		modelProvider, _, _ := strings.Cut(model.Key, "/")
		if modelProvider == provider {
			return true
		}
	}
	return false
}

func fixedAPIKeyProvider(provider string) bool {
	switch provider {
	case "anthropic", "byteplus", "byteplus-plan", "cohere", "deepseek", "meta", "minimax", "mistral", "moonshot", "novita", "nvidia", "ollama-cloud", "together", "volcengine", "volcengine-plan", "xiaomi":
		return true
	default:
		return false
	}
}
