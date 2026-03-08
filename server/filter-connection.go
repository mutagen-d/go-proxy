package server

import (
	"regexp"
	"strings"
)

type ConnectionType string

const (
	Direct  ConnectionType = "direct"
	Proxy   ConnectionType = "proxy"
	Blocked ConnectionType = "blocked"
)

type ConnectionFilter struct {
	Direct  []*regexp.Regexp
	Blocked []*regexp.Regexp
}

func (m *ConnectionFilter) Filter(dstAddr string) ConnectionType {
	parts := strings.Split(dstAddr, ":")
	host := parts[0]
	if m.Direct != nil {
		for _, regex := range m.Direct {
			if regex.MatchString(host) {
				return Direct
			}
		}
	}
	if m.Blocked != nil {
		for _, regex := range m.Blocked {
			if regex.MatchString(host) {
				return Blocked
			}
		}
	}
	return Proxy
}

var ConnFilter = ConnectionFilter{
	Blocked: []*regexp.Regexp{
		regexp.MustCompile(`^(localhost|127\.0\.0\.1)$`), // localhost
		regexp.MustCompile(`^(10\.\d+\.\d+\.\d+|172\.(?:1[6-9]|2[0-9]|3[01])\.\d+\.\d+|192\.168\.\d+\.\d+)$`), // private IP
	},
}
