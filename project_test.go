//go:generate ./scripts/contract-gen.sh

package killcord

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	s1 := New()
	s2 := &Session{}

	if reflect.TypeOf(s1) != reflect.TypeOf(s2) {
		t.Fail()
	}
}
