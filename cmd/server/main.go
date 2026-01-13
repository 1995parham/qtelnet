package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/quic-go/quic-go"
)

func main() {
	listener, err := quic.ListenAddr("0.0.0.0:8080", generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("QUIC echo server listening on :8080")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("accept connection error: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn *quic.Conn) {
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("accept stream error: %v", err)
			return
		}

		go handleStream(conn, stream)
	}
}

func handleStream(conn *quic.Conn, stream *quic.Stream) {
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		log.Printf("read error: %v", err)
		return
	}

	fmt.Printf("Received: %s", string(data))

	// Echo back via a new stream
	respStream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Printf("open stream error: %v", err)
		return
	}
	defer respStream.Close()

	response := fmt.Sprintf("Echo: %s", string(data))
	respStream.Write([]byte(response))
}

// generateTLSConfig creates a self-signed certificate for testing
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		DNSNames:     []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		log.Fatal(err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo"},
	}
}
