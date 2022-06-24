package proxyhandler

import (
	"github.com/mailsac/smtpd"
	"net"
	"testing"
)

func Test_ProxyHandlerV1(t *testing.T) {
	t.Run("ehlo", func(t *testing.T) {
		proxyIPs := []string{"10.0.0.1", "10.0.0.2"}
		p := ProxyHandlerV1{TrustIPs: proxyIPs}
		p.EHLO()
	})
	t.Run("throws error when proxy not supported", func(t *testing.T) {
		proxyIPs := []string{"10.0.0.5"}
		p := ProxyHandlerV1{TrustIPs: proxyIPs}
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		// Do some stuff
		conn := &smtpd.Conn{
			Conn: client,
		}
		err := p.Handle(conn, "TCP4 10.0.0.99 2.2.2.2 33372 25")
		if err == nil || err.Error() != "PROXY not allowed from 'pipe'" {
			t.Errorf("Should not have allowed proxy: %s", err.Error())
			t.Fail()
		}
	})
	t.Run("sets and trusts IP when proxied", func(t *testing.T) {
		proxyIPs := []string{"pipe"} // would be IP address in real life
		p := ProxyHandlerV1{TrustIPs: proxyIPs}
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		// Do some stuff
		conn := &smtpd.Conn{
			Conn: client,
		}
		err := p.Handle(conn, "TCP4 45.76.28.175 10.0.0.12 33372 25")
		if err != nil {
			t.Error(err)
			return
		}
		if conn.ForwardedForIP != "45.76.28.175" {
			t.Errorf("forward addr not set as expected (%s)", conn.ForwardedForIP)
		}
	})
}
