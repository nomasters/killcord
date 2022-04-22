package killcord

import (
	"testing"
	"time"
)

func TestPublishKey(t *testing.T) {
	session := New()
	session.Config.Type = "bogus"
	if err := session.PublishKey(); err != nil {
		if err.Error() != "project type must be `owner` or `publisher`, received: bogus\n" {
			t.Fail()
		}
	}
}

func TestPublishThresholdReached(t *testing.T) {
	testTable := []struct {
		checkin   time.Time
		threshold int64
		expected  bool
	}{
		{time.Now(), defaultPublishThreshold, false},
		{time.Now().AddDate(0, 0, -1), defaultPublishThreshold, false},
		{time.Now().AddDate(0, 0, -4), defaultPublishThreshold, true},
	}

	for _, test := range testTable {
		if r := publishThresholdReached(test.checkin, test.threshold); r != test.expected {
			t.Errorf("checkin: %v threshold: %v was %v, expected %v", test.checkin, test.threshold, r, test.expected)
		}
	}
}
