package net

import (
	"net"
	"strings"
)

// IPNetSet
type IPNetSet map[string]*net.IPNet

// ParseIPNets
func ParseIPNets(specs ...string) (IPNetSet, error) {
	ipnetset := make(IPNetSet)
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		_, ipnet, err := net.ParseCIDR(spec)
		if err != nil {
			return nil, err
		}
		k := ipnet.String()
		ipnetset[k] = ipnet
	}
	return ipnetset, nil
}

// Insert
func (s IPNetSet) Insert(items ...*net.IPNet) {
	for _, item := range items {
		s[item.String()] = item
	}
}

// Delete
func (s IPNetSet) Delete(items ...*net.IPNet) {
	for _, item := range items {
		delete(s, item.String())
	}
}

// Get
func (s IPNetSet) Get() *IPNetSet {
	return &s
}