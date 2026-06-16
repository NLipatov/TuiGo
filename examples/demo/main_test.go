package main

import "testing"

func TestKeyTextLabelNamesZeroWidthText(t *testing.T) {
	got := keyTextLabel("\u0301")
	want := "U+0301"
	if got != want {
		t.Fatalf("keyTextLabel() = %q, want %q", got, want)
	}
}

func TestKeyTextLabelNamesSpace(t *testing.T) {
	got := keyTextLabel(" ")
	want := "space"
	if got != want {
		t.Fatalf("keyTextLabel() = %q, want %q", got, want)
	}
}
