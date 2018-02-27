package killcord

import (
	"testing"
)

// need to write real tests for this...

func TestGetStatus(t *testing.T) {
	session := New()
	err := session.GetStatus()
	if err != nil {
		t.Fail()
	}
}

func TestGetProjectStatus(t *testing.T) {
	session := New()
	err := session.getProjectStatus()
	if err != nil {
		t.Fail()
	}
}

func TestGetContractStatus(t *testing.T) {
	session := New()
	err := session.getContractStatus()
	if err != nil {
		t.Fail()
	}

}

func TestGetPublisherStatus(t *testing.T) {
	session := New()
	err := session.getPublisherStatus()
	if err != nil {
		t.Fail()
	}
}

func TestGetPayloadStatus(t *testing.T) {
	session := New()
	err := session.getPayloadStatus()
	if err != nil {
		t.Fail()
	}
}
