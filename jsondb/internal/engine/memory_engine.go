package engine

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"jsondb/internal/config"
	"jsondb/internal/encryption"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Match struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type KeyData struct {
	Value     []byte    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

type MemoryEngine struct {
	shards    []*engineShard
	numShards int
	encryptor     *encryption.Encryptor
	useEncryption bool
	debug         bool
	dumpPath      string
}

type engineShard struct {
	data map[string]*KeyData
	mu   sync.RWMutex
}

type DumpData struct {
	Version   int                     `json:"version"`
	Timestamp time.Time              `json:"timestamp"`
	Shards    map[int]map[string]*KeyData `json:"shards"`
}

func NewMemoryEngine(cfg *config.Config) (*MemoryEngine, error) {
	var encryptor *encryption.Encryptor
	var err error

	// Setup dump directory
	dumpPath := "data/dump"  // Default path
	if cfg.DumpPath != "" {
		dumpPath = cfg.DumpPath
	}

	if cfg.EnableEncryption {
		if cfg.EncryptionKey == "" {
			return nil, fmt.Errorf("encryption enabled but no key provided")
		}
		encryptor, err = encryption.NewEncryptor(cfg.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize encryptor: %v", err)
		}
		if cfg.Debug {
			log.Printf("Encryption enabled with key length: %d", len(cfg.EncryptionKey))
		}
	}

	numShards := runtime.NumCPU() * 2
	shards := make([]*engineShard, numShards)
	for i := 0; i < numShards; i++ {
			shards[i] = &engineShard{
				data: make(map[string]*KeyData),
				mu:   sync.RWMutex{},
			}
	}

	me := &MemoryEngine{
		shards:        shards,
		numShards:     numShards,
		encryptor:     encryptor,
		useEncryption: cfg.EnableEncryption,
		debug:         cfg.Debug,
		dumpPath:      dumpPath,
	}

	if cfg.DumpMemoryOn {
		// Ensure minimum dump interval
		dumpInterval := cfg.DumpMemoryEverySecond
		if dumpInterval < 1 {
			dumpInterval = 1 // Set minimum interval to 1 second
		}

		if cfg.Debug {
			log.Printf("Memory dump enabled. Path: %s, Interval: %d seconds", 
				dumpPath, dumpInterval)
		}

		go func() {
			ticker := time.NewTicker(time.Duration(dumpInterval) * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				if err := me.DumpToDisk(); err != nil {
					log.Printf("Failed to dump memory: %v", err)
				} else if cfg.Debug {
					log.Printf("Successfully dumped memory to disk")
				}
			}
		}()

		// Create dump directory if it doesn't exist
		if err := os.MkdirAll(dumpPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create dump directory: %v", err)
		}

		if cfg.RestoreMemoryDumpAtStart {
			if err := me.RestoreFromDisk(); err != nil {
				log.Printf("Failed to restore memory dump: %v", err)
			} else if cfg.Debug {
				log.Printf("Successfully restored memory from disk")
			}
		}
	}

	return me, nil
}

func (me *MemoryEngine) getShard(key string) *engineShard {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return me.shards[hash.Sum32()%uint32(me.numShards)]
}

func (me *MemoryEngine) Set(key string, value interface{}) error {
	if me.debug {
		log.Printf("Setting key %s with value type: %T", key, value)
	}

	// Convert value to JSON bytes
	var jsonData []byte
	var err error

	if value == nil {
		jsonData = []byte("null")
	} else if str, ok := value.(string); ok {
		// If it's already a JSON string, use it directly
		if (strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")) ||
		   (strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")) {
			jsonData = []byte(str)
		} else {
			jsonData, err = json.Marshal(str)
		}
	} else {
		jsonData, err = json.Marshal(value)
	}
	
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	var dataToStore []byte
	if me.useEncryption && me.encryptor != nil {
		if me.debug {
			log.Printf("Encrypting data for key: %s", key)
		}
		dataToStore, err = me.encryptor.Encrypt(jsonData)
		if err != nil {
			return fmt.Errorf("encryption failed: %v", err)
		}
	} else {
		dataToStore = jsonData
	}

	shard := me.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = &KeyData{
		Value:     dataToStore,
		ExpiresAt: time.Time{},
	}

	return nil
}

func (me *MemoryEngine) SetWithTTL(key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		return errors.New("TTL must be positive")
	}

	shard := me.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	
	shard.data[key] = &KeyData{
		Value:     value,
		ExpiresAt: expiresAt,
	}
	
	return nil
}

func (me *MemoryEngine) Get(key string) ([]byte, error) {
	if me.debug {
		log.Printf("Getting key: %s", key)
	}

	shard := me.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	data, exists := shard.data[key]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if !data.ExpiresAt.IsZero() && data.ExpiresAt.Before(time.Now()) {
		delete(shard.data, key)
		return nil, ErrKeyNotFound
	}

	if me.useEncryption && me.encryptor != nil {
		if me.debug {
			log.Printf("Decrypting data for key: %s", key)
		}
		decrypted, err := me.encryptor.Decrypt(data.Value)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %v", err)
		}
		return decrypted, nil
	}

	return data.Value, nil
}

