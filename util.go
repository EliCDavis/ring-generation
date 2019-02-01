package main

import "math"

type Orientation int

const (
	Colinear Orientation = iota
	Clockwise
	Counterclockwise
)

// To find orientation of ordered triplet (p, q, r).
// The function returns following values
// 0 --> p, q and r are colinear
// 1 --> Clockwise
// 2 --> Counterclockwise
func calculateOrientation(p, q, r *Vector2) Orientation {
	// See https://www.geeksforgeeks.org/orientation-3-ordered-points/
	// for details of below formula.
	val := ((q.y - p.y) * (r.x - q.x)) - ((q.x - p.x) * (r.y - q.y))

	if val == 0 {
		return Colinear
	}

	if val > 0 {
		return Clockwise
	}

	return Counterclockwise
}

// Intersect determines whether two lines intersect eachother
func doIntersect(l, other *Line) bool {
	// Find the four orientations needed for general and
	// special cases
	o1 := calculateOrientation(l.p1, l.p2, other.p1)
	o2 := calculateOrientation(l.p1, l.p2, other.p2)
	o3 := calculateOrientation(other.p1, other.p2, l.p1)
	o4 := calculateOrientation(other.p1, other.p2, l.p2)

	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special Cases
	// l.p1, l.p2 and other.p1 are colinear and other.p1 lies on segment l.p1l.p2
	if o1 == Colinear && onSegment(l.p1, other.p1, l.p2) {
		return true
	}

	// l.p1, l.p2 and other.p2 are colinear and other.p2 lies on segment l.p1l.p2
	if o2 == Colinear && onSegment(l.p1, other.p2, l.p2) {
		return true
	}

	// p2, other.p2 and l.p1 are colinear and l.p1 lies on segment p2other.p2
	if o3 == Colinear && onSegment(other.p1, l.p1, other.p2) {
		return true
	}

	// p2, other.p2 and l.p2 are colinear and l.p2 lies on segment p2other.p2
	if o4 == 0 && onSegment(other.p1, l.p2, other.p2) {
		return true
	}

	return false // Doesn't fall in any of the above cases
}

// Given three colinear points p, q, r, the function checks if
// point q lies on line segment 'pr'
func onSegment(p, q, r *Vector2) bool {
	return q.x <= math.Max(p.x, r.x) && q.x >= math.Min(p.x, r.x) && q.y <= math.Max(p.y, r.y) && q.y >= math.Min(p.y, r.y)
}

// Returns true if the point p lies inside the polygon[] with n vertices
func isInside(bounds []Vector2, p *Vector2) bool {

	// There must be at least 3 vertices in polygon[]
	if len(bounds) < 3 {
		return false
	}

	// Create a point for line segment from p to infinite
	extreme := NewVector2(math.Inf(0), p.y)

	// Count intersections of the above line with sides of polygon
	count := 0
	i := 0
	for {
		next := (i + 1) % len(bounds)

		// Check if the line segment from 'p' to 'extreme' intersects
		// with the line segment from 'polygon[i]' to 'polygon[next]'
		if doIntersect(NewLine(&bounds[i], &bounds[next]), NewLine(p, extreme)) {
			// If the point 'p' is colinear with line segment 'i-next',
			// then check if it lies on segment. If it lies, return true,
			// otherwise false
			if calculateOrientation(&bounds[i], p, &bounds[next]) == Colinear {
				return onSegment(&bounds[i], p, &bounds[next])
			}

			count++
		}
		i = next
		if i == 0 {
			break
		}
	}

	// Return true if count is odd, false otherwise
	return count%2 == 1
}
