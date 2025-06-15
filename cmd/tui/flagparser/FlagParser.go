package flagparser

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/matejeliash/goposc/internal/netinfo"
	"github.com/matejeliash/goposc/internal/pinger"
	"github.com/matejeliash/goposc/internal/portscanner"
)

type FlagParser struct {
	IpRange  []string
	Ports    []int
	FoundIps []string
}

func (fp *FlagParser) Parse() error {
	ipRangeFlag := flag.String("ip-range", "", "IP range in start-end format (e.g. 192.168.1.10-192.168.1.20)")
	ipFlag := flag.String("ip", "", "IP adress (e.g. 192.168.1.10)")
	ipsFlag := flag.String("ips", "", "IP adresses (e.g. 192.168.1.10,192.168.0.111)")
	pingFlag := flag.Bool("ping", false, "use ping to find if devices with specified id apdresses are active")
	portFlag := flag.Int("port", 0, "specified port")
	portsFlag := flag.String("ports", "", "specified ports")
	portScanFlag := flag.Bool("port-scan", false, "use port scan to fin if service is active")
	getBannerFlag := flag.Bool("get-banner", false, "use port scan to fin if service is active")
	netInfoFlag := flag.Bool("net-info", false, "get basic network information")
	flag.Parse()

	// if *ipRangeFlag == "" && *ipFlag == "" && *ipsFlag == "" {
	// 	return fmt.Errorf("you must specify ip or ip-range or ips")
	// }

	// if *ipRangeFlag != "" && *ipFlag != "" {
	// 	return fmt.Errorf("you must specify either ip or ip-range")
	// }

	if *ipFlag != "" {
		if *ipFlag == "localhost" || isIPv4(*ipFlag) || isValidDomain(*ipFlag) {
			fp.IpRange = append(fp.IpRange, *ipFlag)
		} else {
			return fmt.Errorf("specified ip %s address is incorrect", *ipFlag)

		}

	}
	if *ipRangeFlag != "" {
		err := fp.ParseIpRange(*ipRangeFlag)
		if err != nil {
			return err
		}

	}

	if *ipsFlag != "" {
		splittedIps := strings.Split(*ipsFlag, ",")
		for _, ip := range splittedIps {
			if ip == "localhost" {
				fp.IpRange = append(fp.IpRange, ip)
			} else if isIPv4(ip) {
				fp.IpRange = append(fp.IpRange, ip)
			} else {
				return fmt.Errorf("specified ip %s address is incorrect", ip)

			}
		}
	}

	if *portFlag > 0 {
		fp.Ports = append(fp.Ports, *portFlag)
		//fmt.Println(fp.Ports)

	}

	if *portsFlag != "" {

		if *portsFlag == "all" {
			for p := 1; p < 65535; p++ {
				fp.Ports = append(fp.Ports, p)
			}

		} else {
			splittedPorts := strings.Split(*portsFlag, ",")
			for _, portStr := range splittedPorts {
				port, err := strconv.Atoi(portStr)
				if err != nil {
					return fmt.Errorf("you specified wront port number %s", portStr)
				}
				fp.Ports = append(fp.Ports, port)
			}

		}
	}

	if *portScanFlag {
		if len(fp.IpRange) == 0 {
			return fmt.Errorf("you must specify at least one ip adress")
		}

		if len(fp.Ports) == 0 {
			return fmt.Errorf("you must specify at least one port")
		}

		for _, ip := range fp.IpRange {
			portscanner.ScanPortsOfIP(ip, fp.Ports)

		}
	}

	if *getBannerFlag {
		if len(fp.IpRange) == 0 {
			return fmt.Errorf("you must specify at least one ip adress")
		}

		if len(fp.Ports) == 0 {
			return fmt.Errorf("you must specify at least one port")
		}

		portscanner.GetBanner(fp.IpRange[0], fp.Ports[0])

	}

	if *pingFlag {
		//pinger.PingIpRange(fp.IpRange)
		p := pinger.NewPinger()
		fi := p.PingAllIPs(fp.IpRange)
		for _, ip := range fi {
			fp.FoundIps = append(fp.FoundIps, ip)
		}
		//fmt.Println(fp.IpRange)
	}

	if *netInfoFlag {
		ni := &netinfo.NetworkInfo{}

		ni.GetInfos()
		ni.PrintInfo()
	}

	return nil

}

// ipToUint32 converts net.IP to uint32
func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// uint32ToIP converts uint32 to net.IP
func uint32ToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}

func (fp *FlagParser) ParseIpRange(ipRangeStr string) error {

	splittedIpRange := strings.Split(ipRangeStr, "-")
	if len(splittedIpRange) != 2 {
		return fmt.Errorf("invalid IP range")
	}
	startIP := net.ParseIP(splittedIpRange[0])
	endIP := net.ParseIP(splittedIpRange[1])

	if startIP == nil || endIP == nil {
		return fmt.Errorf("invalid IP range")
	}

	startIPInt := ipToUint32(startIP)
	endIPInt := ipToUint32(endIP)
	if startIPInt > endIPInt {
		return fmt.Errorf("invalid IP range")
	}

	for i := startIPInt; i <= endIPInt; i++ {
		fp.IpRange = append(fp.IpRange, uint32ToIP(i).String())
	}

	return nil
}

func isIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && ip.To4() != nil
}

// func (fp *FlagParser) ParseIps(ipParam string) {
// 	if ipParam == "localhost" {
// 		fp.IpRange = append(fp.IpRange, ipParam)
// 		return
// 	}

// 	if ipParam == "subnet" {
// 		return
// 	}

// 	if strings.Contains(ipParam, ",") {
// 		splitted := strings.Split(ipParam, ",")

// 	}

// }

func isValidDomain(domain string) bool {
	var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}
