package engine

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"jsondb/internal/config"
)

func TestMemoryEngine_SetAndGet(t *testing.T) {
	cfg := &config.Config{
		EnableEncryption: false,
		Debug:           true,
	}
	
	engine, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	tests := []struct {
		name  string
		key   string
		value interface{}
		want  interface{}
	}{
		{
			name:  "String Value",
			key:   "test1",
			value: "hello",
			want:  "hello",
		},
		{
			name:  "Integer Value",
			key:   "test2",
			value: 42,
			want:  float64(42), // JSON numbers are float64
		},
		{
			name: "Complex Value",
			key:  "test3",
			value: map[string]interface{}{
				"name": "test",
				"age":  30,
			},
			want: map[string]interface{}{
				"name": "test",
				"age":  float64(30),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := engine.Set(tt.key, tt.value); err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}

			got, err := engine.Get(tt.key)
			if err != nil {
				t.Fatalf("Failed to get value: %v", err)
			}

			var decoded interface{}
			if err := json.Unmarshal(got, &decoded); err != nil {
				// If unmarshal fails, try comparing as strings
				if string(got) != tt.want.(string) {
					t.Errorf("Value mismatch for key %s: got %q, want %q",
						tt.key, string(got), tt.want)
				}
				return
			}

			if !reflect.DeepEqual(decoded, tt.want) {
				t.Errorf("Value mismatch for key %s: got %v, want %v",
					tt.key, decoded, tt.want)
			}
		})
	}
}

func TestMemoryEngine_GetByPattern(t *testing.T) {
	cfg := &config.Config{
		EnableEncryption: false,
		Debug:           false,
		DumpMemoryOn:    false,
		DumpPath:        "dump",
	}
	
	engine, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Setup test data
	testData := map[string]interface{}{
		"user:1": map[string]string{"name": "John"},
		"user:2": map[string]string{"name": "Jane"},
		"post:1": map[string]string{"title": "Hello"},
	}

	for k, v := range testData {
		if err := engine.Set(k, v); err != nil {
			t.Fatalf("Failed to set test data: %v", err)
		}
	}

	tests := []struct {
		name          string
		pattern       string
		expectedKeys  int
		shouldContain string
	}{
		{
			name:          "User Pattern",
			pattern:       "user:*",
			expectedKeys:  2,
			shouldContain: "user:",
		},
		{
			name:          "Single Pattern",
			pattern:       "post:1",
			expectedKeys:  1,
			shouldContain: "post:1",
		},
		{
			name:          "No Match Pattern",
			pattern:       "nothing:*",
			expectedKeys:  0,
			shouldContain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := engine.GetByPattern(tt.pattern)
			if err != nil {
				t.Errorf("GetByPattern() error = %v", err)
				return
			}

			if len(results) != tt.expectedKeys {
				t.Errorf("GetByPattern() got %v results, want %v", len(results), tt.expectedKeys)
			}

			if tt.shouldContain != "" && len(results) > 0 {
				found := false
				for _, result := range results {
					if strings.Contains(result.Key, tt.shouldContain) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetByPattern() results did not contain expected pattern %s", tt.shouldContain)
				}
			}
		})
	}
}

func TestMemoryEngine_TTL(t *testing.T) {
	cfg := &config.Config{
		EnableEncryption: false,
		Debug:           false,
		DumpMemoryOn:    false,
		DumpPath:        "dump",
	}
	
	engine, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	
	tests := []struct {
		name        string
		key         string
		value       string
		ttl         time.Duration
		expectedTTL time.Duration
	}{
		{
			name:        "No Expiry",
			key:         "test:no-ttl",
			value:       "value",
			ttl:         0,
			expectedTTL: -1,
		},
		{
			name:        "With Expiry",
			key:         "test:with-ttl",
			value:       "value",
			ttl:         60 * time.Second,
			expectedTTL: 60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ttl > 0 {
				err := engine.SetWithTTL(tt.key, []byte(tt.value), tt.ttl)
				if err != nil {
					t.Fatalf("SetWithTTL failed: %v", err)
				}
			} else {
				err := engine.Set(tt.key, tt.value)
				if err != nil {
					t.Fatalf("Set failed: %v", err)
				}
			}

			got, err := engine.TTL(tt.key)
			if err != nil {
				t.Fatalf("TTL failed: %v", err)
			}

			if tt.expectedTTL >= 0 && (got < tt.expectedTTL-time.Second || got > tt.expectedTTL+time.Second) {
				t.Errorf("TTL = %v, want approximately %v", got, tt.expectedTTL)
			}
		})
	}
}

func TestMemoryEngineDumpAndRestore(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "jsondb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		EnableEncryption: true,
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
		Debug:           true,
		DumpPath:        tmpDir,
		DumpMemoryOn:    true,
		DumpMemoryEverySecond: 1,
	}
	
	// Create first engine instance and set some data
	engine1, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create first engine: %v", err)
	}

	testData := map[string]interface{}{
		"key1": "string value",
		"key2": float64(42), // Use float64 for numbers
		"key3": map[string]interface{}{"nested": "value"},
		"key4": []interface{}{float64(1), float64(2), float64(3)}, // Use float64 for array numbers
	}

	// Set test data
	for k, v := range testData {
		if err := engine1.Set(k, v); err != nil {
			t.Fatalf("Failed to set key %s: %v", k, err)
		}
	}

	// Dump to disk
	if err := engine1.DumpToDisk(); err != nil {
		t.Fatalf("Failed to dump: %v", err)
	}

	// Create second engine instance and restore
	engine2, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create second engine: %v", err)
	}

	if err := engine2.RestoreFromDisk(); err != nil {
		t.Fatalf("Failed to restore: %v", err)
	}

	// Verify all data
	for k, want := range testData {
		got, err := engine2.Get(k)
		if err != nil {
			t.Errorf("Failed to get key %s after restore: %v", k, err)
			continue
		}

		var decoded interface{}
		if err := json.Unmarshal(got, &decoded); err != nil {
			t.Errorf("Failed to unmarshal value for key %s: %v", k, err)
			continue
		}

		if !reflect.DeepEqual(decoded, want) {
			t.Errorf("Value mismatch for key %s after restore:\ngot:  %#v\nwant: %#v", 
				k, decoded, want)
		}
	}
}

