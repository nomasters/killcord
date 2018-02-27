package killcord

import (
	"testing"
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
