package main

import (
	"fmt"
	"go-proxy/server"
	"go-proxy/tools"
	"log"
	"net"
)

var sshTunnel = server.SSHTunnel{Host: "mysite.com"}

func sshConnection(dstAddr, srcAddr string) (net.Conn, error) {
	return sshTunnel.Forward(dstAddr, srcAddr)
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("error: ", err)
	}
	fmt.Printf("%v server listening on %v\n", tools.Now(), listener.Addr())
	defer listener.Close()
	go func() {
		_, err = sshTunnel.Connect()
		if err != nil {
			panic(err)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%v Warning! Cant accept connection: %v\n", tools.Now(), err)
			continue
		}
		go server.HandleConnection(conn, sshConnection)
	}
}