func TestMemoryEngine_ResetMemory(t *testing.T) {
	cfg := &config.Config{
		EnableEncryption: false,
		Debug:           false,
		DumpMemoryOn:    false,
		DumpPath:        "dump",
	}
	
	engine, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Set test data
	testData := map[string]interface{}{
		"test:key1": "value1",
		"test:key2": "value2",
	}

	// Insert test data
	for k, v := range testData {
		if err := engine.Set(k, v); err != nil {
			t.Fatalf("Failed to set test data: %v", err)
		}
	}

	// Reset memory
	if err := engine.ResetMemory(); err != nil {
		t.Fatalf("Failed to reset memory: %v", err)
	}

	// Verify all keys are gone
	matches, err := engine.GetByPattern("*")
	if err != nil {
		t.Fatalf("Failed to get keys after reset: %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("Expected 0 keys after reset, got %d", len(matches))
	}
}

func TestEncryptionRoundTrip(t *testing.T) {
	cfg := &config.Config{
		EnableEncryption: true,
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
		Debug:           true,
	}
	
	engine, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	testData := []struct {
		key      string
		value    interface{}
		expected interface{}
	}{
		{
			key:      "key1",
			value:    "apple",
			expected: "apple",
		},
		{
			key:      "key2",
			value:    "banana",
			expected: "banana",
		},
		{
			key:      "key3",
			value:    "This is a longer string to test",
			expected: "This is a longer string to test",
		},
		{
			key:      "key4",
			value:    map[string]interface{}{"test": "value"},
			expected: map[string]interface{}{"test": "value"},
		},
		{
			key:      "key5",
			value:    42,
			expected: float64(42), // JSON numbers are float64
		},
	}

	for _, tt := range testData {
		t.Run(tt.key, func(t *testing.T) {
			// Set value
			if err := engine.Set(tt.key, tt.value); err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}

			// Get value back
			got, err := engine.Get(tt.key)
			if err != nil {
				t.Fatalf("Failed to get value: %v", err)
			}

			// Try to unmarshal the result
			var decoded interface{}
			if err := json.Unmarshal(got, &decoded); err != nil {
				// If unmarshal fails, compare as strings
				if str, ok := tt.value.(string); ok {
					if string(got) != str {
						t.Errorf("Value mismatch for key %s: got %q, want %q", 
							tt.key, string(got), str)
					}
					return
				}
				t.Fatalf("Failed to unmarshal value for key %s: %v", tt.key, err)
			}

			// Compare the decoded value with expected
			if !reflect.DeepEqual(decoded, tt.expected) {
				t.Errorf("Value mismatch for key %s:\ngot:  %#v\nwant: %#v", 
					tt.key, decoded, tt.expected)
			}
		})
	}
}

func TestEncryptionWithDifferentTypes(t *testing.T) {
	cfg := &config.Config{
		EnableEncryption: true,
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
		Debug:           true,
	}
	
	engine, err := NewMemoryEngine(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	testCases := []struct {
		name     string
		key      string
		value    interface{}
		expected interface{}
	}{
		{
			name:     "String Value",
			key:      "test:string",
			value:    "Simple string test",
			expected: "Simple string test",
		},
		{
			name:     "Integer Value",
			key:      "test:int",
			value:    42,
			expected: float64(42), // JSON numbers are float64
		},
		{
			name:     "Float Value",
			key:      "test:float",
			value:    3.14,
			expected: 3.14,
		},
		{
			name:     "Boolean Value",
			key:      "test:bool",
			value:    true,
			expected: true,
		},
		{
			name:     "Null Value",
			key:      "test:null",
			value:    nil,
			expected: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Set the value
			if err := engine.Set(tt.key, tt.value); err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}

			// Get the value back
			got, err := engine.Get(tt.key)
			if err != nil {
				t.Fatalf("Failed to get value: %v", err)
			}

			// Unmarshal the result
			var decoded interface{}
			if err := json.Unmarshal(got, &decoded); err != nil {
				// For string values, compare directly
				if str, ok := tt.value.(string); ok {
					if string(got) != str {
						t.Errorf("String value mismatch:\ngot:  %s\nwant: %s", 
							string(got), str)
					}
					return
				}
				// For null values, check if the error is due to null
				if tt.value == nil && len(got) == 0 {
					return // Empty response is fine for null
				}
				t.Fatalf("Failed to unmarshal value: %v", err)
			}

			// Compare the values
			if !reflect.DeepEqual(decoded, tt.expected) {
				t.Errorf("Value mismatch:\ngot:  %#v\nwant: %#v", decoded, tt.expected)
			}
		})
	}
}