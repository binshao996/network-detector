package timer

import (
	"encoding/json"
	"fmt"
	"network-detector/detector"
	"network-detector/libs/logger"
	"sync"
	"time"
)

type Task struct {
}

// 上报数据结构
type ReportData struct {
	IP               string `json:"ip"`
	DnsResult        string `json:"dns_result"`
	PingResult       string `json:"ping_result"`
	TCPResult        string `json:"tcp_result"`
	MTRHopStatistics string `json:"mtr_hop_statistics"`
	MTRResult        string `json:"mtr_result"`
	Domain           string `json:"domain"`
	Ctime            int64  `json:"ctime"`
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) Run(domain string) {
	logger.Printf("Running task with domain: %s", domain)

	reportData := ReportData{
		IP:               "",
		PingResult:       "",
		TCPResult:        "",
		MTRHopStatistics: "",
		Domain:           domain,
		MTRResult:        "",
		Ctime:            time.Now().UnixNano() / int64(time.Millisecond),
	}

	// 同步执行DNS解析获取IP地址
	dns := detector.NewDNS()
	ip, err := dns.ResolveHost(domain)
	if err != nil {
		reportData.DnsResult = fmt.Sprintf("Failed to resolve host '%s': %s", domain, err.Error())
		logger.Errorf("Failed to resolve host '%s': %s\n", domain, err.Error())
	} else {
		reportData.IP = ip
		reportData.DnsResult = "200"

		var wg sync.WaitGroup
		// 异步执行Ping测试
		wg.Add(1)
		go func() {
			defer wg.Done()
			ping := detector.NewPing()
			_, err := ping.TryPing(ip)
			if err != nil {
				reportData.PingResult = fmt.Sprintf("Failed to ping ip '%s': %s", ip, err.Error())
			} else {
				reportData.PingResult = "200"
			}
		}()

		// 异步执行TCP 443连接测试
		wg.Add(1)
		go func() {
			defer wg.Done()
			tcp := detector.NewTCPConnection()
			err := tcp.TryHTTPSConnection(ip)
			if err != nil {
				reportData.TCPResult = fmt.Sprintf("TCP 443 error '%s': %s", ip, err.Error())
			} else {
				reportData.TCPResult = "200"
			}
		}()

		// 异步执行MTR诊断
		wg.Add(1)
		go func() {
			defer wg.Done()
			mtr := detector.NewMTR()
			mtrResult, err := mtr.TryMTR(ip)

			if err != nil {
				reportData.MTRResult = fmt.Sprintf("Failed to mtr ip '%s': %s", ip, err.Error())
			} else {
				jsonResult, err := json.Marshal(mtrResult.Result)
				if err != nil {
				}
				reportData.MTRHopStatistics = string(jsonResult)
				reportData.MTRResult = mtrResult.MTRStatus
			}
		}()

		// 等待所有goroutine完成
		wg.Wait()
	}

	// // 上报后台server，如果有诉求的话
	// result, err := httpRequest.Post(reportUrl, reportData, "application/json")
	// logger.Printf("report reulst %v\n", result)
	// if err != nil {
	// 	logger.Errorf("report error %s\n", err.Error())
	// }
}
