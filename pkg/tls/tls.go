package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/sazonovItas/go-moco-proxy/pkg/config"
)

// newTLS function returns tls config with defaults settings.
func newTLS() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
}

// NewTLSConfig function returns tls config from given config.
func NewTLSConfig(c config.TLSConfig, isServer bool) (*tls.Config, error) {
	var (
		caPool   *x509.CertPool
		pairCert tls.Certificate
		err      error
	)

	tlsCfg := newTLS()
	if c.SNI != "" {
		tlsCfg.ServerName = c.SNI
	}

	if c.CaCert != "" {
		caPool, err = getCertPool(c.CaCert)
		if err != nil {
			return nil, fmt.Errorf("failed to load cert pool: %w", err)
		}

		tlsCfg.RootCAs = caPool
	}

	if c.Cert != "" && c.Key != "" {
		pairCert, err = tls.LoadX509KeyPair(c.Cert, c.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to load cert: %s", err)
		}

		tlsCfg.Certificates = []tls.Certificate{pairCert}
	}

	if c.IsMutual() {
		if isServer {
			tlsCfg.GetConfigForClient = func(_ *tls.ClientHelloInfo) (*tls.Config, error) {
				return &tls.Config{
					RootCAs:      caPool,
					Certificates: []tls.Certificate{pairCert},
					ClientAuth:   tls.RequireAndVerifyClientCert,
					MinVersion:   tls.VersionTLS13,
				}, nil
			}
		} else {
			tlsCfg.VerifyConnection = func(cs tls.ConnectionState) error {
				return nil
			}
		}
	}

	return tlsCfg, nil
}
