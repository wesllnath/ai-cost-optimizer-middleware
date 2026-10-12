package routing

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type AIRequest struct {
	Prompt      string `json:"prompt"`
	Model       string `json:"model"`
	Temperature float64 `json:"temperature"`
}

type AIResponse struct {
	Response string `json:"response"`
	Cached   bool   `json:"cached"`
	SavedUSD float64 `json:"saved_usd"`
}

func ComputePromptHash(prompt string) string {
	cleaned := strings.TrimSpace(strings.ToLower(prompt))
	hash := sha256.New()
	hash.Write([]byte(cleaned))
	return hex.EncodeToString(hash.Sum(nil))
}

func DetermineCheaperModel(model string) string {
	switch model {
	case "gpt-4o", "gpt-4-turbo", "claude-3-5-sonnet":
		return "gpt-4o-mini"
	default:
		return model
	}
}
