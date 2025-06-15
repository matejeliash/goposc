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
}

func NewPinger() *Pinger {
	return &Pinger{
		ConcurrencyLimit: make(chan struct{}, 100),
		Results:          make(chan string, 256),
	}
}

func (p *Pinger) PingIP(ip string) {

	defer p.Wg.Done()
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 1

	err = pinger.Run()
	if err != nil {
		return
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		p.Results <- ip

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

	// fmt.Println("Devices found:")
	// for ip := range p.Results {
	// 	fmt.Println(ip)
	// }
	// Collect found IPs
	var foundIPs []string
	for ip := range p.Results {
		fmt.Println("Device found:", ip)
		foundIPs = append(foundIPs, ip)
	}

	return foundIPs

}

// func PingIpRange(ipSlice []string) {

// 	//subnet := "192.168.0." // change to your subnet prefix
// 	var wg sync.WaitGroup
// 	results := make(chan string, 256)

// 	concurrencyLimit := make(chan struct{}, 100) // limit concurrent pings to 50

// 	for _, ip := range ipSlice {
// 		//ip := fmt.Sprintf("%s%d", subnet, i)
// 		wg.Add(1)
// 		concurrencyLimit <- struct{}{} // acquire a token
// 		go func(ip string) {
// 			defer func() { <-concurrencyLimit }() // release token
// 			pingIP(ip, &wg, results)
// 		}(ip)
// 	}

// 	// Close results channel after all pings finish
// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	fmt.Println("Devices found:")
// 	for ip := range results {
// 		fmt.Println(ip)
// 	}
// }

// func pingIP(ip string, wg *sync.WaitGroup, results chan<- string) {
// 	defer wg.Done()
// 	pinger, err := probing.NewPinger(ip)
// 	if err != nil {
// 		return
// 	}
// 	pinger.Count = 1
// 	pinger.Timeout = time.Second * 1

// 	err = pinger.Run()
// 	if err != nil {
// 		return
// 	}
// 	stats := pinger.Statistics()
// 	if stats.PacketsRecv > 0 {
// 		results <- ip
// 	}
// }
