package vnet

import (
	"context"
	"errors"
	"fmt"
	"net"
)

func GetServAddr(servAddr string) (string, error) {
	addrTcp, err := net.ResolveTCPAddr(networkTypeTCP, servAddr)
	if err != nil {
		return "", err
	}

	addr := addrTcp.String()
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	if len(host) == 0 {
		return GetLocalHost(addrTcp)
	}

	return addr, nil
}

// Open port listening and return the service address
func ListenServAddr(ctx context.Context, addr string) (net.Listener, string, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, "", err
	}

	netListen, err := net.Listen(tcpAddr.Network(), tcpAddr.String())
	if err != nil {
		return nil, "", err
	}

	return netListen, tcpAddr.String(), nil
}

func GetLocalHost(addrTcp net.Addr) (string, error) {
	addr := addrTcp.String()
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if len(host) == 0 {
		host = "0.0.0.0"
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return "", fmt.Errorf("net.ParseIP error, host: %s", host)
	}

	realAddr := addr
	if ip.IsUnspecified() {
		internalIP, err := GetInternalIP()
		if err != nil {
			return "", err
		}
		realAddr = net.JoinHostPort(internalIP, port)
	}

	return realAddr, nil
}

func GetInternalIP() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no internal ip")
}
