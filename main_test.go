// A simple test

package main

import "testing"

func TestHello(t *testing.T) {
	expected := "Hello"
	actual := hello()
	if actual != expected {
		t.Errorf("Test failed", expected, actual)
	}
}
