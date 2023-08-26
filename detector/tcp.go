package detector

import (
	"fmt"
	"net"
	"net/netip"
)

type TCPConnection struct {
}

func NewTCPConnection() *TCPConnection {
	return &TCPConnection{}
}

func (t *TCPConnection) TryHTTPSConnection(ip string) error {
	// try prehandle, like firewall
	if err := t.PreHandle(); err != nil {
		return err
	}

	httpsServer := net.TCPAddrFromAddrPort(netip.AddrPortFrom(netip.MustParseAddr(ip), 443))
	conn, err := net.DialTCP("tcp", nil, httpsServer)
	if err != nil {
		return fmt.Errorf("对 %s:443 建立连接失败，因为 '%w', 请联系IT协助（可能是防火墙原因）", ip, err)
	}
	_ = conn.Close()
	return nil
}
