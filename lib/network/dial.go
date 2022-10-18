package network

import (
	"io"
	"log"
	"net"
)

type ISODialer struct {
	conn net.Conn
}

func NewDialer(host, uri string) (*ISODialer, bool) {
	conn, err := net.Dial("tcp", host+":"+uri)
	if err != nil {
		log.Fatalf("conn dial failed: %s", err)
	}

	return &ISODialer{
		conn: conn,
	}, true
}

func (dialer *ISODialer) Write(payload []byte) error {
	n, err := dialer.conn.Write(payload)
	if err != nil {
		log.Printf("dialer write failed: %s", err)
		return err
	}

	log.Printf("written %d bytes", n)
	return nil
}

func (dialer *ISODialer) Read() {
	readBuffer := make([]byte, 1024)
	for {
		n, err := dialer.conn.Read(readBuffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("conn issue: %s", err)
			}
		}
		log.Printf("read: %s", readBuffer[:n])
	}
}
