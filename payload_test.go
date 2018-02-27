package killcord

import (
	"testing"
)

func TestingSetPayloadRPCPath(t *testing.T) {
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

// func (s *Session) setPayloadRPCPath() {
// 	if s.Options.Payload.RPCURL != "" {
// 		payloadRPCPath = s.Options.Payload.RPCURL
// 		return
// 	}
// 	if s.Config.Payload.RPCURL != "" {
// 		payloadRPCPath = s.Config.Payload.RPCURL
// 		return
// 	}
// 	payloadRPCPath = defaultpayloadRPCPath
// }
