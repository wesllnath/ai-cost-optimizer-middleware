package routing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type RouteHandler struct {
	EncryptionKey []byte
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages json.RawMessage `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewRouteHandler(key []byte) *RouteHandler {
	return &RouteHandler{
		EncryptionKey: key,
	}
}

func (h *RouteHandler) HandleInference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error":"method not allowed"}`))
		return
	}

	var incomingReq struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
	}

	if err := json.NewDecoder(r.Body).Decode(&incomingReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid request payload"}`))
		return
	}

	if incomingReq.Prompt == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"prompt cannot be empty"}`))
		return
	}

	targetModel := DetermineCheaperModel(incomingReq.Model)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"upstream provider key missing"}`))
		return
	}

	rawMessages, _ := json.Marshal([]map[string]string{
		{"role": "user", "content": incomingReq.Prompt},
	})

	openAIReqPayload := OpenAIRequest{
		Model:    targetModel,
		Messages: rawMessages,
	}

	bodyBytes, err := json.Marshal(openAIReqPayload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed to construct payload"}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openai.com", bytes.NewBuffer(bodyBytes))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"upstream request construction failure"}`))
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: 35 * time.Second,
	}

	httpRes, err := client.Do(httpReq)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"failed to connect to OpenAI"}`))
		return
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"upstream provider returned an error status"}`))
		return
	}

	var openAIRes OpenAIResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&openAIRes); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed to parse upstream response"}`))
		return
	}

	var aiTextOutput string
	if len(openAIRes.Choices) > 0 {
		aiTextOutput = openAIRes.Choices[0].Message.Content
	}

	var savedUSD float64
	if incomingReq.Model != targetModel {
		savedUSD = 0.015
	}

	finalRes := AIResponse{
		Response: aiTextOutput,
		Cached:   false,
		SavedUSD: savedUSD,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(finalRes)
}
