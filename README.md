# Agent Registration Server Go Client

A Go client library for interacting with the Agent Registration Server (ARS).

## Installation

```bash
go get github.com/yourusername/ars/sdk/go
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	
	arsclient "github.com/yourusername/ars/sdk/go"
)

func main() {
	// Create a new client
	client := arsclient.NewClient("https://ars.example.com")
	
	// Generate a key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	
	// Convert public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}
	
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	
	// Register the agent
	response, err := client.RegisterAgent(arsclient.AgentDetails{
		Name:         "My Agent",
		Version:      "1.0.0",
		Endpoint:     "https://myagent.example.com/api",
		Capabilities: []string{"query", "response"},
		PublicKey:    base64.StdEncoding.EncodeToString(publicKeyPEM),
		Metadata: map[string]string{
			"description": "My awesome agent",
		},
	})
	if err != nil {
		log.Fatalf("Failed to register agent: %v", err)
	}
	
	fmt.Printf("Registered agent with ID: %s\n", response.Agent.ID)
	
	// Discover other agents
	agents, err := client.DiscoverAgents(
		[]string{"query"},
		nil,
		5,
		0,
	)
	if err != nil {
		log.Fatalf("Failed to discover agents: %v", err)
	}
	
	fmt.Printf("Found %d agents with 'query' capability\n", agents.Total)
	
	// Deregister when done
	if err := client.DeregisterAgent(); err != nil {
		log.Fatalf("Failed to deregister agent: %v", err)
	}
}
```

## API Reference

### Client

```go
client := arsclient.NewClient(serverURL)
```

Create a new client with the server URL.

### Methods

- `RegisterAgent(agent AgentDetails) (*RegisterResponse, error)`: Register a new agent
- `DiscoverAgents(capabilities []string, metadataFilter map[string]string, limit, offset int) (*DiscoverResponse, error)`: Find agents matching criteria
- `GetAgent(agentID string) (*AgentDetails, error)`: Get details for a specific agent
- `UpdateAgent(details AgentDetails) (*AgentDetails, error)`: Update your agent
- `DeregisterAgent() error`: Remove your agent from the registry
- `VerifyAgent(agentID string, challenge, signature []byte) (bool, error)`: Verify another agent's identity

## Error Handling

All methods return errors if the server returns an error status. You should handle these errors in your code.

## License

MIT
