/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package net

import (
	"net"
	"strings"
)

// IPNetSet maps string to net.IPNet.
type IPNetSet map[string]*net.IPNet

// ParseIPNets parses string slice to IPNetSet.
func ParseIPNets(specs ...string) (IPNetSet, error) {
	ipnetset := make(IPNetSet)
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		_, ipnet, err := net.ParseCIDR(spec)
		if err != nil {
			return nil, err
		}
		k := ipnet.String() // In case of normalization
		ipnetset[k] = ipnet
	}
	return ipnetset, nil
}

// Insert adds items to the set.
func (s IPNetSet) Insert(items ...*net.IPNet) {
	for _, item := range items {
		s[item.String()] = item
	}
}

// Delete removes all items from the set.
func (s IPNetSet) Delete(items ...*net.IPNet) {
	for _, item := range items {
		delete(s, item.String())
	}
}
