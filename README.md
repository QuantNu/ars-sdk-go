# ARS Go SDK

Go SDK for the Agent Registration Server (ARS).

## Overview

The ARS Go SDK provides a client library for interacting with the Agent Registration Server, enabling Go applications to register agents, discover capabilities, manage sessions, and execute operations across different agent protocols.

## Installation

### For Users

```bash
go get github.com/quantnu/ars-go-sdk
```

### For Developers

Clone the repository and install dependencies:

```bash
git clone https://github.com/quantnu/ars.git
cd ars/sdk/go
go mod tidy
```

## Building Independently

The Go SDK can be built independently using the included Makefile:

```bash
# Download dependencies, run tests, and build
make

# Individual steps
make deps      # Download dependencies
make build     # Build the project
make test      # Run tests
make clean     # Clean build artifacts
make lint      # Run linters
make fmt       # Format code
```

If you don't have `make` available, you can use these commands directly:

```bash
# Download dependencies
go mod tidy

# Build the project
go build -o bin/ars-client ./...

# Run tests
go test -v ./...

# Format code
go fmt ./...
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    ars "github.com/quantnu/ars-go-sdk"
)

func main() {
    // Create client
    client, err := ars.NewClient("https://ars.example.com")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Register an agent
    agent, token, err := client.RegisterAgent(context.Background(), ars.AgentRegistration{
        Name:           "Example Agent",
        Description:    "An example agent demonstrating basic functionality",
        Capabilities:   []string{"translate", "summarize"},
        Endpoint:       "https://example.com/agent",
        Protocol:       ars.ProtocolMCP,
        ProtocolVersion: "1.0",
        PublicKey:      "example-public-key",
    })
    if err != nil {
        log.Fatalf("Failed to register agent: %v", err)
    }

    fmt.Printf("Registered agent with ID: %s\n", agent.ID)
    fmt.Printf("Agent token: %s\n", token)

    // Discover agents with specific capabilities
    agents, err := client.DiscoverAgents(context.Background(), ars.DiscoveryRequest{
        Capabilities: []string{"translate"},
        TrustLevels:  []ars.TrustLevel{ars.TrustLevelVerified, ars.TrustLevelPartner},
    })
    if err != nil {
        log.Fatalf("Failed to discover agents: %v", err)
    }

    fmt.Printf("Found %d agents\n", len(agents))

    // Execute a task
    result, err := client.ExecuteTask(context.Background(), "translate", map[string]interface{}{
        "text":           "Hello world",
        "sourceLanguage": "en",
        "targetLanguage": "fr",
    })
    if err != nil {
        log.Fatalf("Failed to execute task: %v", err)
    }

    fmt.Printf("Translation result: %v\n", result)
}
```

## Session Management

```go
package main

import (
    "context"
    "fmt"
    "log"

    ars "github.com/quantnu/ars-go-sdk"
)

func main() {
    // Create client
    client, err := ars.NewSessionClient("https://ars.example.com")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create a session
    session, sessionToken, err := client.CreateSession(context.Background(), ars.SessionData{
        Context: map[string]interface{}{
            "user": "user-123",
            "preferences": map[string]interface{}{
                "language": "en",
            },
        },
    })
    if err != nil {
        log.Fatalf("Failed to create session: %v", err)
    }

    fmt.Printf("Created session with ID: %s\n", session.ID)

    // Execute task using session context
    result, err := client.ExecuteTaskWithSession(context.Background(), "translate", map[string]interface{}{
        "text":           "Hello world",
        "targetLanguage": "fr",
    }, session.ID)
    if err != nil {
        log.Fatalf("Failed to execute task: %v", err)
    }

    fmt.Printf("Translation result: %v\n", result)

    // Update session with new information
    err = client.UpdateSession(context.Background(), session.ID, ars.SessionData{
        Context: map[string]interface{}{
            "history": []map[string]interface{}{
                {
                    "task":   "translate",
                    "input":  map[string]string{"text": "Hello world", "targetLanguage": "fr"},
                    "output": "Bonjour le monde",
                },
            },
        },
    })
    if err != nil {
        log.Fatalf("Failed to update session: %v", err)
    }
}
```

## Features

- **Agent Registration**: Register agents with the ARS
- **Agent Discovery**: Find agents based on capabilities and trust levels
- **Trust Verification**: Verify agent identity and trust levels
- **Session Management**: Maintain stateful interactions between agents
- **Cross-Protocol Operation**: Work with agents across different protocols
- **Error Handling**: Comprehensive error handling and reporting
- **Context Support**: Full Go context.Context support for cancellation and timeouts

## Contributing

Contributions are welcome! Please see the main repository's CONTRIBUTING.md for guidelines.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
