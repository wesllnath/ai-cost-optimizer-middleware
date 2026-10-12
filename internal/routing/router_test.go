package routing

import (
	"testing"
)

func TestComputePromptHashConsistency(t *testing.T) {
	hash1 := ComputePromptHash("  Analyze This Report  ")
	hash2 := ComputePromptHash("analyze this report")

	if hash1 != hash2 {
		t.Errorf("expected hashes to be identical, got %s and %s", hash1, hash2)
	}
}

func TestDetermineCheaperModelRouting(t *testing.T) {
	cheapModel := DetermineCheaperModel("gpt-4o")
	if cheapModel != "gpt-4o-mini" {
		t.Errorf("expected gpt-4o-mini, got %s", cheapModel)
	}

	standardModel := DetermineCheaperModel("custom-llm-v1")
	if standardModel != "custom-llm-v1" {
		t.Errorf("expected custom-llm-v1 to remain unchanged, got %s", standardModel)
	}
}
