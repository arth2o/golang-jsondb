# JsonDB System Documentation

A high-performance, in-memory key-value store with JSON support, pattern matching, and TTL capabilities.

## Overview

JsonDB is a TCP-based server implementation providing fast, secure storage with support for complex data structures and automatic expiration. Built with Go 1.22.2, it offers robust client libraries for PHP, Laravel, and NodeJS.

## Features

### Core Capabilities

- ✨ In-memory storage with sharding
- 🔒 Password-based authentication
- 🔐 Optional AES encryption
- ⏱️ TTL support for entries
- 🔍 Pattern matching for key lookups
- 💾 Memory dump/restore functionality
- 🚦 Rate limiting protection

### Client Libraries

- 🎯 Laravel/PHP Adaptor
- 🟢 NodeJS Adaptor
- 🐘 Pure PHP Adaptor

## Technical Specifications

### Server Architecture

#### Storage Engine

- Sharded in-memory storage
- Thread-safe operations
- Pattern matching support
- TTL management
- Persistence through memory dumps

#### Security Features

- Password authentication
- AES-256 encryption (optional)
- Rate limiting
- Connection timeouts
- Session management

### Configuration

#### Environment Support

- Development
- Production
- Testing

#### Configurable Parameters

- Port
- Server password
- Encryption settings
- Environment mode
- Connection limits
- Debug options

## Data Operations

### Basic Commands

```plaintext
SET key value [ttl]
GET key
DEL key
TTL key
PATTERN pattern
```

### Data Types Supported

- Strings
- Numbers
- JSON objects
- Binary data
- Complex nested structures

## Monitoring & Management

### Health Checks

- Connection status
- Latency monitoring
- Memory usage
- Error rates

### Persistence

- Configurable dump intervals
- Automatic restoration
- Dump file validation

## Performance Considerations

### Limitations

- Single-node architecture
- Memory-bound storage
- Synchronous operations
- TCP-only communication

### Best Practices

- Use connection pooling
- Implement retry mechanisms
- Monitor memory usage
- Regular health checks

## Testing Infrastructure

### Test Categories

- Unit tests
- Integration tests
- Performance tests
- Security tests

### Test Tools

- Mock server
- Benchmark utilities
- Coverage reporting

## Client Library Guidelines

### Implementation Standards

- Connection pooling
- Error handling
- Automatic reconnection
- Rate limiting
- Type safety

## Future Compatibility

### Stability Guarantees

- Command compatibility
- Authentication mechanism
- Data format support
- Error handling patterns
- Adaptor interface stability

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Support

For support, please use the issue tracker or contact the development team.

## Acknowledgments

- Go Team for the excellent networking stack
- Contributors and testers
- Early adopters and feedback providers

---

**Note**: This documentation is maintained by the development team. For the latest updates and changes, please refer to the changelog.

### Client Library Documentation

- [Laravel Client Documentation](docs/laravel-client.md)
- [NodeJS Client Documentation](docs/nodejs-client.md)
- [PHP Client Documentation](docs/php-client.md)
