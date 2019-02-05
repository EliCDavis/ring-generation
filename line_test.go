package main

import "testing"

func TestIntersect(t *testing.T) {
	l1 := NewLine(NewVector2(0, 0), NewVector2(1, 1))
	l2 := NewLine(NewVector2(0, 1), NewVector2(1, 0))
	if doIntersect(l1, l2) == false {
		t.Error("Lines should have interesected")
	}
}
