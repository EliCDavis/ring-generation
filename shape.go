package main

import "math"

// Shape is a flat (2D) arrangement of points.
type Shape struct {
	points []*Vector2
	center *Vector2
}

// NewShape creates a new shape with a center computed by averaging points positions
func NewShape(bounds []*Vector2) *Shape {
	center := NewVector2(0, 0)
	for _, point := range bounds {
		center = center.Add(point)
	}

	return &Shape{bounds, center.Scale(1.0 / float64(len(bounds)))}
}

// NewShapeWithCustomCenter creates a shape with a center you set
func NewShapeWithCustomCenter(bounds []*Vector2, center *Vector2) *Shape {
	return &Shape{bounds, center}
}

// Translate Moves all points over by the specified amount
func (s *Shape) Translate(amount *Vector2) {
	for i, point := range s.points {
		s.points[i] = point.Add(amount)
	}
}

// Len returns the number of points in the polygon
func (s Shape) Len() int {
	return len(s.points)
}

// Swap switches two points indeces so the polygon is ordered a different way
func (s *Shape) Swap(i, j int) {
	s.points[i], s.points[j] = s.points[j], s.points[i]
}

// PointClosestToCenter returns the index of the closest point to the center of the polygon
func (s Shape) PointClosestToCenter() int {
	bestDistance := math.Inf(0)
	curPoint := -1

	for i, point := range s.points {
		dist := s.center.Distance(point)
		if dist < bestDistance {
			bestDistance = dist
			curPoint = i
		}
	}

	return curPoint
}

// Less determines which point is more orriented more clockwise from the center than the other
func (s Shape) Less(i, j int) bool {
	a := s.points[i]
	b := s.points[j]

	if a.x-s.center.x >= 0 && b.x-s.center.x < 0 {
		return true
	}

	if a.x-s.center.x < 0 && b.x-s.center.x >= 0 {
		return false
	}

	if a.x-s.center.x == 0 && b.x-s.center.x == 0 {
		if a.y-s.center.y >= 0 || b.y-s.center.y >= 0 {
			return a.y > b.y
		}
		return b.y > a.y
	}

	// compute the cross product of vectors (center -> a) x (center -> b)
	det := (a.x-s.center.x)*(b.y-s.center.y) - (b.x-s.center.x)*(a.y-s.center.y)
	if det < 0 {
		return true
	}
	if det > 0 {
		return false
	}

	// points a and b are on the same line from the center
	// check which point is closer to the center
	d1 := (a.x-s.center.x)*(a.x-s.center.x) + (a.y-s.center.y)*(a.y-s.center.y)
	d2 := (b.x-s.center.x)*(b.x-s.center.x) + (b.y-s.center.y)*(b.y-s.center.y)
	return d1 > d2
}
