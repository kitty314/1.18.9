package inbound

import (
	"net"
	"net/http"

	C "github.com/kitty314/1.18.9/constant"
)

// NewHTTPS receive CONNECT request and return ConnContext
func NewHTTPS(request *http.Request, conn net.Conn, additions ...Addition) (net.Conn, *C.Metadata) {
	metadata := parseHTTPAddr(request)
	metadata.Type = C.HTTPS
	metadata.RawSrcAddr = conn.RemoteAddr()
	metadata.RawDstAddr = conn.LocalAddr()
	ApplyAdditions(metadata, WithSrcAddr(conn.RemoteAddr()), WithInAddr(conn.LocalAddr()))
	ApplyAdditions(metadata, additions...)
	return conn, metadata
}
