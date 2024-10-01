package internal

import (
	"crypto/tls"
	"os"
)

func GetCertificates() (tls.Certificate, error) {
	bytesCert, err := os.ReadFile(certFile)
	if err != nil {
		return tls.Certificate{}, err
	}
	bytesKey, err := os.ReadFile(keyFile)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.X509KeyPair(bytesCert, bytesKey)
}
