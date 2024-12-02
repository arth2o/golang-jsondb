package config

import (
	"testing"
)

func TestNewDevelopmentConfig(t *testing.T) {
    tests := []struct {
        name string
        want *Config
    }{
        {
            name: "Default Development Config",
            want: &Config{
                Port:             5555,
                Environment:      Development,
                EnableEncryption: false,
                MaxConnections:   -1,
                Debug:           false,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NewDevelopmentConfig()
            if got.Port != tt.want.Port {
                t.Errorf("Port = %v, want %v", got.Port, tt.want.Port)
            }
            if got.Environment != tt.want.Environment {
                t.Errorf("Environment = %v, want %v", got.Environment, tt.want.Environment)
            }
            if got.EnableEncryption != tt.want.EnableEncryption {
                t.Errorf("EnableEncryption = %v, want %v", got.EnableEncryption, tt.want.EnableEncryption)
            }
            if got.MaxConnections != tt.want.MaxConnections {
                t.Errorf("MaxConnections = %v, want %v", got.MaxConnections, tt.want.MaxConnections)
            }
        })
    }
}

func TestConfigChaining(t *testing.T) {
    cfg := NewDevelopmentConfig().
        WithPassword("testpass").
        WithEncryptionKey("12345678901234567890123456789012")

    if cfg.Password != "testpass" {
        t.Errorf("Password = %v, want %v", cfg.Password, "testpass")
    }

    if cfg.EncryptionKey != "12345678901234567890123456789012" {
        t.Errorf("EncryptionKey not set correctly")
    }
}

func TestNewTestConfig(t *testing.T) {
    cfg := NewTestConfig()

    tests := []struct {
        name     string
        got      interface{}
        expected interface{}
    }{
        {"Port", cfg.Port, 6380},
        {"Password", cfg.Password, "test-password"},
        {"Environment", cfg.Environment, Testing},
        {"MaxConnections", cfg.MaxConnections, 2},
        {"Debug", cfg.Debug, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.got != tt.expected {
                t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
            }
        })
    }
}
