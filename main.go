package main

import (
	"fmt"
	"go-proxy/server"
	"go-proxy/tools"
	"log"
	"net"
)

func main() {
	config := tools.ParseConfig()
	var sshTunnel = server.SSHTunnel{
		User: config.SSHUser,
		Host: config.SSHHostname,
		Port: config.SSHPort,
		Password: config.SSHPassword,
		PrivateKeyPath: config.SSHKeyPath,
	}

	sshConnection := func(dstAddr, srcAddr string) (net.Conn, error) {
		return sshTunnel.Forward(dstAddr, srcAddr)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", config.Port))
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
