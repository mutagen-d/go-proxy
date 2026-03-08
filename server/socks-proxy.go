package server

import "net"

func HandleSocks5(conn net.Conn, data []byte, createConnection CreateConnection) {
	if createConnection == nil {
		createConnection = directConnection
	}
	// TODO
}

func HandleSocks4(conn net.Conn, data []byte, createConnection CreateConnection) {
	if createConnection == nil {
		createConnection = directConnection
	}
	// TODO
}