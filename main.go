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

func makeSquare(
	bottomLeft *Vector3,
	topLeft *Vector3,
	topRight *Vector3,
	bottomRight *Vector3,

	bottomLeftUV *Vector2,
	topLeftUV *Vector2,
	topRightUV *Vector2,
	bottomRightUV *Vector2,
) ([]*Polygon, error) {
	polys := make([]*Polygon, 2)

	poly, err := NewPolygonWithTexture(
		[]Vector3{*bottomLeft, *topLeft, *bottomRight},
		[]Vector3{*bottomLeft, *topLeft, *bottomRight},
		[]Vector2{*bottomLeftUV, *topLeftUV, *bottomRightUV},
	)
	if err != nil {
		return nil, err
	}
	polys[0] = poly

	poly, err = NewPolygonWithTexture(
		[]Vector3{*topLeft, *topRight, *bottomRight},
		[]Vector3{*topLeft, *topRight, *bottomRight},
		[]Vector2{*topLeftUV, *topRightUV, *bottomRightUV},
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

	numTimesForTextureToRepeat := 4

	angleIncrement := (1.0 / float64(sides)) * 2.0 * math.Pi
	for sideIndex := 0; sideIndex < sides; sideIndex++ {
		angle := angleIncrement * float64(sideIndex)
		angleNext := angleIncrement * (float64(sideIndex) + 1)

		bottomLeftUV := NewVector2(float64(sideIndex)/float64(sides/numTimesForTextureToRepeat), 0)
		for bottomLeftUV.x > 1 {
			bottomLeftUV.x -= 1.0
		}

		topLeftUV := NewVector2(float64(sideIndex)/float64(sides/numTimesForTextureToRepeat), 1.0)
		for topLeftUV.x > 1 {
			topLeftUV.x -= 1.0
		}

		topRightUV := NewVector2(float64(sideIndex+1)/float64(sides/numTimesForTextureToRepeat), 1.0)
		for topRightUV.x > 1 {
			topRightUV.x -= 1.0
		}

		bottomRightUV := NewVector2(float64(sideIndex+1)/float64(sides/numTimesForTextureToRepeat), 0)
		for bottomRightUV.x > 1 {
			bottomRightUV.x -= 1.0
		}

		// outer
		square, err := makeSquare(
			NewVector3(math.Cos(angle)*outerRadius, 0, math.Sin(angle)*outerRadius),
			NewVector3(math.Cos(angle)*outerRadius, ringHeight, math.Sin(angle)*outerRadius),
			NewVector3(math.Cos(angleNext)*outerRadius, ringHeight, math.Sin(angleNext)*outerRadius),
			NewVector3(math.Cos(angleNext)*outerRadius, 0, math.Sin(angleNext)*outerRadius),
			bottomLeftUV,
			topLeftUV,
			topRightUV,
			bottomRightUV,
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
			bottomRightUV,
			topRightUV,
			topLeftUV,
			bottomLeftUV,
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
			bottomLeftUV,
			topLeftUV,
			topRightUV,
			bottomRightUV,
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
			bottomLeftUV,
			topLeftUV,
			topRightUV,
			bottomRightUV,
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
