package detector

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const ExceedTime = 5

type Ping struct {
}

func NewPing() *Ping {
	return &Ping{}
}

type PingResult struct {
	HasPermission bool
	PeerIP        string
	Pass          bool
}

func (*Ping) TryPing(ip string) (*PingResult, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return &PingResult{HasPermission: false}, nil
	}
	defer func() {
		_ = conn.Close()
	}()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return nil, fmt.Errorf("marshal message error \"%w\", please show it to your developers", err)
	}

	_ = conn.SetWriteDeadline(time.Now().Add(ExceedTime * time.Second))
	if _, err := conn.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(ip)}); err != nil {
		return nil, fmt.Errorf("ping ip %s 失败，因为 \"%w\", 请询问IT协助（可能是防火墙原因）", ip, err)
	}

	rb := make([]byte, 1500)
	_ = conn.SetReadDeadline(time.Now().Add(ExceedTime * time.Second))
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return nil, fmt.Errorf("ping ip %s 失败，因为 \"%w\", 请询问IT协助（可能是防火墙原因）", ip, err)
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
	if err != nil {
		return nil, fmt.Errorf("逻辑错误 \"%w\", 请询问研发协助", err)
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return &PingResult{HasPermission: true, PeerIP: peer.String(), Pass: peer.String() == ip}, nil
	case ipv4.ICMPTypeTimeExceeded:
		return nil, fmt.Errorf("请求 %s 超时", ip)
	case ipv4.ICMPTypeDestinationUnreachable:
		return nil, fmt.Errorf("目标地址 %s 不可达", ip)
	default:
		return nil, fmt.Errorf("不期望得到的ping回包相应类型 %d", rm.Type)
	}
}
