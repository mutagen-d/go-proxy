# SSH Tunnel Proxy Server

A lightweight Go proxy server that tunnels traffic through SSH, supporting HTTP/HTTPS and SOCKS4/SOCKS4A/SOCKS5 protocols simultaneously on the same port, with password and key-based authentication.

## Features

- SSH tunneling for secure proxy connections
- Support for both password and private key authentication
- Multi-protocol support - Single server handles all protocols:
    - HTTP/HTTPS proxy
    - SOCKS4 proxy
    - SOCKS4A proxy (domain name support)
    - SOCKS5 proxy (with authentication support)
- Protocol auto-detection - Automatically detects and handles the appropriate protocol
- ️Configurable via command-line flags
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

### Supported Protocols

The proxy server automatically detects and handles all these protocols on the same port:
| Protocol | Features | Auto-detected |
|----------|----------|---------------|
| HTTP | Standard HTTP proxy with CONNECT method for HTTPS | ✅ |
| SOCKS4 | Basic SOCKS4 without authentication | ✅ |
| SOCKS4A | SOCKS4 with domain name support (no authentication) | ✅ |
| SOCKS5 | Full SOCKS5 with optional authentication | ✅ |

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

### Client Configuration Examples

#### HTTP Proxy
```bash
# Using curl with HTTP proxy
curl -x http://localhost:8080 https://api.example.com

# Using wget
wget -e use_proxy=yes -e http_proxy=localhost:8080 https://example.com
```

#### SOCKS4/SOCKS4A
```bash
# Using curl with SOCKS4
curl --socks4 localhost:8080 https://api.example.com

# Using curl with SOCKS4A (domain names supported)
curl --socks4a localhost:8080 https://api.example.com

# Using ssh through SOCKS4
ssh -o ProxyCommand="nc -X 4 -x localhost:8080 %h %p" user@example.com
```

#### SOCKS5
```bash
# Using curl with SOCKS5
curl --socks5 localhost:8080 https://api.example.com
# or
curl -x socks5://localhost:8080 https://api.example.com
# or
curl -x socks5h://localhost:8080 https://api.example.com

# Using SSH through SOCKS5
ssh -o ProxyCommand="nc -X 5 -x localhost:8080 %h %p" user@example.com
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