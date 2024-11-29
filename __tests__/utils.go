package tests

import "testing"

func SendError(t *testing.T, got interface{}, want interface{}) {
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
