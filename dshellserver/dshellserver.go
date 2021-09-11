package main

import (
	"crypto/rand"
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

//go:embed server.pem
var serverPem []byte

//go:embed server.key
var serverKey []byte

func main() {
	portFlag := flag.Int("p", 9001, "Listen port")

	flag.Parse()

	cert, err := tls.X509KeyPair(serverPem, serverKey)
	if err != nil {
		log.Fatalf("Could not load keys: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader

	log.Printf("Waiting for connection on port %d", *portFlag)

	listener, err := tls.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *portFlag), &config)
	if err != nil {
		log.Fatalf("Could not listen: %s", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Error accepting client - %s", err)
	}

	defer conn.Close()

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalf("Error creating raw terminal: %s", err)
	}

	defer terminal.Restore(0, oldState)

	log.Printf("Connection recieved from %s", conn.RemoteAddr())

	go func() {
		io.Copy(conn, os.Stdin)
	}()
	io.Copy(os.Stdout, conn)
}
