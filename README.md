# SSH Tunnel Proxy Server

A lightweight Go proxy server that tunnels traffic through SSH, supporting password and key-based authentication.

## Features

- SSH tunneling for secure proxy connections
- Support for both password and private key authentication
- ️Configurable via command-line flags
- Environment variable support
- Simple and lightweight

## Installation

```bash
# Clone the repository
git clone https://github.com/mutagen-d/go-proxy.git
cd go-proxy

# Build the binary
go build -o ssh-proxy main.go

# Or install directly
go install
```

## Usage

### Basic Usage

```bash
# Using password authentication
./ssh-proxy -ssh-user=john -ssh-host=example.com:22 -ssh-password=secret

# Using SSH key authentication
./ssh-proxy -ssh-user=john -ssh-host=example.com:22 -ssh-key=/path/to/private_key

# Custom port
./ssh-proxy -ssh-user=john -ssh-host=example.com -port=9090
```

### Command-line Options

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `-port` | Local proxy server port | `8080` | No |
| `-ssh-host` | SSH server address (host:port) | - | **Yes** |
| `-ssh-user` | SSH username | `root` | No |
| `-ssh-password` | SSH password | - | No* |
| `-ssh-key` | Path to SSH private key | `~/.ssh/id_rsa` | No* |
| `-help` | Show help message | - | No |

\* If no `-ssh-password` provided, then `-ssh-key` is used

### Examples

1. **Basic password authentication:**
   ```bash
   ./ssh-proxy -ssh-user=admin -ssh-host=192.168.1.100 -ssh-password=secret123
   ```

2. **Using SSH key with custom port:**
   ```bash
   ./ssh-proxy -ssh-user=admin -ssh-host=192.168.1.100 -ssh-key=~/.ssh/id_rsa -port=9090
   ```

3. **Different SSH port:**
   ```bash
   ./ssh-proxy -ssh-user=admin -ssh-host=example.com:2222 -ssh-password=secret
   ```

## Configuration Priority

1. Command-line flags (highest priority)
2. Default values (lowest priority)

## Security Notes

- Command-line arguments may be visible in process lists
- SSH key authentication is more secure than passwords
- Always use SSH keys with passphrases in production

## Development

### Prerequisites

- Go 1.16 or higher
- SSH server for testing

### Building from Source

```bash
go mod init go-proxy
go mod tidy
go build -o ssh-proxy main.go
```

### Running Tests

```bash
go test -v ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Go standard library `flag` package
- SSH tunneling concepts from the Go crypto library

## Support

For issues and questions:
- Open an issue on GitHub
- Check the [Go flag package documentation](https://pkg.go.dev/flag)
- Review SSH best practices