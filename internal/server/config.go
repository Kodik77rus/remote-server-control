package server

import (
	"crypto/tls"
	"log"
	"path/filepath"
)

const (
	_certificate = "certificate.pem"
	_publicKey   = "key.pem"
	_port        = ":8443"
)

//Server config
type Config struct {
	port   string
	tlsCfg *tls.Config
}

//Server config constructor
func NewConfig(certPath string) *Config {
	cnfg, err := loadCertificate(certPath)
	checkFatalError(err)

	return &Config{
		port:   _port,
		tlsCfg: cnfg,
	}
}

//load Ssl certificate
func loadCertificate(certPaht string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(
		filepath.Join(certPaht, _certificate),
		filepath.Join(certPaht, _publicKey),
	)

	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func checkFatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
