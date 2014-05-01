package stack

import "testing"

func TestText(t *testing.T) {
	result := Text()
	expected := "Hello Worlds"
	if result != expected {
		t.Errorf("Expected %q got %q", expected, result)
	}
}
