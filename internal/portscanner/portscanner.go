package portscanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type PortScanner struct {
	ConcurrencyLimit chan struct{}
	Results          chan int
	Wg               sync.WaitGroup
	TimeoutMs        int64
}

func NewPortScanner(timeoutMs int64) *PortScanner {

	return &PortScanner{
		ConcurrencyLimit: make(chan struct{}, 100),
		Results:          make(chan int, 1000),
		TimeoutMs:        timeoutMs,
	}
}

func (ps *PortScanner) ScanPort(ip string, port int) {
	defer ps.Wg.Done()
	address := fmt.Sprintf("%s:%d", ip, port)
	//fmt.Println(address)
	conn, err := net.DialTimeout("tcp", address, time.Duration(ps.TimeoutMs)*time.Millisecond)
	if err != nil {
		return // port is closed or filtered
	}
	conn.Close()
	ps.Results <- port

}

//	func (ps *PortScanner) ScanPort(ip string, port int) {
//		defer ps.Wg.Done()
//		address := fmt.Sprintf("%s:%d", ip, port)
//		fmt.Println("Scanning port:", port)
//		conn, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
//		if err != nil {
//			fmt.Println("Port closed or filtered:", port)
//			return
//		}
//		conn.Close()
//		fmt.Println("Port open found:", port)
//		ps.Results <- port
//	}
func (ps *PortScanner) ScanPortsOfIP(ip string, ports []int) []int {
	ps.Results = make(chan int, 1000)

	for _, port := range ports {
		ps.Wg.Add(1)
		ps.ConcurrencyLimit <- struct{}{} // acquire a token
		go func(p int) {
			defer func() { <-ps.ConcurrencyLimit }() // release the token
			ps.ScanPort(ip, p)
		}(port)
	}

	go func() {
		ps.Wg.Wait()
		close(ps.Results)
	}()

	// Collect found IPs
	var openedPorts []int
	for port := range ps.Results {
		fmt.Println("Opened port:", port)
		openedPorts = append(openedPorts, port)
	}

	return openedPorts
}

// func (p *Pinger) PingAllIPs(IPSlice []string) []string {

// 	for _, ip := range IPSlice {
// 		//ip := fmt.Sprintf("%s%d", subnet, i)
// 		p.Wg.Add(1)
// 		p.ConcurrencyLimit <- struct{}{} // acquire a token
// 		go func(ip string) {
// 			defer func() { <-p.ConcurrencyLimit }() // release token
// 			p.PingIP(ip)

// 		}(ip)
// 	}

// 	// Close results channel after all pings finish
// 	go func() {
// 		p.Wg.Wait()
// 		close(p.Results)
// 	}()

// 	// Collect found IPs
// 	var foundIPs []string
// 	for ip := range p.Results {
// 		fmt.Println("Device found with ip:", ip)
// 		foundIPs = append(foundIPs, ip)
// 	}

// 	return foundIPs

// }

// func ScanPortsOfIP(ip string, ports []int) {

// 	//fmt.Println(ports)
// 	var wg sync.WaitGroup
// 	results := make(chan int, 1000)

// 	// Limit concurrency to avoid overwhelming the system
// 	concurrency := 100
// 	sem := make(chan struct{}, concurrency)

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	for _, port := range ports {
// 		wg.Add(1)
// 		sem <- struct{}{} // acquire a token
// 		go func(p int) {
// 			defer func() { <-sem }() // release the token
// 			scanPort(ip, port, &wg, results)
// 		}(port)
// 	}

// 	fmt.Printf("Open ports on %s:\n", ip)
// 	for port := range results {
// 		fmt.Printf("    Port %d is open\n", port)
// 	}

// }

// func GetBanner(ip string, port int) {

// 	var defaultProbeStrings = map[int]string{
// 		21:   "USER anonymous\r\n",
// 		23:   "\r\n",
// 		25:   "HELO example.com\r\n",
// 		80:   "GET / HTTP/1.0\r\n\r\n",
// 		8080: "GET / HTTP/1.0\r\n\r\n",
// 		110:  "USER test\r\n",
// 		143:  "A1 CAPABILITY\r\n",
// 		443:  "GET / HTTP/1.0\r\n\r\n", // you should actually use TLS for 443
// 		3306: "",                       // MySQL sends banner automatically
// 	}

// 	address := fmt.Sprintf("%s:%d", ip, port)
// 	//fmt.Println(address)
// 	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
// 	if err != nil {
// 		return // port is closed or filtered
// 	}
// 	defer conn.Close()

// 	// Construct a minimal HTTP/1.0 GET request
// 	request := fmt.Sprintf(defaultProbeStrings[port], ip)

// 	// Send the HTTP request
// 	_, err = conn.Write([]byte(request))
// 	if err != nil {
// 		return
// 	}

// 	// Set a short read deadline so we don't hang forever
// 	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

// 	// Read up to 1024 bytes from the connection
// 	buffer := make([]byte, 1024)
// 	n, err := conn.Read(buffer)
// 	if err != nil || n <= 0 {
// 		return
// 	}

// 	fmt.Println(string(buffer))

// }

// func (p *Pinger) PingIP(ip string) {

// 	defer p.Wg.Done()
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
// 		p.Results <- ip

// 	}

// }
