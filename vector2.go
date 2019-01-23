package main

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
