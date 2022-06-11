package main

import "testing"

func Test_dummy(t *testing.T) {
	result := 3
	if result != 3 {
		t.Error("incorrect result: expect 3, got", result)
	}
}
