# smtpd-proxyhandler
HAProxy SMTP proxy-protocol handler for extending the golang [mailsac mailproto/smtpd fork](https://github.com/mailsac/smtpd) library.

https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt

v1 protocol is supported.
v2 protocol is not supported (TBD).

This package allows an upstream proxy to modify the end user IP address on an SMTP connection.

## Usage

```go
import (
    proxyhandler "github.com/mailsac/smtpd-proxyhandler"
    "github.com/mailsac/smtpd"
)

server := smtpd.NewServer(messageHandler)

// ... set up server, then add upstream proxy IP addresses
allowProxyIPs := []string{"10.0.0.1", "10.0.0.2"}
handler := proxyhandler.ProxyHandlerV1{ TrustIPs: allowProxyIPs }
server.Extend("PROXY", &handler)
```
