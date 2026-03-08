package server

import (
	"fmt"
	"go-proxy/tools"
	"net"
	"slices"
	"strings"
)

type CreateConnection func(dstAddr, srcAddr string) (net.Conn, error)

func HandleConnection(conn net.Conn, createConnection CreateConnection) {
	defer conn.Close()
	buf := make([]byte, 16384)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("%v Error %v\n", tools.Now(), err)
		tools.Send(conn, 500, "Internal Server Error", "read error")
		return
	}
	if isSocks4(buf[:n]) {
		Socks4Proxy(conn, buf[:n], createConnection)
	} else if isSocks5(buf[:n]) {
		Socks5Proxy(conn, buf[:n], createConnection)
	} else if isHttp(buf[:n]) {
		HandleHTTP(conn, buf[:n], createConnection)
	} else {
		// TODO
		fmt.Printf("%v error!\n", tools.Now())
	}
}

func logConnection(dstAddr, srcAddr string, proto string) {
	connType := ConnFilter.Filter(dstAddr)
	ip := strings.Split(srcAddr, ":")[0]
	fmt.Printf("%s %-8s %-8s %-12s %v\n", tools.Now(), connType, strings.ToUpper(proto), ip, dstAddr)
}

func directConnection(dstAddr, srcAddr string) (net.Conn, error) {
	return net.Dial("tcp", dstAddr)
}

func isSocks5(data []byte) bool {
	if data[0] != 0x05 {
		return false
	}
	if len(data) < 2 {
		return false
	}
	size := data[1]
	if len(data) < (2 + int(size)) {
		return false
	}
	auth := data[2:2+size]
	return slices.Contains(auth, 0x00) || slices.Contains(auth, 0x02)
}

func isSocks4(data []byte) bool {
	if data[0] != 0x04 {
		return false
	}
	if len(data) < 2 {
		return false
	}
	return data[1] == 0x01 || data[1] == 0x02 || data[1] == 0x03
}

func isHttp(data []byte) bool {
	req := tools.ParseRequest(data)
	return tools.IsSupportedMethod(req.Method)
}