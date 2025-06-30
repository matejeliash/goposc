package pinger

import (
	"fmt"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type Pinger struct {
	ConcurrencyLimit chan struct{}
	Results          chan string
	Wg               sync.WaitGroup
	TimeoutMs        int64
}

func NewPinger(timeoutMs int64) *Pinger {
	return &Pinger{
		ConcurrencyLimit: make(chan struct{}, 100), // limit to only 100 concurrent goroutines
		Results:          make(chan string, 256),   // channel storing only relevant used IPs
		TimeoutMs:        timeoutMs,
	}
}

func (p *Pinger) PingIP(ip string) {

	defer p.Wg.Done()
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return
	}
	pinger.Count = 1
	pinger.Timeout = time.Millisecond * time.Duration(p.TimeoutMs)

	err = pinger.Run()
	if err != nil {
		return
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		p.Results <- ip
		//fmt.Println("ip found: " + ip)

	}

}

func (p *Pinger) PingAllIPs(IPSlice []string) []string {

	for _, ip := range IPSlice {
		//ip := fmt.Sprintf("%s%d", subnet, i)
		p.Wg.Add(1)
		p.ConcurrencyLimit <- struct{}{} // acquire a token
		go func(ip string) {
			defer func() { <-p.ConcurrencyLimit }() // release token
			p.PingIP(ip)

		}(ip)
	}

	// Close results channel after all pings finish
	go func() {
		p.Wg.Wait()
		close(p.Results)
	}()

	// Collect found IPs
	var foundIPs []string
	for ip := range p.Results {
		fmt.Println("Device found with ip:", ip)
		foundIPs = append(foundIPs, ip)
	}

	return foundIPs

}
