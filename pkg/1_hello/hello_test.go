package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello(" world")
	want := "hello World"
	if got != want {
		t.Errorf("got '%q' want '%q'", got, want)
	}
}
