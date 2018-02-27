// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package dscp provides helper functions to apply DSCP / ECN / CoS flags to sockets.
package dscp

import (
	"fmt"
	"net"
	"time"
)

// DialTCPWithTOS is similar to net.DialTCP but with the socket configured
// to the use the given ToS (Type of Service), to specify DSCP / ECN / class
// of service flags to use for incoming connections.
func DialTCPWithTOS(laddr, raddr *net.TCPAddr, tos byte) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", laddr, raddr)
	if err != nil {
		return nil, err
	}
	if err = setTOS(raddr.IP, conn, tos); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, err
}

// DialTimeoutWithTOS is similar to net.DialTimeout but with the socket configured
// to the use the given ToS (Type of Service), to specify DSCP / ECN / class
// of service flags to use for incoming connections.
func DialTimeoutWithTOS(network, address string, timeout time.Duration, tos byte) (net.Conn,
	error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	var ip net.IP
	// Unfortunately we have to explicitly switch on the address type here to
	// avoid calling net.ResolveIpAddr(), as this would resolve the address
	// again leading to a potentially different result.
	switch addr := conn.RemoteAddr().(type) {
	case *net.TCPAddr:
		ip = addr.IP
	case *net.UDPAddr:
		ip = addr.IP
	case *net.IPAddr:
		ip = addr.IP
	case *net.IPNet:
		ip = addr.IP
	default:
		conn.Close()
		return nil, fmt.Errorf("DialTimeoutWithTOS: cannot set TOS on a %s socket", network)
	}
	if err = setTOS(ip, conn, tos); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, err
}
