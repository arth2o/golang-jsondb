package testutil

import (
	"net"
	"testing"
	"time"
)

// WaitForServer attempts to connect to a server for a specified duration
func WaitForServer(t *testing.T, address string, timeout time.Duration) error {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        conn, err := net.Dial("tcp", address)
        if err == nil {
            conn.Close()
            return nil
        }
        time.Sleep(100 * time.Millisecond)
    }
    return net.ErrClosed
}

// GetFreePort returns an available port number
func GetFreePort() (int, error) {
    addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
    if err != nil {
        return 0, err
    }
    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return 0, err
    }
    defer l.Close()
    return l.Addr().(*net.TCPAddr).Port, nil
}