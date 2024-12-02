package main

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
    // Set test environment
    os.Setenv("GO_ENV", "test")
    os.Setenv("PORT", "6380")  // Use test port
    
    // Create temporary .env.test if it doesn't exist
    envFile := filepath.Join(".", ".env.test")
    if _, err := os.Stat(envFile); os.IsNotExist(err) {
        content := []byte("PORT=6380\nENVIRONMENT=testing\nENABLE_ENCRYPTION=false\nMAX_CONNECTIONS=2\nDEBUG=true")
        if err := os.WriteFile(envFile, content, 0644); err != nil {
            t.Fatalf("Failed to create test env file: %v", err)
        }
        defer os.Remove(envFile)
    }
    
    done := make(chan bool)
    go func() {
        main()
        done <- true
    }()

    // Wait for server to start
    time.Sleep(100 * time.Millisecond)

    t.Run("Server Starts and Accepts Connections", func(t *testing.T) {
        conn, err := net.Dial("tcp", "localhost:6380")
        if err != nil {
            t.Fatalf("Failed to connect to server: %v", err)
        }
        conn.Close()
    })

    // Cleanup
    if conn, err := net.Dial("tcp", "localhost:6380"); err == nil {
        conn.Close()
    }
}
