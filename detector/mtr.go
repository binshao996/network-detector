package detector

import (
	"fmt"
	"runtime"
	"time"

	"github.com/tonobo/mtr/pkg/hop"
	"github.com/tonobo/mtr/pkg/mtr"
	"golang.org/x/net/icmp"
)

type MTR struct {
}

func NewMTR() *MTR {
	return &MTR{}
}

type MTRResult struct {
	HasPermission bool
	Result        []hop.HopStatistic
	MTRStatus     string
}

func (p *MTR) TryMTR(ip string) (*MTRResult, error) {
	// test permission
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return &MTRResult{HasPermission: false}, nil
	}
	_ = conn.Close()
	// try prehandle, like firewall
	if err := p.PreHandle(); err != nil {
		return nil, err
	}

	// try mtr
	m, ch, err := mtr.NewMTR(ip, "0.0.0.0", time.Second, 100*time.Millisecond, time.Millisecond, 20, 10, 50, false)
	if err != nil {
		return nil, fmt.Errorf("未知的错误'%w', 请将它展示给你的研发看", err)
	}
	go func() {
		m.Run(ch, 20)
		close(ch)
	}()
	// read until ch empty
	for range ch {
	}
	// 用与判断是否防火墙限制
	allFailed := true
	mtrStatus := "200"

	result := make([]hop.HopStatistic, len(m.Statistic))
	for i := 1; i <= len(m.Statistic); i++ {
		if m.Statistic[i].Lost != m.Statistic[i].Sent {
			allFailed = false
		}
		t := m.Statistic[i]
		result[i-1] = hop.HopStatistic{
			Dest:       t.Dest,
			Timeout:    t.Timeout,
			Sent:       t.Sent,
			TTL:        t.TTL,
			Targets:    t.Targets,
			Last:       t.Last,
			Best:       t.Best,
			Worst:      t.Worst,
			SumElapsed: t.SumElapsed,
			Lost:       t.Lost,
		}
	}

	if allFailed {
		mtrStatus = "500"
		// All Failed At Windows, Please check the computer's firewall
		// 到时候对方需要在windows高级防火墙的入站规则里面，打开虚拟机监控(回显请求-ICMPV4-In)
		if runtime.GOOS == "windows" {
			return nil, fmt.Errorf("所有mtr探测均失败, 可能是因为你电脑防火墙设置的原因, 请联系研发或者IT协助")
		}
	}

	// Construct the MTR result object
	mtrResult := &MTRResult{
		HasPermission: true,
		Result:        result,
		MTRStatus:     mtrStatus,
	}

	return mtrResult, nil
}
