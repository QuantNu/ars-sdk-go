package arsclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client for interacting with the Agent Registration Server
type Client struct {
	ServerURL string
	AgentID   string
	AuthToken string
	HttpClient *http.Client
}

// AgentDetails represents an agent's information
type AgentDetails struct {
	ID           string            `json:"id,omitempty"`
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Endpoint     string            `json:"endpoint"`
	Capabilities []string          `json:"capabilities"`
	PublicKey    string            `json:"publicKey"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	RegisteredAt int64             `json:"registeredAt,omitempty"`
	UpdatedAt    int64             `json:"updatedAt,omitempty"`
}

// RegisterResponse is returned when registering an agent
type RegisterResponse struct {
	Agent AgentDetails `json:"agent"`
	Token string       `json:"token"`
}

// DiscoverResponse is returned when discovering agents
type DiscoverResponse struct {
	Agents []AgentDetails `json:"agents"`
	Total  int            `json:"total"`
}

// VerifyResponse is returned when verifying an agent
type VerifyResponse struct {
	Verified bool `json:"verified"`
}

// NewClient creates a new ARS client
func NewClient(serverURL string) *Client {
	return &Client{
		ServerURL: strings.TrimSuffix(serverURL, "/"),
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterAgent registers a new agent with the ARS
func (c *Client) RegisterAgent(agent AgentDetails) (*RegisterResponse, error) {
	url := fmt.Sprintf("%s/api/agents", c.ServerURL)
	
	body, err := json.Marshal(agent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent: %w", err)
	}
	
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var response RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Save credentials for future calls
	c.AgentID = response.Agent.ID
	c.AuthToken = response.Token
	
	return &response, nil
}

// DiscoverAgents discovers agents based on capabilities and metadata
func (c *Client) DiscoverAgents(capabilities []string, metadataFilter map[string]string, limit, offset int) (*DiscoverResponse, error) {
	baseURL := fmt.Sprintf("%s/api/agents/discover", c.ServerURL)
	params := url.Values{}
	
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	
	if offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}
	
	for _, cap := range capabilities {
		params.Add("capability", cap)
	}
	
	if len(metadataFilter) > 0 {
		metadataJSON, err := json.Marshal(metadataFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata filter: %w", err)
		}
		params.Set("metadata", string(metadataJSON))
	}
	
	url := baseURL
	if len(params) > 0 {
		url = fmt.Sprintf("%s?%s", baseURL, params.Encode())
	}
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discovery failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var response DiscoverResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &response, nil
}

// GetAgent gets information about a specific agent
func (c *Client) GetAgent(agentID string) (*AgentDetails, error) {
	url := fmt.Sprintf("%s/api/agents/%s", c.ServerURL, agentID)
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get agent failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var agent AgentDetails
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &agent, nil
}

// UpdateAgent updates agent details
func (c *Client) UpdateAgent(details AgentDetails) (*AgentDetails, error) {
	if c.AgentID == "" || c.AuthToken == "" {
		return nil, fmt.Errorf("agent not registered or not authenticated")
	}
	
	url := fmt.Sprintf("%s/api/protected/agents/%s", c.ServerURL, c.AgentID)
	
	body, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent details: %w", err)
	}
	
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var agent AgentDetails
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &agent, nil
}

// DeregisterAgent deregisters an agent
func (c *Client) DeregisterAgent() error {
	if c.AgentID == "" || c.AuthToken == "" {
		return fmt.Errorf("agent not registered or not authenticated")
	}
	
	url := fmt.Sprintf("%s/api/protected/agents/%s", c.ServerURL, c.AgentID)
	
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("deregistration failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	// Clear client state
	c.AgentID = ""
	c.AuthToken = ""
	
	return nil
}

// VerifyAgent verifies another agent's identity using challenge-response
func (c *Client) VerifyAgent(agentID string, challenge, signature []byte) (bool, error) {
	url := fmt.Sprintf("%s/api/agents/%s/verify", c.ServerURL, agentID)
	
	payload := map[string]string{
		"challenge": base64.StdEncoding.EncodeToString(challenge),
		"signature": base64.StdEncoding.EncodeToString(signature),
	}
	
	body, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("verification failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var response VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return response.Verified, nil
}
