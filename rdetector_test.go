package main

import (
	"testing"
)

func TestMain(m *testing.M) {
	sliceData := []string{"a", "b", "c"}
	element := "d"
	got := elementExists(sliceData, element)
	want := false

	if got != want {
		m.Run()
	}

	sliceData = []string{"a", "b", "c"}
	element = "b"
	got = elementExists(sliceData, element)
	want = true
	if got == want {
		m.Run()
	}
}
