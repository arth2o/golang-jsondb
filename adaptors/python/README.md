# JsonStorage Python Client

A Python client library for interacting with the JsonStorage TCP service. This client provides a robust interface for storing and retrieving JSON data through a TCP connection.

## Features

- TCP-based communication with JsonStorage server
- JSON data storage and retrieval
- Pattern matching support
- TTL (Time To Live) functionality
- Interactive testing mode
- Comprehensive error handling

## Requirements

- Python 3.8+
- python-dotenv

## Installation

1. Clone the repository
2. Install dependencies:

```bash
pip install -r requirements.txt
```

## Basic Usage

```python
from json_storage_connection import JsonStorageConnection

# Initialize connection
storage = JsonStorageConnection('localhost', 5555)
storage.connect()

# Basic operations
storage.set("my:key", "Hello World")
value = storage.get("my:key")
storage.delete("my:key")

# TTL operations
storage.set("temp:key", "Expires in 5 seconds", 5)
ttl = storage.ttl("temp:key")
```

## Available Test Scripts

1. Basic Connection Test:

```bash
python3 test_connection.py
```

2. Raw Socket Testing:

```bash
python3 raw_test.py
```

3. Basic Commands Example:

```bash
python3 basic_commands_example.py
```

## Authentication

The server requires authentication by default. Configure it using:

1. Environment variable:

```bash
export SERVER_PASSWORD=your_password
```

2. .env file in the jsondb directory:

```
SERVER_PASSWORD=your_password
```

## Available Commands

- `SET key value [EX seconds]` - Store a value with optional expiration
- `GET key` - Retrieve a value
- `DELETE key` - Delete a value
- `TTL key` - Get remaining time to live for a key
- `EXPIRE key seconds` - Set expiration time for a key
- `PING` - Test server connection

## Example Scripts

The basic_commands_example.py demonstrates core functionality:

```python:adaptors/python/basic_commands_example.py
startLine: 8
endLine: 51
```

For more complex testing scenarios, see raw_test.py:

```python:adaptors/python/raw_test.py
startLine: 151
endLine: 177
```

## License

MIT License - see LICENSE file for details
