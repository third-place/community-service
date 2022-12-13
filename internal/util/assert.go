package util

import (
	"testing"
)

func Assert(t *testing.T, condition bool) {
	if condition {
		return
	}
	t.Fail()
}
