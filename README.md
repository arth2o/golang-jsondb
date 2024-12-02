# JsonDB

A high-performance, in-memory JSON database server written in Go with PHP client support.

## Overview

JsonDB is a lightweight, TCP-based JSON storage server that provides fast in-memory storage for JSON data with optional encryption. It's designed to be a simple, fast alternative to Redis for JSON-specific use cases.

## Features

- TCP server implementation
- In-memory storage engine
- Key-based JSON storage
- TTL support for keys
- Optional encryption
- Environment-based configuration
- Multiple database support (0-4)
- PHP client library included
- Connection pooling
- Concurrent access support

## Requirements

- Go 1.21 or higher
- For PHP client: PHP 7.4 or higher
- Environment configuration file

## Installation

1. Clone the repository:

## Configuration

Create a `.env.development` file (or `.env.production` for production):

```bash
## copy the env example and fill up with correct details
cp env.example .env.development
```

### Configuration Options

- `PORT`: Server listening port (default: 5555)
- `SERVER_PASSWORD`: Authentication password
- `ENCRYPTION_KEY`: 32-byte key for data encryption
- `ENVIRONMENT`: development/production/testing
- `ENABLE_ENCRYPTION`: true/false
- `MAX_CONNECTIONS`: Maximum concurrent connections (-1 for unlimited)
- `DUMP_MEMORY_ON`: Enable/disable memory dumping functionality (true/false)
- `DUMP_MEMORY_EVERY_SECOND`: Interval in seconds between memory dumps
- `RESTORE_MEMORY_DUMP_AT_START`: Restore last memory dump when server starts (true/false)
- `DEBUG`: Enable debug mode for additional logging (true/false)

### Memory Persistence

When enabled, JsonDB can persist its in-memory data through memory dumps:

- Memory dumps are created automatically based on the configured interval
- Data can be automatically restored when the server restarts
- Useful for development and scenarios requiring data persistence without a full database

To enable memory persistence:

```env
DUMP_MEMORY_ON=true
DUMP_MEMORY_EVERY_SECOND=2  # Dumps every 2 seconds
RESTORE_MEMORY_DUMP_AT_START=true
```

## Usage

### Starting the Server

```bash
make run
or
cd jsondb
go build -o bin/server cmd/server/main.go
./bin/server
```

### Basic Commands

```bash
# Test server connection
PING                                    # Returns: PONG

# Key-Value Operations
SET key value                          # Store a value
SET key value EX seconds              # Store with expiration time
GET key                               # Retrieve a value
DEL key                               # Delete a key (alias for DELETE)
DELETE key                            # Delete a key

# Key Pattern Matching
KEYS pattern                          # Find keys matching pattern
KEYS user:*                          # Example: Find all keys starting with "user:"

# TTL (Time To Live) Operations
TTL key                               # Get remaining time to live
                                     # Returns:
                                     #   -2 if key doesn't exist
                                     #   -1 if key exists but has no expiry
                                     #   seconds remaining if key has expiry
EXPIRE key seconds                    # Set expiration time on existing key
                                     # Returns:
                                     #   1 if timeout was set
                                     #   0 if key doesn't exist

# Memory Management
RESET_MEMORY                          # Clear all stored data

# Persistence Operations
DUMP                                  # Manually trigger a memory dump to disk
                                     # Returns: OK on success
RESTORE                              # Manually restore data from the latest dump
                                     # Returns: OK on success, or error message
```

Examples:

```bash
# Basic Storage
SET user:1 "John Doe"
> OK

# Storage with TTL
SET session:abc "session_data" EX 3600
> OK

# Pattern Matching
KEYS user:*
> ["user:1", "user:2", "user:3"]

# TTL Operations
TTL session:abc
> 3600

EXPIRE user:1 7200
> 1

# Manual Persistence
DUMP
> OK

RESTORE
> OK
```

Note: The DUMP and RESTORE commands are particularly useful when `DUMP_MEMORY_ON` is set to false in your configuration, allowing manual control over data persistence.

