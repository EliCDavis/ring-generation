package main

import "github.com/EliCDavis/vector"

// Line represents a line segment
type Line struct {
	p1 *vector.Vector2
	p2 *vector.Vector2
}

// NewLine create a new
func NewLine(p1, p2 *vector.Vector2) *Line {
	return &Line{p1, p2}
}

// Intersection finds where two lines intersect
// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
func (l Line) Intersection(other *Line) *vector.Vector2 {

	s1_x := l.p2.X() - l.p1.X()
	s1_y := l.p2.Y() - l.p1.Y()

	s2_x := other.p2.X() - other.p1.X()
	s2_y := other.p2.Y() - other.p1.Y()

	s := (-s1_y*(l.p1.X()-other.p1.X()) + s1_x*(l.p1.Y()-other.p1.Y())) / (-s2_x*s1_y + s1_x*s2_y)
	t := (s2_x*(l.p1.Y()-other.p1.Y()) - s2_y*(l.p1.X()-other.p1.X())) / (-s2_x*s1_y + s1_x*s2_y)

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		return vector.NewVector2(l.p1.X()+(t*s1_x), l.p1.Y()+(t*s1_y))
	}

	return nil
}
