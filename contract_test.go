package killcord

import (
	"testing"
)

func TestSetEthereumRPCPath(t *testing.T) {
	s := New()
	s.setEthereumRPCPath()
	if ethereumRPCPath != defaultETHRPCPDev {
		t.Fail()
	}
	if s.Config.Contract.Mode != "testnet" {
		t.Fail()
	}
	s.Config.Contract.Mode = "mainnet"
	s.setEthereumRPCPath()
	if ethereumRPCPath != defaultETHRPCDProd {
		t.Fail()
	}
	s.Config.Contract.RPCURL = "t1"
	s.setEthereumRPCPath()
	if ethereumRPCPath != "t1" {
		t.Fail()
	}
	s.Options.Contract.RPCURL = "t2"
	if ethereumRPCPath != "t2" {
		t.Fail()
	}
}
