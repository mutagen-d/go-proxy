package tools

import (
	"fmt"
	"net"
	"strings"
)

type Socks4Context struct {
	Conn net.Conn
}

func (v *Socks4Context) Fail(err error, code int) bool {
	if err == nil {
		return false
	}
	srcIp := strings.Split(v.Conn.RemoteAddr().String(), ":")[0]
	fmt.Printf("%v %-12v %02v %v\n", Now(), srcIp, code, err)
	v.response(0x5B) // REJECTED
	return true
}
func (v *Socks4Context) FailConnect(err error, code int) bool {
	if err == nil {
		return false
	}
	srcIp := strings.Split(v.Conn.RemoteAddr().String(), ":")[0]
	fmt.Printf("%v %-12v %02v %v\n", Now(), srcIp, code, err)
	v.response(0x5C) // REJECTED - cannot connect
	return true
}
func (v *Socks4Context) Success() {
	v.response(0x5A)
}
func (v *Socks4Context) Reject() {
	v.response(0x5B)
}
func (v *Socks4Context) NoUserId() {
	v.response(0x5D)
}
func (v *Socks4Context) response(code byte) {
	v.Conn.Write([]byte{0, code, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01})
}