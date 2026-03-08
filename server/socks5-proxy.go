package server

import (
	"fmt"
	"go-proxy/tools"
	"io"
	"net"
	"strings"
)

func Socks5Proxy(conn net.Conn, data []byte, createConnection CreateConnection) {
	if createConnection == nil {
		createConnection = directConnection
	}
	reader := tools.NewBytesReader(data)
	_, err := reader.Seek(1, io.SeekStart)
	ctx := tools.Socks5Context{Conn: conn}
	if ctx.Fail(err, 1) {
		return
	}
	if ctx.Fail(err, 2) {
		return
	}
	size, err := reader.ReadUint8()
	if ctx.Fail(err, 3) {
		return
	}
	authType, err := reader.Read(int(size))
	if ctx.Fail(err, 4) {
		return
	}
	ctx.AuthType = authType
	isAuthRequired := false
	if isAuthRequired && ctx.NotSupportedAuth(0x02, 5) {
		return
	}
	if !isAuthRequired && ctx.NotSupportedAuth(0x00, 6) {
		return
	}
	if isAuthRequired {
		ctx.PwdAuth()
		onSocks5PasswordAuth(conn, createConnection)
		return
	}
	ctx.NoAuth()
	buf := make([]byte, 16384)
	n, err := conn.Read(buf)
	if ctx.Fail(err, 7) {
		return
	}
	onSocks5Connection(conn, buf[:n], createConnection)
}

func onSocks5PasswordAuth(conn net.Conn, createConnection CreateConnection) {
	buf := make([]byte, 16384)
	n, err := conn.Read(buf)
	ctx := tools.Socks5Context{Conn: conn}
	if ctx.Fail(err, 8) {
		return
	}
	reader := tools.NewBytesReader(buf[:n])
	reader.Seek(1, io.SeekStart)
	useridSize, err := reader.ReadUint8()
	if ctx.Fail(err, 9) {
		return
	}
	userid, err := reader.ReadString(int(useridSize))
	if ctx.Fail(err, 10) {
		return
	}
	passwordSize, err := reader.ReadUint8()
	if ctx.Fail(err, 11) {
		return
	}
	password, err := reader.ReadString(int(passwordSize))
	if ctx.Fail(err, 12) {
		return
	}
	if userid != "" && password != "" {
		ctx.AuthSuccess()
		n, err = conn.Read(buf)
		if ctx.Fail(err, 13) {
			return
		}
		onSocks5Connection(conn, buf[:n], createConnection)
	} else {
		ctx.AuthFailed(14)
	}
}

func onSocks5Connection(conn net.Conn, data []byte, createConnection CreateConnection) {
	reader := tools.NewBytesReader(data)

	ctx := tools.Socks5Context{Conn: conn}
	_, err := reader.Seek(1, io.SeekStart)
	if ctx.XFail(err, 15) {
		return
	}
	command, err := reader.ReadUint8()
	if ctx.XFail(err, 16) {
		return
	}
	_, err = reader.Seek(1, io.SeekCurrent)
	if ctx.XFail(err, 17) {
		return
	}
	addressType, err := reader.ReadUint8()
	if ctx.XFail(err, 18) {
		return
	}
	var dstHost string
	var dstPort uint16
	var address []byte
	switch addressType {
	case 0x01:
		ip, err := reader.Read(4)
		if ctx.XFail(err, 19) {
			return
		}
		dstHost = fmt.Sprintf("%v.%v.%v.%v", ip[0], ip[1], ip[2], ip[3])
		dstPort, err = reader.ReadUint16BE()
		if ctx.XFail(err, 20) {
			return
		}
		_, err = reader.Seek(-6, io.SeekCurrent)
		if ctx.XFail(err, 21) {
			return
		}
		address, err = reader.ReadAll()
		if ctx.XFail(err, 22) {
			return
		}
	case 0x03:
		domSize, err := reader.ReadUint8()
		if ctx.XFail(err, 23) {
			return
		}
		host, err := reader.Read(int(domSize))
		if ctx.XFail(err, 24) {
			return
		}
		dstHost = string(host)
		dstPort, err = reader.ReadUint16BE()
		if ctx.XFail(err, 25) {
			return
		}
		_, err = reader.Seek(-(int64(domSize + 3)), io.SeekCurrent)
		if ctx.XFail(err, 26) {
			return
		}
		address, err = reader.ReadAll()
		if ctx.XFail(err, 27) {
			return
		}
	case 0x04:
		host, err := reader.ReadStringHex(16)
		if ctx.XFail(err, 28) {
			return
		}
		dstHost = strings.Join(host, ":")
		dstPort, err = reader.ReadUint16BE()
		if ctx.XFail(err, 29) {
			return
		}
		_, err = reader.Seek(-18, io.SeekCurrent)
		if ctx.XFail(err, 30) {
			return
		}
		address, err = reader.ReadAll()
		if ctx.XFail(err, 31) {
			return
		}
	default:
		ctx.AddressTypeUnsupported()
		return
	}
	if command == 0x01 { // CONNECT
		dstAddr := fmt.Sprintf("%v:%v", dstHost, dstPort)
		srcAddr := conn.RemoteAddr().String()
		logConnection(dstAddr, srcAddr, "socks5")
		target, err := createConnection(dstAddr, srcAddr)
		if ctx.XFail(err, 32) {
			return
		}
		defer target.Close()
		ctx.AddressType = addressType
		ctx.Address = address
		ctx.Success()
		go io.Copy(target, conn)
		io.Copy(conn, target)
	}
}
