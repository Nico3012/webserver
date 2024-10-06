package webserver

// exported functions:

// this function creates a web server from certificate files
// the server automatically reacts to certificate file updates
func CreateWebServer(certPath string, certKeyPath string) {}

// this function creates a web server from certificate authority (ca) files
// the server automatically reacts to certificate authority (ca) file updates
// the server automatically renews its certificates before they expire
// compared to Let’s Encrypt, this function does not need a cache location to store certificates, because self signed certificates can be generated as often as needed
func CreateWebServerWithSelfSignedCertificate(caPath string, caKeyPath string) {}

// this function creates a web server from Let’s Encrypt certificate authority (ca)
// the certificates will be renewed automatically
// testing boolean indicates if the staging environment is used
// compared to self signed, this function needs a cache location to store certificates, because Let’s Encrypt limits the amount of certificates, that can be generated within a few days
func CreateWebServerWithLetsEncryptCertificate(testing bool, cacheJsonPath string) {}
