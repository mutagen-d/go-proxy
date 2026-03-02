package server

import (
	"fmt"
	"go-proxy/tools"
	"io"
	"net"
	"net/url"
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
	HandleHTTP(conn, buf[:n], createConnection)
}

func logConnection(dstAddr, srcAddr string) {
	connType := ConnFilter.Filter(dstAddr)
	ip := strings.Split(srcAddr, ":")[0]
	fmt.Printf("%s %-8s %-12s %v\n", tools.Now(), connType, ip, dstAddr)
}

func directConnection(dstAddr, srcAddr string) (net.Conn, error) {
	return net.Dial("tcp", dstAddr)
}

func HandleHTTP(conn net.Conn, data []byte, createConnection CreateConnection) {
	// HTTP proxy
	if createConnection == nil {
		createConnection = directConnection
	}
	req := tools.ParseRequest(data)
	srcAddr := conn.RemoteAddr().String()
	if strings.ToUpper(req.Method) == "CONNECT" {
		logConnection(req.Url, srcAddr)
		target, err := createConnection(req.Url, srcAddr)
		if err != nil {
			fmt.Printf("%v Error %v %v\n", tools.Now(), err, req.Url)
			tools.Send(conn, 500, "Internal Server Error", fmt.Sprintf("connect error (url = %v)", req.Url))
			return
		}
		tools.Send(conn, 200, "OK", "")
		defer target.Close()
		go io.Copy(target, conn)
		io.Copy(conn, target)
		return
	}
	if !tools.IsSupportedMethod(req.Method) {
		fmt.Printf("%v Unsupported Method %v\n", tools.Now(), req.Method)
		tools.Send(conn, 405, "Unsupported Method", "")
		return
	}
	u, err := url.Parse(req.Url)
	if err != nil {
		fmt.Printf("%v Error %v %v\n", tools.Now(), err, req.Url)
		tools.Send(conn, 500, "Internal Server Error", fmt.Sprintf("invalid url (url = %v)", req.Url))
		return
	}
	host := u.Hostname()
	port := u.Port()
	if port == "" && (u.Scheme == "http" || u.Scheme == "ws") {
		port = "80"
	}
	if port == "" && (u.Scheme == "https" || u.Scheme == "wss") {
		port = "443"
	}
	origin := fmt.Sprintf("%v:%v", host, port)
	logConnection(origin, srcAddr)
	target, err := createConnection(origin, srcAddr)
	if err != nil {
		fmt.Printf("%v Error %v %v\n", tools.Now(), err, origin)
		tools.Send(conn, 500, "Internal Server Error", fmt.Sprintf("connect error (url = %v)", req.Url))
		return
	}
	defer target.Close()
	target.Write(data)
	go io.Copy(target, conn)
	io.Copy(conn, target)
}
