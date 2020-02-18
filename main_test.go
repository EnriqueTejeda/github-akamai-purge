package main

import "testing"

func TestIsPurgeableAsset(t *testing.T) {
	got := isPurgeableAsset("image.jpeg")
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestAkamaiMakeUrl(t *testing.T) {
	got := akamaiMakeUrl("www.foo.com", "file.jpeg")
	want := "https://www.foo.com/file.jpeg"

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
