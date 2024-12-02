package server

import (
	"bufio"
	"fmt"
	"jsondb/internal/config"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	cfg := &config.Config{
		Port:     6380,
		Password: "testpass",
		Debug:    true,
	}

	srv, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	// Small wait for server to start
	time.Sleep(10 * time.Millisecond)

	// Test connection with timeout
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", cfg.Port), time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Set shorter timeouts
	conn.SetReadDeadline(time.Now().Add(time.Second))
	conn.SetWriteDeadline(time.Now().Add(time.Second))

	reader := bufio.NewReader(conn)

	// Read AUTH_REQUIRED prompt
	prompt, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read auth prompt: %v", err)
	}
	if strings.TrimSpace(prompt) != "AUTH_REQUIRED" {
		t.Fatalf("Expected AUTH_REQUIRED prompt, got: %s", prompt)
	}

	// Authenticate
	fmt.Fprintf(conn, "AUTH %s\n", cfg.Password)
	authResponse, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read auth response: %v", err)
	}
	if strings.TrimSpace(authResponse) != "OK" {
		t.Fatalf("Authentication failed: %s", authResponse)
	}

	// Test PING
	fmt.Fprintf(conn, "PING\n")
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read PING response: %v", err)
	}
	if strings.TrimSpace(response) != "PONG" {
		t.Fatalf("Expected PONG, got: %s", response)
	}
}

func TestServerCommands(t *testing.T) {
	cfg := &config.Config{
		Port:     6381,
		Password: "testpass",
		Debug:    true,
	}

	srv, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	// Small wait for server to start
	time.Sleep(10 * time.Millisecond)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", cfg.Port), time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Set shorter timeouts
	conn.SetReadDeadline(time.Now().Add(time.Second))
	conn.SetWriteDeadline(time.Now().Add(time.Second))

	reader := bufio.NewReader(conn)

	// Handle initial AUTH_REQUIRED
	prompt, _ := reader.ReadString('\n')
	if strings.TrimSpace(prompt) != "AUTH_REQUIRED" {
		t.Fatalf("Expected AUTH_REQUIRED prompt, got: %s", prompt)
	}

	// Authenticate
	fmt.Fprintf(conn, "AUTH %s\n", cfg.Password)
	authResponse, _ := reader.ReadString('\n')
	if strings.TrimSpace(authResponse) != "OK" {
		t.Fatalf("Authentication failed: %s", authResponse)
	}

	// Test basic commands
	commands := []struct {
		cmd      string
		expected string
	}{
		{"SET test:key value123", "OK"},
		{"GET test:key", "value123"},
		{"DEL test:key", "OK"},
		{"GET test:key", "nil"},
	}

	for _, cmd := range commands {
		// Reset deadlines for each command
		conn.SetReadDeadline(time.Now().Add(time.Second))
		conn.SetWriteDeadline(time.Now().Add(time.Second))

		fmt.Fprintf(conn, "%s\n", cmd.cmd)
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Command '%s' failed: %v", cmd.cmd, err)
		}
		if strings.TrimSpace(response) != cmd.expected {
			t.Errorf("Command '%s': got %q, want %q", cmd.cmd, strings.TrimSpace(response), cmd.expected)
		}
	}
}
