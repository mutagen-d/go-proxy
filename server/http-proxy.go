package server

import (
	"fmt"
	"go-proxy/tools"
	"io"
	"net"
	"net/url"
	"strings"
)

func HandleHTTP(conn net.Conn, data []byte, createConnection CreateConnection) {
	// HTTP proxy
	if createConnection == nil {
		createConnection = directConnection
	}
	req := tools.ParseRequest(data)
	srcAddr := conn.RemoteAddr().String()
	// fmt.Printf("%v request:\n%v\n", tools.Now(), string(data))
	if strings.ToUpper(req.Method) == "CONNECT" {
		logConnection(req.Url, srcAddr, "http")
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
	logConnection(origin, srcAddr, "http")
	target, err := createConnection(origin, srcAddr)
	if err != nil {
		fmt.Printf("%v Error %v %v\n", tools.Now(), err, origin)
		tools.Send(conn, 500, "Internal Server Error", fmt.Sprintf("connect error (url = %v)", req.Url))
		return
	}
	fixedData := tools.FixFast(data, origin)
	// fmt.Printf("%v write 1 (%v - %v):\n%v\n", tools.Now(), len(fixedData), len(data), string(fixedData))
	_, err = target.Write(fixedData)
	if err != nil {
		fmt.Printf("%v Error %v %v\n", tools.Now(), err, origin)
		tools.Send(conn, 500, "Internal Server Error", fmt.Sprintf("write error (url = %v)", req.Url))
		return
	}
	defer target.Close()
	go io.Copy(target, &tools.FixHttp{Reader: conn})
	io.Copy(conn, &tools.FixHttp{Reader: target})
}
