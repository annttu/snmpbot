package client

import (
	"net/url"
)

type Config struct {
	Logging   Logging
	Community string `json:"community"`
	Addr      string `json:"addr"` // host or host:port
	OID       string `json:"oid"`
}

// Parse a pseudo-URL config string:
//  [community "@"] Host
func (config *Config) Parse(str string) error {
	str = "snmp://" + str

	if parseURL, err := url.Parse(str); err != nil {
		return err
	} else {
		return config.ParseURL(parseURL)
	}
}

func (config *Config) ParseURL(configURL *url.URL) error {
	if configURL.User != nil {
		config.Community = configURL.User.Username()
	}

	//log.Printf("ParseConfig %s: url=%#v\n", str, configUrl)
	config.Addr = configURL.Host

	if configURL.Path != "" {
		config.OID = configURL.Path[1:]
	} else {
		config.OID = ""
	}

	return nil
}

func (config Config) String() string {
	str := ""

	if config.Community != "" {
		str += config.Community + "@"
	}

	str += config.Addr

	if config.OID != "" {
		str += "/" + config.OID
	}

	return str
}
