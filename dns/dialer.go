package dns

// export functions from tunnel module

import "github.com/kitty314/1.18.9/tunnel"

const RespectRules = tunnel.DnsRespectRules

type dnsDialer = tunnel.DNSDialer

var newDNSDialer = tunnel.NewDNSDialer
