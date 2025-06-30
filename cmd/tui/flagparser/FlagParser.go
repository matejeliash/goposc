package flagparser

import (
	"flag"
	"fmt"

	"github.com/matejeliash/goposc/internal/netinfo"
	"github.com/matejeliash/goposc/internal/numbermanager"
	"github.com/matejeliash/goposc/internal/pinger"
	"github.com/matejeliash/goposc/internal/portscanner"
)

type FlagParser struct {
	IpRange  []string
	Ports    []int
	FoundIps []string
}

func (fp *FlagParser) Parse() error {
	//ipRangeFlag := flag.String("ip-range", "", "IP range in start-end format (e.g. 192.168.1.10-192.168.1.20)")
	//ipFlag := flag.String("ip", "", "IP adress (e.g. 192.168.1.10)")
	ipsFlag := flag.String("ips", "", "IP adresses (e.g. 192.168.1.10,192.168.0.111)")
	pingFlag := flag.Bool("ping", false, "use ping to find if devices with specified id apdresses are active")
	portsFlag := flag.String("ports", "", "specified ports")
	portScanFlag := flag.Bool("port-scan", false, "use port scan to fin if service is active")
	getBannerFlag := flag.Bool("get-banner", false, "use port scan to fin if service is active")
	netInfoFlag := flag.Bool("net-info", false, "get basic network information")
	flag.Parse()

	if *ipsFlag != "" {
		ips, err := numbermanager.IpsFromPromInput(*ipsFlag)
		if err != nil {
			return err
		}
		fp.IpRange = append(fp.IpRange, ips...)
	}

	if *portsFlag != "" {
		ports, err := numbermanager.PortsFromPromInput(*portsFlag)
		if err != nil {
			return err
		}
		fp.Ports = append(fp.Ports, ports...)

	}

	if *pingFlag {
		p := pinger.NewPinger()
		fi := p.PingAllIPs(fp.IpRange)
		fp.FoundIps = append(fp.FoundIps, fi...)
	}

	if *portScanFlag {
		if len(fp.IpRange) == 0 {
			return fmt.Errorf("you must specify at least one ip adress")
		}

		if len(fp.Ports) == 0 {
			return fmt.Errorf("you must specify at least one port")
		}

		if len(fp.FoundIps) > 0 {

			for _, ip := range fp.FoundIps {
				portscanner.ScanPortsOfIP(ip, fp.Ports)
			}

		} else {
			for _, ip := range fp.IpRange {
				portscanner.ScanPortsOfIP(ip, fp.Ports)

			}

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

	if *netInfoFlag {
		ni := &netinfo.NetworkInfo{}

		ni.GetInfos()
		ni.PrintInfo()
	}

	return nil

}
