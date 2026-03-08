package server

import (
	"fmt"
	"go-proxy/tools"
	"io"
	"net"
	"strings"
)

func Socks4Proxy(conn net.Conn, data []byte, createConnection CreateConnection) {
	if createConnection == nil {
		createConnection = directConnection
	}
	reader := tools.NewBytesReader(data)
	ctx := tools.Socks4Context{Conn: conn}
	_, err := reader.Seek(1, io.SeekStart) // skip version byte
	if ctx.Fail(err, 1) {
		return
	}
	command, err := reader.ReadUint8()
	if ctx.Fail(err, 2) {
		return
	}
	dstPort, err := reader.ReadUint16BE()
	if ctx.Fail(err, 3) {
		return
	}
	ip, err := reader.Read(4)
	if ctx.Fail(err, 4) {
		return
	}
	_, err = reader.ReadStringNT() // read user id
	if ctx.Fail(err, 5) {
		return
	}
	dstHost := fmt.Sprintf("%v.%v.%v.%v", ip[0], ip[1], ip[2], ip[3])
	if strings.HasPrefix(dstHost, "0.0.0.") {
		dstHost, err = reader.ReadStringNT()
		if ctx.Fail(err, 6) {
			return
		}
	}
	if command == 0x01 {
		dstAddr := fmt.Sprintf("%v:%v", dstHost, dstPort)
		srcAddr := conn.RemoteAddr().String()
		logConnection(dstAddr, srcAddr, "socks4")
		target, err := createConnection(dstAddr, srcAddr)
		if ctx.FailConnect(err, 7) {
			return
		}
		ctx.Success()
		defer target.Close()
		go io.Copy(target, conn)
		io.Copy(conn, target)
	} else {
		ctx.Reject()
		return
	}
}
