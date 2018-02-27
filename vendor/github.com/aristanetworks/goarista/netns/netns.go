// Copyright (c) 2016 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package netns provides a utility function that allows a user to
// perform actions in a different network namespace
package netns

import (
	"fmt"
)

const (
	netNsRunDir = "/var/run/netns/"
	selfNsFile  = "/proc/self/ns/net"
)

// Callback is a function that gets called in a given network namespace.
// The user needs to check any errors from any calls inside this function.
type Callback func() error

// File descriptor interface so we can mock for testing
type handle interface {
	close() error
	fd() int
}

// The file descriptor associated with a network namespace
type nsHandle int

// setNsByName wraps setNs, allowing specification of the network namespace by name.
// It returns the file descriptor mapped to the given network namespace.
func setNsByName(nsName string) error {
	netPath := netNsRunDir + nsName
	handle, err := getNs(netPath)
	if err != nil {
		return fmt.Errorf("Failed to getNs: %s", err)
	}
	err = setNs(handle)
	handle.close()
	if err != nil {
		return fmt.Errorf("Failed to setNs: %s", err)
	}
	return nil
}
