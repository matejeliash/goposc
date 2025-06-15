package main

import (
	"flag"

	"github.com/matejeliash/goposc/cmd/tui/flagparser"
)

// func main() {
// var ni netinfo.NetworkInfo

// ni.GetInfos()
// fmt.Println(ni.NetInfos)
func main() {

	var fp flagparser.FlagParser
	err := fp.Parse()
	if err != nil {
		flag.Usage()
		panic(err)
	}

	//ni := &netinfo.NetworkInfo{}
	//ni.GetInfos()
	//fmt.Println(ni.NetInfos)

	//ipRange := flag.String("ip-range", "", "IP range")

	// var ni netinfo.NetworkInfo

	// ni.GetInfos()
	// fmt.Println(ni.NetInfos)

	//fmt.Println(GetPublicIp())

	// interfaces, err := net.Interfaces()
	// if err != nil {
	// 	panic(err)
	// }

	// //var ni []NetInterface

	// for _, iface := range interfaces {
	// 	addrs, err := iface.Addrs()
	// 	if err != nil {
	// 		continue
	// 	}

	// 	for _, addr := range addrs {
	// 		//fmt.Printf("%v\n", addr)
	// 		var ip net.IP
	// 		switch v := addr.(type) {
	// 		case *net.IPNet:
	// 			ip = v.IP
	// 		case *net.IPAddr:
	// 			ip = v.IP
	// 		}

	// 		if ip == nil || ip.IsLoopback() {
	// 			continue
	// 		}

	// 		var i NetInterface

	// 		if ip.To4() != nil {
	// 			//ipv4 = append(ipv4Addrs, ip.String())
	// 			i.Ipv4 = ip.String()
	// 			ipnet, ok := addr.(*net.IPNet)
	// 			if !ok || ipnet.IP.To4() == nil {
	// 				continue
	// 			}

	// 			networkIpv4 := calculateNetworkIP(ipnet.IP, ipnet.Mask)
	// 			networkSize, _ := ipnet.Mask.Size()
	// 			i.NetworkIpv4 = networkIpv4.To4().String() + "/" + strconv.Itoa(networkSize)
	// 		} else if ip.To16() != nil {
	// 			//ipv6 = append(ipv6Addrs, ip.String())
	// 			i.Ipv6 = ip.String()
	// 		}
	// 	}

	// }

	// subnet := "192.168.0." // change to your subnet prefix
	// var wg sync.WaitGroup
	// results := make(chan string, 256)

	// concurrencyLimit := make(chan struct{}, 100) // limit concurrent pings to 50

	// for i := 1; i <= 254; i++ {
	// 	ip := fmt.Sprintf("%s%d", subnet, i)
	// 	wg.Add(1)
	// 	concurrencyLimit <- struct{}{} // acquire a token
	// 	go func(ip string) {
	// 		defer func() { <-concurrencyLimit }() // release token
	// 		pingIP(ip, &wg, results)
	// 	}(ip)
	// }

	// // Close results channel after all pings finish
	// go func() {
	// 	wg.Wait()
	// 	close(results)
	// }()

	// fmt.Println("Devices found:")
	// for ip := range results {
	// 	fmt.Println(ip)
	// }
}
