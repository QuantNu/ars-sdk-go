package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"

	arsclient "github.com/yourusername/ars/sdk/go"
)

func main() {
	// Create a new client
	client := arsclient.NewClient("https://ars.example.com")

	// Generate a key pair for testing
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
		Name:         "Example Go Agent",
		Version:      "1.0.0",
		Endpoint:     "https://agent.example.com/api",
		Capabilities: []string{"query", "response", "calculation"},
		PublicKey:    base64.StdEncoding.EncodeToString(publicKeyPEM),
		Metadata: map[string]string{
			"description": "Example agent for demonstration",
			"creator":     "ARS SDK",
			"language":    "Go",
		},
	})
	if err != nil {
		log.Fatalf("Failed to register agent: %v", err)
	}

	fmt.Printf("Registered agent with ID: %s\n", response.Agent.ID)
	fmt.Printf("Authentication token: %s\n", response.Token)

	// Discover agents with the "query" capability
	discoveryResponse, err := client.DiscoverAgents(
		[]string{"query"},
		nil,
		5,
		0,
	)
	if err != nil {
		log.Fatalf("Failed to discover agents: %v", err)
	}

	fmt.Printf("Found %d agents with 'query' capability\n", discoveryResponse.Total)
	for i, agent := range discoveryResponse.Agents {
		fmt.Printf("  %d. %s (%s)\n", i+1, agent.Name, agent.ID)
	}

	// Update the agent
	updatedAgent, err := client.UpdateAgent(arsclient.AgentDetails{
		Version: "1.0.1",
		Metadata: map[string]string{
			"status": "active",
		},
	})
	if err != nil {
		log.Fatalf("Failed to update agent: %v", err)
	}

	fmt.Printf("Updated agent to version %s\n", updatedAgent.Version)

	// Deregister the agent
	if err := client.DeregisterAgent(); err != nil {
		log.Fatalf("Failed to deregister agent: %v", err)
	}

	fmt.Println("Agent deregistered successfully")
}
