package utils

import (
	"fmt"
	"net"
)

func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsValidPort(port int) bool {
	return port > 0 && port < 65536
}

func CheckBind(ip string, port int) bool {
	if !IsValidIP(ip) || !IsValidPort(port) {
		return false
	}

	addr := fmt.Sprintf("%s:%d", ip, port)

	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	defer tcpListener.Close()

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return false
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return false
	}
	defer udpConn.Close()

	return true
}
