package tools

import (
	"errors"
	"go-ssip/app/common/consts"
	"net"
	"strconv"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.FreePortAddress)
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func GetFreePortInRange(ip string) (int, error) {
	startPort := 8000
	endPort := 9000
	for i := startPort; i < endPort; i++ {
		if i == 8080 || i == 8888 {
			continue
		}
		if isPortOpen(ip, i) {
			return i, nil
		}
	}
	return 0, errors.New("no free port")
}

func isPortOpen(ip string, port int) bool {
	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
