package selfsigned

import "crypto/tls"

func GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {}
