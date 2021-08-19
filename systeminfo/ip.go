package systeminfo

import (
	"errors"
	"fmt"
	"github.com/kbinani/win"
	"go-windows-monitor/utils/log"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"unsafe"
)

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addr, err := iface.MulticastAddrs()
		log.Infof("==%s,", addr)

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("failed")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetDns() []string {
	dns := make([]string, 0)
	info := win.FIXED_INFO_W2KSP1{}
	size := uint32(unsafe.Sizeof(info))
	r := win.GetNetworkParams(&info, &size)
	fmt.Println("r=", r, "size=", size)
	if r == 0 {
		for ai := &info.DnsServerList; ai != nil; ai = ai.Next {
			d := fmt.Sprintf("%v.%v.%v.%v", ai.Context&0xFF, (ai.Context>>8)&0xFF, (ai.Context>>16)&0xFF, (ai.Context>>24)&0xFF)
			dns = append(dns, d)
		}
	} else {
		newBuffers := make([]byte, size)
		netParams := (win.PFIXED_INFO)(unsafe.Pointer(&newBuffers[0]))
		win.GetNetworkParams(netParams, &size)
		for ai := &netParams.DnsServerList; ai != nil; ai = ai.Next {
			d := fmt.Sprintf("%v.%v.%v.%v", ai.Context&0xFF, (ai.Context>>8)&0xFF, (ai.Context>>16)&0xFF, (ai.Context>>24)&0xFF)
			dns = append(dns, d)
		}
	}

	// r == win.ValueOverflow
	return dns
}

func LocalIp() (string, string) {
	var finalIp string
	var gateWay string
	cmd := exec.Command("cmd", "/c", "ipconfig")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", ""
	}

	defer out.Close()
	if err := cmd.Start(); err != nil {
		return "", ""
	}

	opBytes, err := ioutil.ReadAll(out)
	if err != nil {
		return "", ""
	}

	str := string(opBytes)
	log.Infof("str=%s", str)
	var strs = strings.Split(str, "\r\n")
	if 0 != len(strs) {
		var havingFinalIp4 bool = false
		var cnt int = 0
		for index, value := range strs {
			vidx := strings.Index(value, "IPv4")
			//说明已经找到该ip
			if vidx != -1 {
				ip4lines := strings.Split(value, ":")
				if len(ip4lines) == 2 {
					cnt = index
					havingFinalIp4 = true
					ip4str := ip4lines[1]
					finalIp = strings.TrimSpace(ip4str)
				}
			}


			vidx = strings.Index(value, "默认网关")
			log.Infof("===%s %d", value,  vidx)
			if vidx != -1 {
				ip4lines := strings.Split(value, ":")
				log.Infof("value=%+v", ip4lines)

				if len(ip4lines) == 2 {
					ip4str := ip4lines[1]
					gateWay = strings.TrimSpace(ip4str)
				}
			}

			if havingFinalIp4 && index == cnt+2 {
				lindex := strings.Index(value, ":")
				if -1 != lindex {
					lines := strings.Split(value, ":")
					if len(lines) == 2 {
						fip := lines[1]
						if strings.TrimSpace(fip) != "" {
							break
						}
					}
				}
				havingFinalIp4 = false
				finalIp = ""
			}
		}
	}

	return finalIp, gateWay

}
