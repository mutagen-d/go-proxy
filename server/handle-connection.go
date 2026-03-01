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

func HandleHTTP(conn net.Conn, data []byte, createConnection CreateConnection) {
	// HTTP proxy
	if createConnection == nil {
		createConnection = func(dstAddr, srcAddr string) (net.Conn, error) {
			return net.Dial("tcp", dstAddr)
		}
	}
	req := tools.ParseRequest(data)
	srcAddr := conn.RemoteAddr().String()
	if strings.ToUpper(req.Method) == "CONNECT" {
		fmt.Printf("%v proxy %v %v\n", tools.Now(), conn.RemoteAddr(), req.Url)
		target, err := createConnection(req.Url, srcAddr)
		if err != nil {
			fmt.Printf("%v Error %v\n", tools.Now(), err)
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
		fmt.Printf("%v Error %v\n", tools.Now(), err)
		tools.Send(conn, 500, "Internal Server Error", fmt.Sprintf("invalid url (url = %v)", req.Url))
		return
	}
	host := u.Host
	port := u.Port()
	if port == "" {
		port = "80"
	}
	origin := fmt.Sprintf("%v:%v", host, port)
	fmt.Printf("%v proxy %v %v %v\n", tools.Now(), conn.RemoteAddr(), conn.LocalAddr(), origin)
	target, err := createConnection(origin, srcAddr)
	if err != nil {
		fmt.Printf("%v Error %v\n", tools.Now(), err)
		tools.Send(conn, 500, "Internal Server Error", fmt.Sprintf("connect error (url = %v)", req.Url))
		return
	}
	defer target.Close()
	target.Write(data)
	go io.Copy(target, conn)
	io.Copy(conn, target)
}

