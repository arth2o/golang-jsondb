package server

import (
	"bufio"
	"fmt"
	"io"
	"jsondb/internal/config"
	"jsondb/internal/engine"
	"log"
	"net"
	"strings"
	"time"
)
type ClientConnection struct {
    ID         string
    LastAccess time.Time
    Connected  bool
    Conn       net.Conn
    Reader     *bufio.Reader
    Authenticated bool
}

type Server struct {
    Engine    *engine.MemoryEngine
    Password  string
    Debug     bool
    Config    *config.Config
    Listener  net.Listener
    isRunning bool
    shutdownCh chan struct{}
}

func NewServer(cfg *config.Config) (*Server, error) {
    eng, err := engine.NewMemoryEngine(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create engine: %v", err)
    }

    return &Server{
        Engine:     eng,
        Password:   cfg.Password,
        Debug:      cfg.Debug,
        Config:     cfg,
        isRunning:  false,
        shutdownCh: make(chan struct{}),
    }, nil
}

func (s *Server) Start() error {
    var err error
    s.Listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
    if err != nil {
        return fmt.Errorf("failed to start server: %v", err)
    }

    s.isRunning = true

    // Display server configuration
    log.Printf("Server Configuration:")
    log.Printf("- Port: %d", s.Config.Port)
    log.Printf("- Debug Mode: %v", s.Debug)
    log.Printf("- Encryption Enabled: %v", s.Config.EnableEncryption)
    if s.Config.EnableEncryption {
        log.Printf("- Encryption Key Length: %d", len(s.Config.EncryptionKey))
    }
    log.Printf("- Memory Dump: %v", s.Config.DumpMemoryOn)
    if s.Config.DumpMemoryOn {
        log.Printf("- Dump Path: %s", s.Config.DumpPath)
        log.Printf("- Dump Interval: %d seconds", s.Config.DumpMemoryEverySecond)
        log.Printf("- Restore From Dump: %v", s.Config.RestoreMemoryDumpAtStart)
    }
    log.Printf("- Environment: %s", s.Config.Environment)
    
    fmt.Printf("\nServer listening on port %d\n", s.Config.Port)

    go func() {
        for s.isRunning {
            conn, err := s.Listener.Accept()
            if err != nil {
                if !s.isRunning {
                    return
                }
                log.Printf("Error accepting connection: %v", err)
                continue
            }
            go s.handleConnection(conn)
        }
    }()

    return nil
}

func (s *Server) Stop() error {
    s.isRunning = false
    close(s.shutdownCh)
    if s.Listener != nil {
        return s.Listener.Close()
    }
    return nil
}

func (s *Server) IsRunning() bool {
    return s.isRunning
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    reader := bufio.NewReader(conn)
    client := &ClientConnection{
        Conn:         conn,
        Reader:       reader,
        LastAccess:   time.Now(),
        Connected:    true,
        Authenticated: false,
    }

    // Send authentication prompt
    if _, err := conn.Write([]byte("AUTH_REQUIRED\n")); err != nil {
        if s.Debug {
            log.Printf("Error sending auth prompt: %v", err)
        }
        return
    }

    for {
        command, err := reader.ReadString('\n')
        if err != nil {
            if err != io.EOF && s.Debug {
                log.Printf("Error reading command: %v", err)
            }
            return
        }

        command = strings.TrimSpace(command)
        if command == "" {
            continue
        }

        if s.Debug {
            log.Printf("Received command: %s", command)
        }

        // Handle authentication
        if !client.Authenticated {
            parts := strings.Fields(command)
            if len(parts) != 2 || strings.ToUpper(parts[0]) != "AUTH" {
                conn.Write([]byte("ERROR Authentication required\n"))
                continue
            }
            if parts[1] != s.Password {
                if s.Debug {
                    log.Printf("Authentication failed: invalid password. Got: %s, Expected: %s", parts[1], s.Password)
                }
                conn.Write([]byte("ERROR Invalid password\n"))
                continue
            }
            client.Authenticated = true
            conn.Write([]byte("OK\n"))
            continue
        }

        // Execute authenticated command
        response, err := s.executeCommand(command)
        if err != nil {
            response = fmt.Sprintf("ERROR %s\n", err.Error())
        } else if response == "" {
            response = "OK\n"
        } else {
            response = fmt.Sprintf("%s\n", response)
        }

        if _, err := conn.Write([]byte(response)); err != nil {
            if s.Debug {
                log.Printf("Error writing response: %v", err)
            }
            return
        }
    }
}

func (s *Server) executeCommand(command string) (string, error) {
    parts := strings.Fields(command)
    if len(parts) == 0 {
        return "", fmt.Errorf("empty command")
    }

    cmd := strings.ToUpper(parts[0])
    switch cmd {
    case "PING":
        return "PONG", nil
        
    case "SET":
        if len(parts) < 3 {
            return "", fmt.Errorf("SET command requires key and value")
        }
        key := parts[1]
        // Join the remaining parts to handle JSON with spaces
        value := strings.Join(parts[2:], " ")
        
        // Remove surrounding quotes if present
        value = strings.TrimPrefix(value, "\"")
        value = strings.TrimSuffix(value, "\"")
        
        if err := s.Engine.Set(key, value); err != nil {
            return "", err
        }
        return "OK", nil

    case "GET":
        if len(parts) != 2 {
            return "", fmt.Errorf("GET command requires key")
        }
        
        value, err := s.Engine.Get(parts[1])
        if err != nil {
            if err == engine.ErrKeyNotFound {
                return "nil", nil
            }
            return "", err
        }
        
        // Return the raw JSON string without any encoding
        return string(value), nil

    case "DELETE":
        if len(parts) != 2 {
            return "", fmt.Errorf("DELETE command requires key")
        }
        err := s.Engine.Delete(parts[1])
        if err != nil {
            return "", err
        }
        return "OK", nil

    case "TTL":
        if len(parts) != 2 {
            return "", fmt.Errorf("TTL command requires key")
        }
        ttl, err := s.Engine.TTL(parts[1])
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("%d", ttl), nil

    default:
        return "", fmt.Errorf("unknown command: %s", cmd)
    }
}

func (s *Server) handleResetMemory(args []string) (string, error) {
    if len(args) != 0 {
        return "", fmt.Errorf("RESET_MEMORY command takes no arguments")
    }

    if err := s.Engine.ResetMemory(); err != nil {
        return "", fmt.Errorf("failed to reset memory: %w", err)
    }

    return "OK", nil
}



