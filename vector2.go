package main

import "math"

type Vector2 struct {
	x float64
	y float64
}

func NewVector2(x float64, y float64) *Vector2 {
	return &Vector2{
		x: x,
		y: y,
	}
}

func (v *Vector2) X() float64 {
	return v.x
}

func (v *Vector2) Y() float64 {
	return v.y
}

// Add returns a vector that is the result of two vectors added together
func (v Vector2) Add(other *Vector2) *Vector2 {
	return &Vector2{
		x: v.x + other.x,
		y: v.y + other.y,
	}
}

// Scale multiplies each axis by the specifid value
func (v Vector2) Scale(other float64) *Vector2 {
	return &Vector2{
		x: v.x * other,
		y: v.y * other,
	}
}

// Distance is the euclidian distance between two points
func (v Vector2) Distance(other *Vector2) float64 {
	return math.Sqrt(math.Pow(other.x-v.x, 2.0) + math.Pow(other.y-v.y, 2.0))
}
