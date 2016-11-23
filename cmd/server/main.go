package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {
	addrs := "0.0.0.0:6666"
	listener, error := net.Listen("tcp", addrs)
	if error != nil {
		log.Fatal(error)
	}
	defer listener.Close()
	fmt.Printf("Listening on %v for incoming connecionts\n", addrs)
	for {
		conn, error := listener.Accept()
		if error != nil {
			log.Printf("Failed to get a network connection %v\n", error)
			continue
		}
		log.Printf("Established Connection from %v\n", conn.RemoteAddr())
		buffer := bytes.Buffer{}
		read, error := buffer.ReadFrom(conn)
		if error != nil {
			log.Printf("Failed to read from connection %v\n", error)
			continue
		}
		log.Printf("Received %v bytes from connection: %v", read, buffer.String())
		conn.Close()
	}
}
