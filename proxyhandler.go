package proxyhandler

import (
	"errors"
	"fmt"
	"github.com/mailsac/smtpd"
	"strings"
)

// ProxyHandlerV1 pre checks IPs to see if they are from haproxy
type ProxyHandlerV1 struct {
	TrustIPs []string
}

// Handle implements the expected method for a smtp handler
func (p *ProxyHandlerV1) Handle(conn *smtpd.Conn, methodBody string) error {
	remoteIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
	if !sliceContains(p.TrustIPs, remoteIP) {
		return errors.New("PROXY not allowed from '" + remoteIP + "'")
	}

	phead, err := newProxyHeaderV1(methodBody)
	if err != nil {
		return err
	}

	isHealthCheck := sliceContains(p.TrustIPs, phead.EndUserIP)
	if isHealthCheck {
		return nil
	}

	conn.ForwardedForIP = phead.EndUserIP
	return nil
}

// EHLO also exports expected behavior
func (p *ProxyHandlerV1) EHLO() string {
	return "PROXY"
}

func newProxyHeaderV1(methodBody string) (*ProxyHeaderV1, error) {
	// methodBody: "TCP4 209.85.214.42 45.76.28.175 33372 25"
	//				0	 1			   2			3     4
	// 					 src	  	   dest         src   dest
	methodBodyParts := strings.Split(methodBody, " ")
	if len(methodBodyParts) < 5 {
		return nil, errors.New(fmt.Sprintf("PROXY v1 format is invalid, %s", methodBody))
	}
	return &ProxyHeaderV1{
		ProtoName:   methodBodyParts[0],
		EndUserIP:   methodBodyParts[1],
		ProxyIP:     methodBodyParts[2],
		EndUserPort: methodBodyParts[3],
		ProxyPort:   methodBodyParts[4],
	}, nil
}

type ProxyHeaderV1 struct {
	ProtoName   string
	EndUserIP   string
	EndUserPort string
	ProxyIP     string
	ProxyPort   string
}

func sliceContains(slice []string, comparator string) bool {
	for _, s := range slice {
		if s == comparator {
			return true
		}
	}
	return false
}
