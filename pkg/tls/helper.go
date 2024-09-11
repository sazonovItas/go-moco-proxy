package tls

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// getCertPool function returns new cert pool from given cert path.
func getCertPool(certPath string) (*x509.CertPool, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		return nil, fmt.Errorf("failed to parse cert from %s", certPath)
	}

	return certPool, nil
}

// getCaCert function returns new certificate from given cert path.
func getCaCert(certPath string) (*x509.Certificate, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	b, _ := pem.Decode(cert)
	if b == nil {
		return nil, fmt.Errorf("failed to parse ca cert from %s", certPath)
	}

	return x509.ParseCertificate(b.Bytes)
}
