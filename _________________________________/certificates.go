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
	"os"
	"time"
)

func generateCaBytes() ([]byte, []byte, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return []byte{}, []byte{}, err
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
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    notBefore,
		NotAfter:     notAfter,

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// create ca and key

	// create certificate bytes
	caBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// PKCS#8 is a standard for storing private key information for any algorithm
	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return caBytes, keyBytes, nil
}

// creates a tls server certificate by using a ca and its key to sign this certificate
func generateCertBytes(caPath string, caKeyPath string, hosts []string) ([]byte, []byte, error) {
	// read ca and key file

	caFile, err := os.ReadFile(caPath)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	caKeyFile, err := os.ReadFile(caKeyPath)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// decode pem

	caBlock, _ := pem.Decode(caFile)
	if caBlock == nil || caBlock.Type != "CERTIFICATE" {
		return []byte{}, []byte{}, errors.New("failed to decode PEM block containing certificate")
	}

	caKeyBlock, _ := pem.Decode(caKeyFile)
	if caKeyBlock == nil || caKeyBlock.Type != "PRIVATE KEY" {
		return []byte{}, []byte{}, errors.New("failed to decode PEM block containing private key")
	}

	// parse ca and key

	caTemplate, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	caKey, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// create new certificate essential stuff

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return []byte{}, []byte{}, err
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
		return []byte{}, []byte{}, err
	}

	// PKCS#8 is a standard for storing private key information for any algorithm
	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return certBytes, keyBytes, nil
}

// can be used for both cert and ca bytes!
func createCertFiles(certBytes []byte, certKeyBytes []byte, certWritePath string, certKeyWritePath string) error {
	// create cert file
	certOut, err := os.Create(certWritePath)
	if err != nil {
		return err
	}
	// write to cert file
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		return err
	}

	// create key file
	keyOut, err := os.Create(certKeyWritePath)
	if err != nil {
		return err
	}
	// write to key file
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: certKeyBytes})
	if err != nil {
		return err
	}

	return nil
}

func GenerateCAFiles(caWritePath string, caKeyWritePath string) error {
	caBytes, caKeyBytes, err := generateCaBytes()
	if err != nil {
		return err
	}

	err = createCertFiles(caBytes, caKeyBytes, caWritePath, caKeyWritePath)
	if err != nil {
		return err
	}

	return nil
}

func GenerateCertFiles(certWritePath string, certKeyWritePath string, caPath string, caKeyPath string, hosts []string) error {
	certBytes, certKeyBytes, err := generateCertBytes(caPath, caKeyPath, hosts)
	if err != nil {
		return err
	}

	err = createCertFiles(certBytes, certKeyBytes, certWritePath, certKeyWritePath)
	if err != nil {
		return err
	}

	return nil
}

// get certificate interface

func GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {}
