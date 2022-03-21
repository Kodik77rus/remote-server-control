package server

import (
	"crypto/tls"
	"log"
	"os"
	"path"
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
func NewConfig() *Config {
	cnfg, err := loadCertificate()
	checkFatalError(err)

	return &Config{
		port:   _port,
		tlsCfg: cnfg,
	}
}

//load Ssl certificate
func loadCertificate() (*tls.Config, error) {
	dirname, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(
		path.Join(dirname, "../../configs/certificate/", _certificate),
		path.Join(dirname, "../../configs/certificate/", _publicKey),
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
