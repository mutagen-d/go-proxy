package tools

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port        int
	SSHHost     string
	SSHHostname string
	SSHPort     int
	SSHUser     string
	SSHPassword string
	SSHKeyPath  string
}

func ParseConfig() *Config {
	// Define flags with defaults and descriptions
	var (
		port        = flag.Int("port", 8080, "Proxy server port (default: 8080)")
		sshHost     = flag.String("ssh-host", "", "SSH host and port (format: hostname:port)")
		sshUser     = flag.String("ssh-user", "root", "SSH username (default: root)")
		sshPassword = flag.String("ssh-password", "", "SSH password (optional if using key)")
		sshKeyPath  = flag.String("ssh-key", "", "Path to SSH private key (optional if using password)")
		showHelp    = flag.Bool("help", false, "Show help message")
	)

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A proxy server that tunnels through SSH.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -ssh-user=john -ssh-host=example.com:2222 -port=9090\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -ssh-user=john -ssh-key=/path/to/key -ssh-host=example.com\n", os.Args[0])
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Validate required fields
	if *sshUser == "" {
		flag.Usage()
		log.Fatal("SSH user is required (use -ssh-user flag)")
	}

	// Parse SSH host and port
	sshHostname, sshPort, err := parseHostPort(*sshHost, 22)
	if err != nil {
		flag.Usage()
		log.Fatalf("Invalid SSH host format: %v", err)
	}

	if *sshHost == "" {
		flag.Usage()
		log.Fatal("SSH host is required (use -ssh-host flag)")
	}

	// Create config
	config := &Config{
		Port:        *port,
		SSHHost:     *sshHost,
		SSHHostname: sshHostname,
		SSHPort:     sshPort,
		SSHUser:     *sshUser,
		SSHPassword: *sshPassword,
		SSHKeyPath:  *sshKeyPath,
	}

	// Display configuration
	printConfig(config)

	return config
}

func parseHostPort(hostPort string, defaultPort int) (hostname string, port int, err error) {
	parts := strings.Split(hostPort, ":")
	hostname = parts[0]
	port = defaultPort

	if len(parts) > 1 {
		port, err = strconv.Atoi(parts[1])
		if err != nil {
			return "", 0, fmt.Errorf("invalid port number: %s", parts[1])
		}
	}
	return hostname, port, nil
}

func printConfig(config *Config) {
	fmt.Println("=================================")
	fmt.Println("Proxy Server Configuration")
	fmt.Println("=================================")
	fmt.Printf("Local Port:     %d\n", config.Port)
	fmt.Printf("SSH Host:       %s\n", config.SSHHost)
	fmt.Printf("SSH Hostname:   %s\n", config.SSHHostname)
	fmt.Printf("SSH Port:       %d\n", config.SSHPort)
	fmt.Printf("SSH User:       %s\n", config.SSHUser)
	fmt.Printf("SSH Password:   %s\n", maskString(config.SSHPassword))
	fmt.Printf("SSH Key Path:   %s\n", config.SSHKeyPath)
	fmt.Println("=================================")
}

func maskString(s string) string {
	if s == "" {
		return "<not set>"
	}
	return "********"
}
