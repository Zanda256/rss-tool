package es

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
)

type EsClient struct {
	*elasticsearch8.Client
}

type Config struct {
	URL      string
	User     string
	Password string
}

func NewESClient(config Config) (*EsClient, error) {
	esCfg := elasticsearch8.Config{
		Addresses: []string{
			config.URL,
		},
		Username: config.User,
		Password: config.User,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	c, err := elasticsearch8.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	return &EsClient{c}, nil
}