# Run Go Test

```bash
./test_go.sh
```

# Test using telnet

```bash
# run the server
cd jsondb
go build -o bin/server cmd/server/main.go
./bin/server
# test
telnet localhost 5555
```

### PHP Client Usage

```php
require_once 'JsonStorageConnectionClass.php';

$storage = new JsonStorageConnection('localhost', 5555);
$storage->connect();

// Store JSON data
$storage->set('user:1', ['name' => 'John', 'age' => 30]);

// Retrieve data
$user = $storage->get('user:1');

// Set TTL
$storage->set('temporary:key', 'value', 3600); // Expires in 1 hour

// Delete key
$storage->del('user:1');
```

## Testing

Run the test suite:

```bash
make test
```

## Performance Testing

The project includes performance testing scripts for both SET and GET operations:

```bash
php adaptors/php/performance_set_test.php
php adaptors/php/performance_get_test.php
```

## Project Structure

```
jsondb/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── server/
│   └── storage/
├── adaptors/
│   └ php/
├── bin/
└── Makefile
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by Redis and DragonFly DB
- Built with Go's standard library
- Uses AES encryption for secure storage

## Support

For support, please open an issue in the GitHub repository.

## TTL (Time To Live) Commands

### SET with Expiration

Sets a key with a specified expiration time in seconds.

```bash
SET key value EX seconds
```

Example:

```bash
SET user:session:123 "session_data" EX 3600  # Expires in 1 hour
```

### TTL

Get the remaining time to live of a key in seconds.

```bash
TTL key
```

Returns:

- `-2` if the key does not exist or has expired
- `-1` if the key exists but has no expiration set
- Remaining time in seconds if the key exists and has an expiration

Example:

```bash
TTL user:session:123  # Returns remaining seconds, or -1/-2
```

### EXPIRE

Set a timeout on key. After the timeout has expired, the key will automatically be deleted.

```bash
EXPIRE key seconds
```

Returns:

- `1` if the timeout was set
- `0` if the key does not exist or the timeout could not be set

Example:

```bash
EXPIRE user:session:123 3600  # Set to expire in 1 hour
```

### TTL Response Examples

```bash
# Key with TTL
SET mykey "value" EX 60
TTL mykey           # Returns 60 (or less, depending on time elapsed)

# Key without TTL
SET mykey "value"
TTL mykey           # Returns -1

# Non-existent or expired key
TTL nonexistentkey  # Returns -2
```

### Notes

- TTL precision is in seconds
- When using SET with EX, the expiration time must be a positive integer
- Keys are automatically deleted once they expire
- Expired keys are removed when accessed by commands like GET, TTL, or pattern matching

## JsonDB Commands Reference

### PING

Test server connection

```
SET key value
SET key value EX seconds
```

Examples:

```
SET user:1 "John Doe"
> OK

SET session:abc "data" EX 3600  # Expires in 1 hour
> OK
```

### GET

Get value of key

```
GET key
```

Example:

```
GET user:1
> "John Doe"

GET nonexistent
> nil
```

### DELETE

Delete a key

```
DELETE key
```

Example:

```
DELETE user:1
> OK
```

### KEYS

Find all keys matching the given pattern

```
KEYS pattern
```

Examples:

```
KEYS user:*      # All keys starting with "user:"
> ["user:1", "user:2"]

KEYS session:*   # All keys starting with "session:"
> ["session:abc", "session:def"]
```

### TTL

Get the remaining time to live of a key

```
TTL key
```

Returns:

- `-2` if key does not exist
- `-1` if key exists but has no expiry
- Remaining time in seconds if key exists with expiry

Examples:

```
TTL session:abc
> 3598  # Seconds remaining

TTL user:1
> -1    # No expiration

TTL unknown
> -2    # Key doesn't exist
```

### EXPIRE

Set a timeout on key

```
EXPIRE key seconds
```

Example:

```
EXPIRE session:abc 7200  # Set to expire in 2 hours
> 1  # Success

