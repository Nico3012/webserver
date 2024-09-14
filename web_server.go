package webserver

import (
	"crypto/tls"
	"encoding/pem"
	"net/http"
)

func listenAndServeTLS(addr string, cert tls.Certificate, handler http.Handler) error {
	/*netListener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(netListener, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})*/

	listener, err := tls.Listen("tcp", addr, &tls.Config{
		Certificates: []tls.Certificate{cert},
		/*CurvePreferences: []tls.CurveID{
			tls.X25519,
		},*/

		// tls 1.3 does not allow to pick a specific cipher or prefer e.g. aes over chacha
		// the used cipher is picked by the server and client automatically.
		// often aes-gcm ist prefered over chacha if both (server and client) support AES hardware acceleration.
		// e.g. raspberry PI 4 does not support AES hardware acceleration and therefore also prefers chacha (as client or as server)
	})
	if err != nil {
		return err
	}

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	tlsConn, ok := conn.(*tls.Conn)
	// 	if !ok {
	// 		return errors.New("failed to cast to TLS connection")
	// 	}

	// 	err = tlsConn.Handshake()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	buf := make([]byte, 1024)
	// 	for {
	// 		n, err := tlsConn.Read(buf)
	// 		if err != nil {
	// 			if err != io.EOF {
	// 				return err
	// 			}
	// 			break
	// 		}

	// 		log.Printf("Received: %s", string(buf[:n]))

	// 		_, err = tlsConn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!"))
	// 		if err != nil {
	// 			break
	// 		}
	// 	}
	// }

	// http module
	// fmt.Println("starting web server...")
	return http.Serve(listener, handler)
}

func CreateWebServerAndCertificate(addr string, caPath string, caKeyPath string, hosts []string, handler http.Handler) error {
	certBytes, certKeyBytes, err := generateCertBytes(caPath, caKeyPath, hosts)
	if err != nil {
		return err
	}

	// create pem blocks
	certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certKeyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: certKeyBytes})

	// create tls certificate and return
	cert, err := tls.X509KeyPair(certPEMBlock, certKeyPEMBlock)

	return listenAndServeTLS(addr, cert, handler)
}

func CreateWebServer(addr string, certPath string, certKeyPath string, handler http.Handler) error {
	cert, err := tls.LoadX509KeyPair(certPath, certKeyPath)
	if err != nil {
		return err
	}

	return listenAndServeTLS(addr, cert, handler)
}
