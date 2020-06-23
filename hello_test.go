package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	want := "Selam Dünya."

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
