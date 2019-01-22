package main

import (
	"bufio"
	"math"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func makeSquare(bottomLeft *Vector3, topLeft *Vector3, topRight *Vector3, bottomRight *Vector3) ([]*Polygon, error) {
	polys := make([]*Polygon, 2)

	poly, err := NewPolygon(
		[]Vector3{*bottomLeft, *topLeft, *bottomRight},
		[]Vector3{*bottomLeft, *topLeft, *bottomRight},
	)
	if err != nil {
		return nil, err
	}
	polys[0] = poly

	poly, err = NewPolygon(
		[]Vector3{*topLeft, *topRight, *bottomRight},
		[]Vector3{*topLeft, *topRight, *bottomRight},
	)
	if err != nil {
		return nil, err
	}
	polys[1] = poly
	return polys, nil
}

func main() {

	outerRadius := 1.2
	innerRadius := 1.0
	ringHeight := .8

	sides := 32
	polys := make([]*Polygon, sides*8)

	angleIncrement := (1.0 / float64(sides)) * 2.0 * math.Pi
	for sideIndex := 0; sideIndex < sides; sideIndex++ {
		angle := angleIncrement * float64(sideIndex)
		angleNext := angleIncrement * (float64(sideIndex) + 1)

		// outer
		square, err := makeSquare(
			NewVector3(math.Cos(angle)*outerRadius, 0, math.Sin(angle)*outerRadius),
			NewVector3(math.Cos(angle)*outerRadius, ringHeight, math.Sin(angle)*outerRadius),
			NewVector3(math.Cos(angleNext)*outerRadius, ringHeight, math.Sin(angleNext)*outerRadius),
			NewVector3(math.Cos(angleNext)*outerRadius, 0, math.Sin(angleNext)*outerRadius),
		)

		check(err)
		polys[(sideIndex * 8)] = square[0]
		polys[(sideIndex*8)+1] = square[1]

		// inner
		square, err = makeSquare(
			NewVector3(math.Cos(angleNext)*innerRadius, 0, math.Sin(angleNext)*innerRadius),
			NewVector3(math.Cos(angleNext)*innerRadius, ringHeight, math.Sin(angleNext)*innerRadius),
			NewVector3(math.Cos(angle)*innerRadius, ringHeight, math.Sin(angle)*innerRadius),
			NewVector3(math.Cos(angle)*innerRadius, 0, math.Sin(angle)*innerRadius),
		)

		check(err)
		polys[(sideIndex*8)+2] = square[0]
		polys[(sideIndex*8)+3] = square[1]

		// top
		square, err = makeSquare(
			NewVector3(math.Cos(angle)*outerRadius, ringHeight, math.Sin(angle)*outerRadius),
			NewVector3(math.Cos(angle)*innerRadius, ringHeight, math.Sin(angle)*innerRadius),
			NewVector3(math.Cos(angleNext)*innerRadius, ringHeight, math.Sin(angleNext)*innerRadius),
			NewVector3(math.Cos(angleNext)*outerRadius, ringHeight, math.Sin(angleNext)*outerRadius),
		)

		check(err)
		polys[(sideIndex*8)+4] = square[0]
		polys[(sideIndex*8)+5] = square[1]

		// bottom
		square, err = makeSquare(
			NewVector3(math.Cos(angle)*innerRadius, 0, math.Sin(angle)*innerRadius),
			NewVector3(math.Cos(angle)*outerRadius, 0, math.Sin(angle)*outerRadius),
			NewVector3(math.Cos(angleNext)*outerRadius, 0, math.Sin(angleNext)*outerRadius),
			NewVector3(math.Cos(angleNext)*innerRadius, 0, math.Sin(angleNext)*innerRadius),
		)

		check(err)
		polys[(sideIndex*8)+6] = square[0]
		polys[(sideIndex*8)+7] = square[1]
	}

	model, err := NewModel(polys)
	check(err)

	f, err := os.Create("out.obj")
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	model.Save(w)
	w.Flush()

}
