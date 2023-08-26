package timer

import (
	"fmt"
	"network-detector/libs/logger"
	"time"
)

type Timer struct {
	ticker      *time.Ticker
	controlChan chan bool
	running     bool
}

func NewTimer() *Timer {
	return &Timer{
		ticker:      nil,
		controlChan: make(chan bool),
		running:     false,
	}
}

func (t *Timer) Start(domain string) {
	if t.running {
		fmt.Println("A timer is running.")
		return
	}

	t.ticker = time.NewTicker(30 * time.Second)

	go func() {
		t.running = true
		for {
			select {
			case <-t.ticker.C:
				println("心跳.")
				logger.Info("心跳.")
				task := NewTask()
				task.Run(domain)
			case <-t.controlChan:
				t.ticker.Stop()
				t.running = false
				logger.Info("Timer stopped.")
				return
			}
		}
	}()
}

func (t *Timer) Stop() {
	if !t.running {
		logger.Info("Timer is not running.")
		return
	}

	t.controlChan <- true
}
