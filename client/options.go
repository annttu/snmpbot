package client

import (
	"flag"
	"github.com/qmsk/snmpbot/snmp"
	"time"
)

const (
	SNMPVersion    = snmp.SNMPv2c
	DefaultTimeout = 1 * time.Second
	DefaultRetry   = uint(3)
	DefaultMaxVars = uint(10)
)

type Options struct {
	Community string
	Timeout   time.Duration
	Retry     uint
	UDP       UDPOptions
	MaxVars   uint
}

func (options *Options) InitFlags() {
	flag.StringVar(&options.Community, "snmp-community", "public", "Default SNMP community")
	flag.DurationVar(&options.Timeout, "snmp-timeout", DefaultTimeout, "SNMP request timeout")
	flag.UintVar(&options.Retry, "snmp-retry", 0, "SNMP request retry")
	flag.UintVar(&options.UDP.Size, "snmp-udp-size", UDPSize, "Maximum UDP recv size")
	flag.UintVar(&options.MaxVars, "snmp-maxvars", DefaultMaxVars, "Maximum request VarBinds")
}