func (me *MemoryEngine) GetByPattern(pattern string) ([]Match, error) {
	if me.debug {
		log.Printf("Getting keys by pattern: %s", pattern)
	}

	var matches []Match
	re, err := regexp.Compile(strings.ReplaceAll(strings.ReplaceAll(pattern, "*", ".*"), "?", "."))
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %v", err)
	}

	// Search through all shards
	for i := 0; i < me.numShards; i++ {
		shard := me.shards[i]
		shard.mu.RLock()

		for key, data := range shard.data {
			if re.MatchString(key) {
				var value []byte
				if me.useEncryption && me.encryptor != nil {
					if me.debug {
						log.Printf("Decrypting matched key %s, length: %d", key, len(data.Value))
					}
					decrypted, err := me.encryptor.Decrypt(data.Value)
					if err != nil {
						shard.mu.RUnlock()
						return nil, fmt.Errorf("failed to decrypt value for key %s: %v", key, err)
					}
					value = decrypted
				} else {
					value = data.Value
				}

				matches = append(matches, Match{
					Key:   key,
					Value: string(value),
				})
			}
		}
		shard.mu.RUnlock()
	}

	return matches, nil
}

func (me *MemoryEngine) Delete(key string) error {
	shard := me.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.data[key]; !exists {
		return ErrKeyNotFound
	}

	delete(shard.data, key)
	return nil
}

func (me *MemoryEngine) TTL(key string) (time.Duration, error) {
	shard := me.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	data, exists := shard.data[key]
	if !exists {
		return -2 * time.Second, nil // Key does not exist
	}

	if data.ExpiresAt.IsZero() {
		return -1 * time.Second, nil // Key exists but has no expiry
	}

	ttl := time.Until(data.ExpiresAt)
	if ttl <= 0 {
		// Key has expired, delete it immediately
		delete(shard.data, key)
		return -2 * time.Second, nil // Return -2 for non-existent key
	}

	return ttl, nil
}

func (me *MemoryEngine) DumpToDisk() error {
	if err := os.MkdirAll(me.dumpPath, 0755); err != nil {
		return fmt.Errorf("failed to create dump directory: %v", err)
	}

	dump := DumpData{
		Version:   1,
		Timestamp: time.Now(),
		Shards:    make(map[int]map[string]*KeyData),
	}

	for i, shard := range me.shards {
		shard.mu.RLock()
		shardData := make(map[string]*KeyData)
		for k, v := range shard.data {
			if !v.ExpiresAt.IsZero() && v.ExpiresAt.Before(time.Now()) {
				continue
			}
			// Create a deep copy of the KeyData
			shardData[k] = &KeyData{
				Value:     append([]byte(nil), v.Value...),
				ExpiresAt: v.ExpiresAt,
			}
		}
		dump.Shards[i] = shardData
		shard.mu.RUnlock()
	}

	tmpFile := filepath.Join(me.dumpPath, "memory.dump.tmp")
	file, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create dump file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(dump); err != nil {
		return fmt.Errorf("failed to encode dump: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	finalPath := filepath.Join(me.dumpPath, "memory.dump")
	if err := os.Rename(tmpFile, finalPath); err != nil {
		return fmt.Errorf("failed to rename dump file: %v", err)
	}

	if me.debug {
		log.Printf("Successfully dumped memory to %s", finalPath)
	}

	return nil
}

func (me *MemoryEngine) RestoreFromDisk() error {
	dumpPath := filepath.Join(me.dumpPath, "memory.dump")
	file, err := os.OpenFile(dumpPath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open dump file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var dump DumpData
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&dump); err != nil {
		return fmt.Errorf("failed to decode dump: %w", err)
	}

	// Clear existing data and restore from dump
	for i, shard := range me.shards {
		shard.mu.Lock()
		shard.data = make(map[string]*KeyData)
		if shardData, ok := dump.Shards[i]; ok {
			for k, v := range shardData {
				// Skip expired keys
				if !v.ExpiresAt.IsZero() && v.ExpiresAt.Before(time.Now()) {
					continue
				}
				shard.data[k] = v
			}
		}
		shard.mu.Unlock()
	}

	return nil
}

func (me *MemoryEngine) ResetMemory() error {
	// Lock all shards while resetting
	for _, shard := range me.shards {
		shard.mu.Lock()
		shard.data = make(map[string]*KeyData)
		shard.mu.Unlock()
	}
	return nil
}

