package storage

import (
	"time"
)

type Organization struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	StripeCustomerID string    `json:"stripe_customer_id"`
	PlanStatus       string    `json:"plan_status"`
	CreatedAt        time.Time `json:"created_at"`
}

type APICredential struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ProviderName   string    `json:"provider_name"`
	EncryptedKey   []byte    `json:"encrypted_key"`
	Nonce          []byte    `json:"nonce"`
	CreatedAt      time.Time `json:"created_at"`
}

type InferenceLog struct {
	ID                   string    `json:"id"`
	OrganizationID       string    `json:"organization_id"`
	PromptTokens         int       `json:"prompt_tokens"`
	CompletionTokens     int       `json:"completion_tokens"`
	OriginalEstimatedCost float64   `json:"original_estimated_cost"`
	ActualCost           float64   `json:"actual_cost"`
	SavedCost            float64   `json:"saved_cost"`
	RoutingStrategy      string    `json:"routing_strategy"`
	CreatedAt            time.Time `json:"created_at"`
}

type StorageRepository interface {
	SaveOrganization(org *Organization) error
	GetOrganization(id string) (*Organization, error)
	SaveCredential(cred *APICredential) error
	GetCredential(orgID string, provider string) (*APICredential, error)
	SaveInferenceLog(log *InferenceLog) error
	GetMetrics(orgID string) ([]InferenceLog, error)
}
