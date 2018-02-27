// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package dscp_test

import (
	"net"
	"testing"

	"github.com/aristanetworks/goarista/dscp"
)

func TestDialTCPWithTOS(t *testing.T) {
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer listen.Close()

	done := make(chan struct{})
	go func() {
		conn, err := listen.Accept()
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		buf := []byte{'!'}
		conn.Write(buf)
		n, err := conn.Read(buf)
		if n != 1 || err != nil {
			t.Fatalf("Read returned %d / %s", n, err)
		} else if buf[0] != '!' {
			t.Fatalf("Expected to read '!' but got %q", buf)
		}
		close(done)
	}()
	conn, err := dscp.DialTCPWithTOS(nil, listen.Addr().(*net.TCPAddr), 40)
	if err != nil {
		t.Fatal("Connection failed:", err)
	}
	defer conn.Close()
	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	if n != 1 || err != nil {
		t.Fatalf("Read returned %d / %s", n, err)
	} else if buf[0] != '!' {
		t.Fatalf("Expected to read '!' but got %q", buf)
	}
	conn.Write(buf)
	<-done
}
