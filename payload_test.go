package killcord

import (
	"testing"
)

func TestSetPayloadRPCPath(t *testing.T) {
	s := New()
	s.setPayloadRPCPath()
	if payloadRPCPath != defaultpayloadRPCPath {
		t.Fail()
	}
	s.Config.Payload.RPCURL = "t1"
	s.setPayloadRPCPath()
	if payloadRPCPath != "t1" {
		t.Fail()
	}
	s.Options.Payload.RPCURL = "t2"
	s.setPayloadRPCPath()
	if payloadRPCPath != "t2" {
		t.Fail()
	}
}
