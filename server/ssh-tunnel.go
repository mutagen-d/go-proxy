package server

import (
	"fmt"
	"go-proxy/tools"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	User           string
	Host           string
	Port           int
	Password       string
	PrivateKeyPath string
	client         *ssh.Client
	mu             sync.RWMutex
}

/**
 * ssh port forwarding:
 * P.S. `dstAddr` is like "google.com:443"
 */
func (s *SSHTunnel) Forward(dstAddr, srcAddr string) (net.Conn, error) {
	_, err := s.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %v", err)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, err := s.client.Dial("tcp", dstAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial through SSH tunnel to %s: %v", dstAddr, err)
	}
	return conn, nil
}

func (s *SSHTunnel) Connect() (*ssh.Client, error) {
	s.mu.RLock()
	client := s.client
	if client == nil {
		s.mu.RUnlock()
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.connect()
	}
	ok := s.helthcheck()
	if ok {
		s.mu.RUnlock()
		return client, nil
	}
	s.mu.RUnlock()
	return s.Reconnect()
}

func (s *SSHTunnel) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.close()
}

func (s *SSHTunnel) Reconnect() (*ssh.Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.close()
	return s.connect()
}

func (s *SSHTunnel) helthcheck() bool {
	if s.client == nil {
		return false
	}
	_, _, err := s.client.SendRequest("keepalive", true, nil)
	return err == nil
}

func (s *SSHTunnel) close() {
	if s.client != nil {
		s.client.Close()
		s.client = nil
	}
}
func (s *SSHTunnel) reconnect() (*ssh.Client, error) {
	s.close()
	return s.connect()
}

func (s *SSHTunnel) connect() (*ssh.Client, error) {
	if s.client != nil {
		return s.client, nil
	}
	start := time.Now()
	var config *ssh.ClientConfig
	user := s.User
	if user == "" {
		user = "root"
	}
	port := s.Port
	if port == 0 {
		port = 22
	}
	host := s.Host
	if host == "" {
		host = "localhost"
	}
	addr := fmt.Sprintf("%v:%v", host, port)
	if s.Password != "" {
		config = &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{ssh.Password(s.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key validation
		}
	} else {
		privateKeyPath := s.PrivateKeyPath
		if privateKeyPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("homedir not founds: %v", err)
			}
			privateKeyPath = filepath.Join(home, ".ssh", "id_rsa")
		}
		key, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key from %s: %v", privateKeyPath, err)
		}
		// Parse the private key
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
		config = &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key validation
		}
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %v", err)
	}
	s.client = client
	duration := time.Now().Sub(start)
	fmt.Printf("%v ssh connected to %v (+%vms)\n", tools.Now(), addr, duration.Milliseconds())
	return client, nil
}
