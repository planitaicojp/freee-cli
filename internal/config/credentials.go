package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const credentialsFile = "credentials.yaml"

// CredentialsStore holds OAuth2 credentials per profile.
type CredentialsStore struct {
	Profiles map[string]OAuthCredentials `yaml:"profiles"`
}

// OAuthCredentials stores OAuth2 tokens for a profile.
type OAuthCredentials struct {
	AccessToken  string    `yaml:"access_token"`
	RefreshToken string    `yaml:"refresh_token"`
	ExpiresAt    time.Time `yaml:"expires_at"`
	ClientID     string    `yaml:"client_id"`
	ClientSecret string    `yaml:"client_secret"`
}

func LoadCredentials() (*CredentialsStore, error) {
	path := filepath.Join(DefaultConfigDir(), credentialsFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &CredentialsStore{Profiles: map[string]OAuthCredentials{}}, nil
		}
		return nil, fmt.Errorf("reading credentials: %w", err)
	}

	var store CredentialsStore
	if err := yaml.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parsing credentials: %w", err)
	}
	if store.Profiles == nil {
		store.Profiles = map[string]OAuthCredentials{}
	}
	return &store, nil
}

func (s *CredentialsStore) Save() error {
	dir := DefaultConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling credentials: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, credentialsFile), data, 0600)
}

func (s *CredentialsStore) Get(profile string) (OAuthCredentials, bool) {
	c, ok := s.Profiles[profile]
	return c, ok
}

func (s *CredentialsStore) Set(profile string, cred OAuthCredentials) {
	s.Profiles[profile] = cred
}

func (s *CredentialsStore) Delete(profile string) {
	delete(s.Profiles, profile)
}

// IsTokenValid returns true if the access token has more than 5 minutes remaining.
func (s *CredentialsStore) IsTokenValid(profile string) bool {
	c, ok := s.Profiles[profile]
	if !ok {
		return false
	}
	return time.Until(c.ExpiresAt) > 5*time.Minute
}
