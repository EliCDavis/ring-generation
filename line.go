package main

// Line represents a line segment
type Line struct {
	p1 *Vector2
	p2 *Vector2
}

// NewLine create a new
func NewLine(p1, p2 *Vector2) *Line {
	return &Line{p1, p2}
}

// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
func (l Line) Intersection(other *Line) *Vector2 {

	s1_x := l.p2.x - l.p1.x
	s1_y := l.p2.y - l.p1.y

	s2_x := other.p2.x - other.p1.x
	s2_y := other.p2.y - other.p1.y

	s := (-s1_y*(l.p1.x-other.p1.x) + s1_x*(l.p1.y-other.p1.y)) / (-s2_x*s1_y + s1_x*s2_y)
	t := (s2_x*(l.p1.y-other.p1.y) - s2_y*(l.p1.x-other.p1.x)) / (-s2_x*s1_y + s1_x*s2_y)

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		return NewVector2(l.p1.x+(t*s1_x), l.p1.y+(t*s1_y))
	}

	return nil
}
