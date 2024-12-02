# JsonStorage PHP Client

A PHP client library for interacting with the JsonStorage TCP service. This client provides a robust interface for storing and retrieving JSON data through a TCP connection.

## Features

- TCP-based communication with JsonStorage server
- JSON data storage and retrieval
- Pattern matching support
- TTL (Time To Live) functionality
- Interactive testing mode
- Performance testing capabilities
- Comprehensive error handling

## Requirements

- PHP 7.4 or higher
- Socket extension enabled
- Running JsonStorage TCP server (default: localhost:5555)

## Installation

1. Clone the repository
2. Ensure the socket extension is enabled in your PHP configuration
3. Include the JsonStorageConnectionClass.php in your project

## Basic Usage

The client provides a simple connection class for interacting with the storage server:

```php
require_once 'JsonStorageConnectionClass.php';

$storage = new JsonStorageConnection('localhost', 5555);
$storage->connect();

// Store a value
$storage->set('test:key', 'Hello World');

// Retrieve a value
$value = $storage->get('test:key');

// Close the connection
$storage->close();
```

## Test Cases

The client includes comprehensive test cases for different data types:

```php:adaptors/php/index.php
startLine: 16
endLine: 33
```

## TTL (Time To Live) Testing

Example of TTL functionality testing:

```php:adaptors/php/ttl_test.php
startLine: 11
endLine: 33
```

## Error Handling

The client implements comprehensive error handling:

- Socket connection errors
- Read/Write operation failures
- Data encoding/decoding issues
- Timeout handling

## Available Commands

- `SET key value [EX seconds]` - Store a value with optional expiration
- `GET key` - Retrieve a value
- `DELETE key` - Delete a value
- `TTL key` - Get remaining time to live for a key
- `PING` - Test server connection
- Pattern matching using `*` in GET/DELETE commands

## Testing Tools

### 1. Basic Testing

Run basic functionality tests:

```bash
php index.php
```

### 2. TTL Testing

Test Time-To-Live functionality:

```bash
php ttl_test.php
```

### 3. Interactive Testing

Use the interactive test client:

```bash
php interactive_test.php
```

### 4. Performance Testing

Run performance benchmarks:

```bash
php performance_json_test.php
```

## Development

The project structure:

```
adaptors/php/
├── JsonStorageConnectionClass.php  # Main client class
├── index.php                       # Basic usage example
├── interactive_test.php            # Interactive testing tool
├── ttl_test.php                    # TTL functionality tests
├── performance_json_test.php       # Performance testing suite
└── raw_test.php                    # Raw socket testing
```

## Data Types Support

The client supports various data types:

- Simple strings
- Numbers
- JSON objects
- Arrays
- Nested structures

## Pattern Matching

The client supports pattern matching for keys using the wildcard character (\*). Example:

- `GET user:*` - Retrieves all keys starting with "user:"
- `DEL test:*` - Deletes all keys starting with "test:"

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see LICENSE file for details

# Authentication

The server requires authentication by default. You can authenticate in two ways:

1. Using environment variable:

```bash
export SERVER_PASSWORD=your_password
php raw_test.php
```

2. Interactive authentication:

```bash
php interactive_test.php
# The client will prompt for password
```

All commands after authentication work as before:

```php
$storage = new JsonStorageConnection('localhost', 5555);
$storage->connect();
$storage->authenticate('your_password'); // New authentication method
$storage->set('test:key', 'Hello World');
```

# Performance Test

## SET

Performance Summary:
Total Time: 7.90 seconds
Average Time per Record: 0.0790 ms
Initial Memory: 537.16 KB
Peak Memory Usage: 615.72 KB
Successful Records: 100000/100000
Operations per Second: 12,661.81

## GET

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

```bash
php performance_json_test.php
```
