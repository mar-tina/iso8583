package main

import (
	"io"
	"log"
	"net"
)

func Listen(port, url string) {
	netListener, err := net.Listen("tcp", url+":"+port)
	if err != nil {
		log.Printf("connection failed: %s", err)
	}

	for {
		conn, err := netListener.Accept()
		if err != nil {
			log.Printf("conn accept failed: %s", err)
		}

		go func(c net.Conn) {
			reader := make([]byte, 1024)
			nBytes, err := c.Read(reader)
			if err != nil {
				if err != io.EOF {
					log.Printf("conn issue: %s", err)
				}
			}

			c.Write([]byte("0810823A000002000000048000000000000004200906139000010906130420042000001031128"))

			log.Printf("read bytes %s", reader[:nBytes])
		}(conn)
	}
}

func main() {
	log.Printf("running")
	Listen("9999", "127.0.0.1")
}
