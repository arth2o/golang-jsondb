package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Environment type for server environment
type Environment string

const (
    Development Environment = "development"
    Production  Environment = "production"
    Testing     Environment = "testing"
)

// Config holds all server configuration
type Config struct {
    Port                    int
    Password                string
    EncryptionKey           string
    Environment             Environment
    EnableEncryption        bool
    MaxConnections          int
    Debug                   bool
    DumpMemoryOn           bool
    DumpMemoryEverySecond  int
    RestoreMemoryDumpAtStart bool
    DumpPath               string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
    // Determine environment
    env := getEnvStr("ENVIRONMENT", "development")
    
    // Load appropriate .env file
    envFile := ".env"
    if env == "development" {
        envFile = ".env.development"
    } else if env == "production" {
        envFile = ".env.production"
    }

    // Load environment file if it exists
    if err := godotenv.Load(envFile); err != nil {
        // Only log as warning, don't return error as env vars might be set directly
        log.Printf("Warning: Error loading %s file: %v", envFile, err)
    } else {
        log.Printf("Successfully loaded environment file: %s", envFile)
    }

    // Create config based on environment
    var cfg *Config
    switch env {
    case "development":
        cfg = NewDevelopmentConfig()
    case "production":
        cfg = NewProductionConfig()
    default:
        cfg = NewDevelopmentConfig()
    }

    // Validate config
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("invalid configuration: %v", err)
    }

    return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
    if c.Port <= 0 {
        return fmt.Errorf("invalid port number: %d", c.Port)
    }
    if c.Password == "" {
        return fmt.Errorf("server password cannot be empty")
    }
    if c.EnableEncryption && c.EncryptionKey == "" {
        return fmt.Errorf("encryption enabled but no key provided")
    }
    if c.DumpMemoryOn && c.DumpPath == "" {
        return fmt.Errorf("memory dump enabled but no dump path provided")
    }
    return nil
}

// NewDevelopmentConfig creates a new configuration for development environment
func NewDevelopmentConfig() *Config {
    return &Config{
        Port:                    getEnvInt("PORT", 5555),
        Password:               getEnvStr("SERVER_PASSWORD", "password"),
        EncryptionKey:          getEnvStr("ENCRYPTION_KEY", ""),
        EnableEncryption:       getEnvBool("ENABLE_ENCRYPTION", false),
        Environment:            Development,
        MaxConnections:         getEnvInt("MAX_CONNECTIONS", -1),
        Debug:                  getEnvBool("DEBUG", true),
        DumpMemoryOn:          getEnvBool("DUMP_MEMORY_ON", false),
        DumpMemoryEverySecond: getEnvInt("DUMP_MEMORY_EVERY_SECOND", 60),
        RestoreMemoryDumpAtStart: getEnvBool("RESTORE_MEMORY_DUMP_AT_START", false),
        DumpPath:              getEnvStr("DUMP_PATH", "data/dump"),
    }
}

// NewProductionConfig creates a new configuration for production environment
func NewProductionConfig() *Config {
    return &Config{
        Port:                    getEnvInt("PORT", 5555),
        Password:               getEnvStr("SERVER_PASSWORD", ""),
        EncryptionKey:          getEnvStr("ENCRYPTION_KEY", ""),
        EnableEncryption:       getEnvBool("ENABLE_ENCRYPTION", true),
        Environment:            Production,
        MaxConnections:         getEnvInt("MAX_CONNECTIONS", 1000),
        Debug:                  getEnvBool("DEBUG", false),
        DumpMemoryOn:          getEnvBool("DUMP_MEMORY_ON", true),
        DumpMemoryEverySecond: getEnvInt("DUMP_MEMORY_EVERY_SECOND", 300),
        RestoreMemoryDumpAtStart: getEnvBool("RESTORE_MEMORY_DUMP_AT_START", true),
        DumpPath:              getEnvStr("DUMP_PATH", "data/dump"),
    }
}

func getEnvStr(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func getEnvBool(key string, fallback bool) bool {
    if value := os.Getenv(key); value != "" {
        return strings.ToLower(value) == "true"
    }
    return fallback
}

func getEnvInt(key string, fallback int) int {
    if value := os.Getenv(key); value != "" {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return fallback
}

// NewTestConfig creates a configuration suitable for testing
func NewTestConfig() *Config {
    return &Config{
        Port:                    6380,
        Password:                "test-password",
        EncryptionKey:           "test-key",
        Environment:             Testing,
        EnableEncryption:        false,
        MaxConnections:          2,
        Debug:                   false,
        DumpMemoryOn:           false,
        DumpMemoryEverySecond:  5,
        RestoreMemoryDumpAtStart: false,
        DumpPath:               "dump",
    }
}

// WithPassword sets the password and returns the config
func (c *Config) WithPassword(password string) *Config {
    c.Password = password
    return c
}

// WithEncryptionKey sets the encryption key and returns the config
func (c *Config) WithEncryptionKey(key string) *Config {
    c.EncryptionKey = key
    return c
}