EXPIRE nonexistent 60
> 0  # Key doesn't exist
```

````

This documentation is based on the test cases and implementations found in:


```57:138:internal/server/server_test.go
        {
            name:    "PING Command",
            command: "PING",
            want:    "PONG",
        },
        {
            name:    "SET Command",
            command: `SET test:key "hello world"`,
            want:    "OK",
        },
        {
            name:    "GET Command",
            command: "GET test:key",
            want:    "hello world",
            setup: func() error {
                return srv.Engine.Set("test:key", "hello world")
            },
        },
        {
            name:    "DELETE Command",
            command: "DELETE test:key",
            want:    "OK",
            setup: func() error {
                return srv.Engine.Set("test:key", "hello world")
            },
        },
        {
            name:    "KEYS Command with Exact Match",
            command: "KEYS test:key",
            want:    `["test:key"]`,
            setup: func() error {
                return srv.Engine.Set("test:key", "hello world")
            },
        },
        {
            name:    "KEYS Command with Pattern",
            command: "KEYS test:*",
            want:    `["test:key1","test:key2"]`,
            setup: func() error {
                if err := srv.Engine.Set("test:key1", "value1"); err != nil {
                    return err
                }
                return srv.Engine.Set("test:key2", "value2")
            },
        },
        {
            name:    "TTL Command - No Expiry",
            command: "TTL test:key",
            want:    "-1",
            setup: func() error {
                return srv.Engine.Set("test:key", "value")
            },
        },
        {
            name:    "EXPIRE Command",
            command: "EXPIRE test:key 60",
            want:    "1",
            setup: func() error {
                return srv.Engine.Set("test:key", "value")
            },
        },
        {
            name:    "TTL Command - With Expiry",
            command: "TTL test:key",
            want:    "60",
            setup: func() error {
                return srv.Engine.SetWithTTL("test:key", []byte("value"), 60*time.Second)
            },
        },
        {
            name:    "SET Command with TTL",
            command: `SET test:key "hello world" EX 60`,
            want:    "OK",
        },
        {
            name:    "Verify TTL after SET",
            command: "TTL test:key",
            want:    "60",
            setup: func() error {
                return srv.Engine.SetWithTTL("test:key", []byte("hello world"), 60*time.Second)
            },
        },
````

### Notes

- TTL precision is in seconds
- When using SET with EX, the expiration time must be a positive integer
- Keys are automatically deleted once they expire
- Expired keys are removed when accessed by commands like GET, TTL, or pattern matching

# Performance Test

I did that via PHP and in an 4 CPU laptop so the result not reflect the real performance.
I can imagine this could be much faster based on the CPU and hardware.
Later I will do more tests with Go and see the results.

## Result of SET test

Performance Summary:
Total Time: 7.90 seconds
Average Time per Record: 0.0790 ms
Initial Memory: 537.16 KB
Peak Memory Usage: 615.72 KB
Successful Records: 100000/100000
Operations per Second: 12,661.81

## Result of GET test

Performance Summary:
Total Time: 6.93 seconds
Average Time per Record: 0.0693 ms
Initial Memory: 537.41 KB
Peak Memory Usage: 616.25 KB
Successful Records: 100000/100000
Operations per Second: 14,428.47

Data Type Distribution:
string: 25000 (25.0%)
json: 50000 (50.0%)
number: 25000 (25.0%)
complex: 0 (0.0%)

# Performance json test

GET Operations:
Total Time: 7.74 seconds
Average Time per GET: 0.0774 ms
Success Rate: 100.00%

Overall Statistics:
Total Time: 20.04 seconds
Initial Memory: 570.12 KB
Peak Memory Usage: 648.61 KB
Operations per Second: 9,980.88

Data Type Distribution:
user: 20026 (20.0%)
product: 19988 (20.0%)
order: 19984 (20.0%)
log: 19947 (19.9%)
metrics: 20055 (20.1%)
