package detector

import (
	"fmt"
	"net"
)

type DNS struct {
}

func NewDNS() *DNS {
	return &DNS{}
}

func (*DNS) ResolveHost(host string) (string, error) {
	ip, err := net.ResolveIPAddr("", host)
	if err != nil {
		return "", fmt.Errorf("解析域名'%s'失败，因为 '%s', 请联系IT做检查（可能是ISP, Hardware的问题）", host, err)
	}
	return ip.String(), nil
}
