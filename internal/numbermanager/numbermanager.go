package numbermanager

import (
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// last possible port number
const maxPort int = 65535

func GetPortsError(input string) error {
	return fmt.Errorf("incorrect string `%s` as --ports flag", input)
}

// get slice of app ports from --ips flag
func PortsFromPromInput(input string) ([]int, error) {

	// shirtcut to get all ports
	if input == "all" {

		// get all possible ports
		ports := make([]int, maxPort)
		for p := 1; p <= maxPort; p++ {
			ports = append(ports, p)
		}
		return ports, nil

	}
	// repomove spaces and split by `,` so multiple ranges and individual ports can be seleced
	inputOrg := input
	input = strings.ReplaceAll(input, " ", "")
	splittedByComma := strings.Split(input, ",")

	portMap := make(map[int]bool)
	for _, sbc := range splittedByComma {
		if strings.Contains(sbc, "-") {
			splittedByDash := strings.Split(sbc, "-")
			if len(splittedByDash) != 2 {
				return nil, GetPortsError(inputOrg)
			}
			portA, err := getPort(splittedByDash[0])
			if err != nil {
				return nil, GetPortsError(inputOrg)
			}
			portB, err := getPort(splittedByDash[1])
			if err != nil {
				return nil, GetPortsError(inputOrg)
			}
			if isPortCorrect(portA) && isPortCorrect(portB) && portA < portB {
				for p := portA; p <= portB; p++ {
					portMap[p] = true
				}
			} else {
				return nil, GetPortsError(inputOrg)

			}

		} else {
			port, err := getPort(sbc)
			if err != nil {
				return nil, GetPortsError(inputOrg)
			}
			portMap[port] = true

		}

	}
	var ports []int
	for k, v := range portMap {
		if v {
			ports = append(ports, k)
		}
	}
	return ports, nil
}

func getPort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, err
	} else {
		return port, nil
	}
}

func isPortCorrect(port int) bool {
	if port <= 0 || port > maxPort {
		return false
	}

	return true
}

func isIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && ip.To4() != nil
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

func isValidDomain(domain string) bool {
	var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}

func ParseIpRange(startIP, endIP string) ([]string, error) {

	if startIP == "" || endIP == "" {
		return nil, fmt.Errorf("invalid IP range")
	}

	startIPInt := ipToUint32(net.ParseIP(startIP))
	endIPInt := ipToUint32(net.ParseIP(endIP))
	if startIPInt >= endIPInt {
		return nil, fmt.Errorf("invalid IP range")
	}

	var ipRange []string

	for i := startIPInt; i <= endIPInt; i++ {
		ipRange = append(ipRange, uint32ToIP(i).String())
	}

	return ipRange, nil
}

func IpsFromPromInput(input string) ([]string, error) {

	// shirtcut to get all ports
	// remove spaces and split by `,` so multiple ranges and individual ports can be seleced
	inputOrg := input
	input = strings.ReplaceAll(input, " ", "")
	splittedByComma := strings.Split(input, ",")

	ipsMap := make(map[string]bool)
	for _, sbc := range splittedByComma {
		// take care of range
		if strings.Contains(sbc, "~") {
			splittedByDash := strings.Split(sbc, "~")
			if len(splittedByDash) != 2 {
				return nil, GetPortsError(inputOrg)
			}
			// check of correct IPv4
			if !isIPv4(splittedByDash[0]) || !isIPv4(splittedByDash[1]) {
				return nil, GetPortsError(inputOrg)
			}
			ipRange, err := ParseIpRange(splittedByDash[0], splittedByDash[1])
			if err != nil {
				return nil, GetPortsError(inputOrg)
			}
			for _, IP := range ipRange {
				ipsMap[IP] = true
			}

			// one IP / donaibn
		} else {
			if !isIPv4(sbc) && !isValidDomain(sbc) && sbc != "localhost" {
				return nil, GetPortsError(inputOrg)
			}
			ipsMap[sbc] = true

		}

	}
	var ips []string
	for k, v := range ipsMap {
		if v {
			ips = append(ips, k)
		}
	}
	return ips, nil
}
