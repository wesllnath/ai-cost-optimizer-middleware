package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type FileRepository struct {
	mu       sync.RWMutex
	filePath string
	Logs     map[string]InferenceLog  `json:"logs"`
	Orgs     map[string]Organization  `json:"orgs"`
	Creds    map[string]APICredential `json:"creds"`
}

func NewFileRepository(path string) (*FileRepository, error) {
	repo := &FileRepository{
		filePath: path,
		Logs:     make(map[string]InferenceLog),
		Orgs:     make(map[string]Organization),
		Creds:    make(map[string]APICredential),
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return repo, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return repo, nil
	}

	_ = json.Unmarshal(data, repo)
	return repo, nil
}

func (r *FileRepository) saveToDisk() error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0600)
}

func (r *FileRepository) SaveOrganization(org *Organization) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Orgs[org.ID] = *org
	return r.saveToDisk()
}

func (r *FileRepository) GetOrganization(id string) (*Organization, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	org, exists := r.Orgs[id]
	if !exists {
		return nil, errors.New("organization not found")
	}
	return &org, nil
}

func (r *FileRepository) SaveCredential(cred *APICredential) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Creds[cred.ID] = *cred
	return r.saveToDisk()
}

func (r *FileRepository) GetCredential(orgID string, provider string) (*APICredential, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, cred := range r.Creds {
		if cred.OrganizationID == orgID && cred.ProviderName == provider {
			return &cred, nil
		}
	}
	return nil, errors.New("credential not found")
}

func (r *FileRepository) SaveInferenceLog(log *InferenceLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Logs[log.ID] = *log
	return r.saveToDisk()
}

func (r *FileRepository) GetMetrics(orgID string) ([]InferenceLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []InferenceLog
	for _, log := range r.Logs {
		if log.OrganizationID == orgID {
			results = append(results, log)
		}
	}
	return results, nil
}
