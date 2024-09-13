package webserver

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

// creates a tls server certificate by using a ca and its key to sign this certificate
func createCertificate(caPath string, keyPath string, hosts []string) (tls.Certificate, error) {
	// read ca and key file

	caFile, err := os.ReadFile(caPath)
	if err != nil {
		return tls.Certificate{}, err
	}

	caKeyFile, err := os.ReadFile(keyPath)
	if err != nil {
		return tls.Certificate{}, err
	}

	// decode pem

	caBlock, _ := pem.Decode(caFile)
	if caBlock == nil || caBlock.Type != "CERTIFICATE" {
		return tls.Certificate{}, errors.New("failed to decode PEM block containing certificate")
	}

	caKeyBlock, _ := pem.Decode(caKeyFile)
	if caKeyBlock == nil || caKeyBlock.Type != "PRIVATE KEY" {
		return tls.Certificate{}, errors.New("failed to decode PEM block containing private key")
	}

	// parse ca and key

	caTemplate, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return tls.Certificate{}, err
	}

	caKey, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return tls.Certificate{}, err
	}

	// create new certificate essential stuff

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}

	subject := pkix.Name{
		// if this information is missing, the certificate may not be trusted:
		CommonName:         "liquipay.de",                                // required by openssl
		Organization:       []string{"Liquipay UG (haftungsbeschränkt)"}, // required by openssl
		OrganizationalUnit: []string{"IT"},                               // required by openssl
		Country:            []string{"DE"},                               // required by openssl
		Province:           []string{"Nordrhein-Westfalen"},              // required by openssl
		Locality:           []string{"Lindlar"},                          // required by openssl
		PostalCode:         []string{"51789"},                            // optional
		StreetAddress:      []string{"Hauptstraße 10"},                   // optional
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(60 * 24 * time.Hour) // 60 days is recommended for certificates

	// create new certificate template

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    notBefore,
		NotAfter:     notAfter,

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// assign hosts to certificate

	for _, host := range hosts {
		ip := net.ParseIP(host)
		if ip != nil {
			certTemplate.IPAddresses = append(certTemplate.IPAddresses, ip)
		} else {
			certTemplate.DNSNames = append(certTemplate.DNSNames, host)
		}
	}

	// create certificate and key bytes

	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, caTemplate, &priv.PublicKey, caKey)
	if err != nil {
		return tls.Certificate{}, nil
	}

	// PKCS#8 is a standard for storing private key information for any algorithm
	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, nil
	}

	// create pem blocks

	certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})

	// create tls certificate and return

	return tls.X509KeyPair(certPEMBlock, keyPEMBlock)
}

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

func CreateWebServerAndCertificate(addr string, caPath string, keyPath string, hosts []string, handler http.Handler) error {
	cert, err := createCertificate(caPath, keyPath, hosts)
	if err != nil {
		return err
	}

	return listenAndServeTLS(addr, cert, handler)
}

func CreateWebServer(addr string, certPath string, keyPath string, handler http.Handler) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}

	return listenAndServeTLS(addr, cert, handler)
}
