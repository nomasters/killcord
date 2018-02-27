// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package dscp

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"syscall"

	"golang.org/x/sys/unix"
)

// This works for the UNIX implementation of netFD, i.e. not on Windows and Plan9.
// conn must either implement syscall.Conn or be a TCPListener.
func setTOS(ip net.IP, conn interface{}, tos byte) error {
	var proto, optname int
	if ip.To4() != nil {
		proto = unix.IPPROTO_IP
		optname = unix.IP_TOS
	} else {
		proto = unix.IPPROTO_IPV6
		optname = unix.IPV6_TCLASS
	}

	switch c := conn.(type) {
	case syscall.Conn:
		return setTOSWithSyscallConn(proto, optname, c, tos)
	case *net.TCPListener:
		// This code is needed to support go1.9. In go1.10
		// *net.TCPListener implements syscall.Conn.
		return setTOSWithTCPListener(proto, optname, c, tos)
	}
	return fmt.Errorf("unsupported connection type: %T", conn)
}

func setTOSWithTCPListener(proto, optname int, conn *net.TCPListener, tos byte) error {
	// A kludge for pre-go1.10 to get the fd of a net.TCPListener
	value := reflect.ValueOf(conn)
	netFD := value.Elem().FieldByName("fd").Elem()
	fd := int(netFD.FieldByName("pfd").FieldByName("Sysfd").Int())
	if err := unix.SetsockoptInt(fd, proto, optname, int(tos)); err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func setTOSWithSyscallConn(proto, optname int, conn syscall.Conn, tos byte) error {
	syscallConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	var setsockoptErr error
	err = syscallConn.Control(func(fd uintptr) {
		if err := unix.SetsockoptInt(int(fd), proto, optname, int(tos)); err != nil {
			setsockoptErr = os.NewSyscallError("setsockopt", err)
		}
	})
	if setsockoptErr != nil {
		return setsockoptErr
	}
	return err
}
