package netinfo

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/fatih/color"
)

type NetworkInfo struct {
	PublicIp string
	NetInfos []NetInterface
	AllIPs   []string
}

type NetInterface struct {
	Name        string
	Ipv4        string
	Ipv6        string
	NetworkIpv6 string
	NetworkIpv4 string
}

func (ni *NetworkInfo) GetInfos() error {

	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	//var ni []NetInterface

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		var i NetInterface
		i.Name = iface.Name

		for _, addr := range addrs {
			//fmt.Printf("%v\n", addr)

			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			//if ip == nil || ip.IsLoopback() {
			//continue
			//}

			if ip.To4() != nil {
				//ipv4 = append(ipv4Addrs, ip.String())
				i.Ipv4 = ip.String()
				ipnet, ok := addr.(*net.IPNet)
				if !ok || ipnet.IP.To4() == nil {
					continue
				}

				networkIpv4 := calculateNetworkIP(ipnet.IP, ipnet.Mask)
				networkSize, _ := ipnet.Mask.Size()
				i.NetworkIpv4 = networkIpv4.To4().String() + "/" + strconv.Itoa(networkSize)
			} else if ip.To16() != nil {
				i.Ipv6 = ip.String()
			}

		}

		ni.NetInfos = append(ni.NetInfos, i)
	}

	return nil
}

func calculateNetworkIP(ip net.IP, mask net.IPMask) net.IP {
	ip = ip.To4()
	network := make(net.IP, len(ip))
	for i := 0; i < len(ip); i++ {
		network[i] = ip[i] & mask[i]
	}
	return network
}

// Returns a slice of all usable host IP addresses in the network
func getAllHosts(cidr string) ([]net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []net.IP

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		// Make a copy since ip will change in the loop
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		ips = append(ips, ipCopy)
	}

	// Remove network and broadcast addresses
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}

	return ips, nil
}

// Increments an IP address by 1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}

func (ni *NetworkInfo) PrintInfo() {

	green := color.New(color.FgGreen).SprintFunc()
	for _, i := range ni.NetInfos {

		fmt.Println(green("Interface"), i.Name+":")
		fmt.Println(green("   Ipv4:"), i.Ipv4)
		fmt.Println(green("   Subnet Ipv4:"), i.NetworkIpv4)
		fmt.Println(green("   Ipv6:"), i.Ipv6)
	}

	ip, err := ni.GetPublicIp()
	if err == nil {
		fmt.Println(green("\nPublic Ipv4:"), ip)

	}

}

func (ni *NetworkInfo) GetPublicIp() (string, error) {

	resp, err := http.Get("https://ifconfig.me/ip")

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	ni.PublicIp = string(body)
	return string(body), nil

}
