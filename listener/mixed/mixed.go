package mixed

import (
	"net"

	"github.com/kitty314/1.18.9/adapter/inbound"
	N "github.com/kitty314/1.18.9/common/net"
	"github.com/kitty314/1.18.9/component/auth"
	C "github.com/kitty314/1.18.9/constant"
	authStore "github.com/kitty314/1.18.9/listener/auth"
	"github.com/kitty314/1.18.9/listener/http"
	"github.com/kitty314/1.18.9/listener/socks"
	"github.com/kitty314/1.18.9/transport/socks4"
	"github.com/kitty314/1.18.9/transport/socks5"
)

type Listener struct {
	listener net.Listener
	addr     string
	closed   bool
}

// RawAddress implements C.Listener
func (l *Listener) RawAddress() string {
	return l.addr
}

// Address implements C.Listener
func (l *Listener) Address() string {
	return l.listener.Addr().String()
}

// Close implements C.Listener
func (l *Listener) Close() error {
	l.closed = true
	return l.listener.Close()
}

func New(addr string, tunnel C.Tunnel, additions ...inbound.Addition) (*Listener, error) {
	return NewWithAuthenticator(addr, tunnel, authStore.Default, additions...)
}

func NewWithAuthenticator(addr string, tunnel C.Tunnel, store auth.AuthStore, additions ...inbound.Addition) (*Listener, error) {
	isDefault := false
	if len(additions) == 0 {
		isDefault = true
		additions = []inbound.Addition{
			inbound.WithInName("DEFAULT-MIXED"),
			inbound.WithSpecialRules(""),
		}
	}

	l, err := inbound.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	ml := &Listener{
		listener: l,
		addr:     addr,
	}
	go func() {
		for {
			c, err := ml.listener.Accept()
			if err != nil {
				if ml.closed {
					break
				}
				continue
			}
			store := store
			if isDefault || store == authStore.Default { // only apply on default listener
				if !inbound.IsRemoteAddrDisAllowed(c.RemoteAddr()) {
					_ = c.Close()
					continue
				}
				if inbound.SkipAuthRemoteAddr(c.RemoteAddr()) {
					store = authStore.Nil
				}
			}
			go handleConn(c, tunnel, store, additions...)
		}
	}()

	return ml, nil
}

func handleConn(conn net.Conn, tunnel C.Tunnel, store auth.AuthStore, additions ...inbound.Addition) {
	bufConn := N.NewBufferedConn(conn)
	head, err := bufConn.Peek(1)
	if err != nil {
		return
	}

	switch head[0] {
	case socks4.Version:
		socks.HandleSocks4(bufConn, tunnel, store, additions...)
	case socks5.Version:
		socks.HandleSocks5(bufConn, tunnel, store, additions...)
	default:
		http.HandleConn(bufConn, tunnel, store, additions...)
	}
}
