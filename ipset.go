package main

import (
	"fmt"
	"net"
)

type IPSet struct {
	networks []*net.IPNet
}

func NewIPSet(networks []string) (*IPSet, error) {
	ips := &IPSet{}
	for _, net := range networks {
		ipnet, err := parseIPOrCIDR(net)
		if err != nil {
			return nil, err
		}
		ips.networks = append(ips.networks, ipnet)
	}
	return ips, nil
}

func parseIPOrCIDR(input string) (*net.IPNet, error) {
	_, ipnet, err := net.ParseCIDR(input)
	if err == nil {
		return ipnet, nil
	}

	// Not a CIDR, construct one from an IP
	ip := net.ParseIP(input)
	if ip == nil {
		return nil, fmt.Errorf("Failed to parse '%s' as an IP or CIDR", input)
	}
	ipnet = &net.IPNet{IP: ip}
	if ip.To4() != nil {
		ipnet.Mask = net.CIDRMask(32, 32)
	} else {
		ipnet.Mask = net.CIDRMask(128, 128)
	}
	return ipnet, nil
}

func (ips *IPSet) Contains(ipString string) bool {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return false
	}
	for _, net := range ips.networks {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}
