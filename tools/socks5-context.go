package tools

import (
	"fmt"
	"net"
	"slices"
	"strings"
)

type Socks5Context struct {
	Conn net.Conn
	AuthType []byte
	AddressType byte
	Address []byte
}
func (v *Socks5Context) Fail(err error, code int) bool {
	if err == nil {
		return false
	}
	srcIp := strings.Split(v.Conn.RemoteAddr().String(), ":")[0]
	fmt.Printf("%v %-12v %02v %v\n", Now(), srcIp, code, err)
	v.Conn.Write([]byte{0x05, 0x01}) // FAILED
	return true
}
func (v *Socks5Context) XFail(err error, code int) bool {
	if err == nil {
		return false
	}
	srcIp := strings.Split(v.Conn.RemoteAddr().String(), ":")[0]
	fmt.Printf("%v %-12v %02v %v\n", Now(), srcIp, code, err)
	v.Conn.Write([]byte{0x05, 0x01, 0x00}) // FAILED
	return true
}
func (v *Socks5Context) NotSupportedAuth(auth byte, code int) bool {
	if slices.Contains(v.AuthType, auth) {
		return false
	}
	err := fmt.Errorf("not supported authentication, recv = %v, sup = %v", v.AuthType, auth)
	srcIp := strings.Split(v.Conn.RemoteAddr().String(), ":")[0]
	fmt.Printf("%v %-12v %02v %v\n", Now(), srcIp, code, err)
	v.Conn.Write([]byte{0x05, 0xFF}) // NOT SUPPORTED AUTHENTICATION 
	return true
}
func (v *Socks5Context) PwdAuth() {
	v.Conn.Write([]byte{0x05, 0x02})
}
func (v *Socks5Context) NoAuth() {
	v.Conn.Write([]byte{0x05, 0x00})
}
func (v *Socks5Context) AuthSuccess() {
	v.Conn.Write([]byte{0x01, 0x00})
}
func (v *Socks5Context) AuthFailed(code int) bool {
	err := fmt.Errorf("Auth Failed")
	return v.Fail(err, code)
}
func (v *Socks5Context) AddressTypeUnsupported() {
	v.Conn.Write([]byte{0x05, 0x08, 0x00}) // ADDRESS TYPE UNSUPPORTED
}
func (v *Socks5Context) Success() {
	v.Conn.Write(append([]byte{0x05, 0x00, 0x00, v.AddressType}, v.Address...)) // SUCCESS
}