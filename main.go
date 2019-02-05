package main

import (
	"bufio"
	"io/ioutil"
	"math"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/pradeep-pyro/triangle"
	"golang.org/x/image/font"
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

	bottomLeftTexture *Vector2,
	topLeftTexture *Vector2,
	topRightTexture *Vector2,
	bottomRightTexture *Vector2,
) ([]*Polygon, error) {
	polys := make([]*Polygon, 2)

	poly, err := NewPolygonWithTexture(
		[]Vector3{*bottomLeft, *topLeft, *bottomRight},
		[]Vector3{*bottomLeft, *topLeft, *bottomRight},
		[]Vector2{*bottomLeftTexture, *topLeftTexture, *bottomRightTexture},
	)
	if err != nil {
		return nil, err
	}
	polys[0] = poly

	poly, err = NewPolygonWithTexture(
		[]Vector3{*topLeft, *topRight, *bottomRight},
		[]Vector3{*topLeft, *topRight, *bottomRight},
		[]Vector2{*topLeftTexture, *topRightTexture, *bottomRightTexture},
	)
	if err != nil {
		return nil, err
	}
	polys[1] = poly
	return polys, nil
}

func makeSquareFrom2D(bottomLeft, topRight *Vector2) []*Polygon {
	polys, _ := makeSquare(
		NewVector3(bottomLeft.X(), 0, bottomLeft.Y()),
		NewVector3(bottomLeft.X(), 0, topRight.Y()),
		NewVector3(topRight.X(), 0, topRight.Y()),
		NewVector3(topRight.X(), 0, bottomLeft.Y()),

		NewVector2(bottomLeft.X(), bottomLeft.Y()),
		NewVector2(bottomLeft.X(), topRight.Y()),
		NewVector2(topRight.X(), topRight.Y()),
		NewVector2(topRight.X(), bottomLeft.Y()),
	)
	return polys
}

// pointWithinBounds checks whether or not a point exists inside a rectangle
func pointWithinBounds(x float64, y float64, width float64, height float64, point Vector2) bool {
	return point.X() >= x && point.X() < x+width && point.Y() >= y && point.Y() < y+height
}

func lineIntersectsRectangle(x float64, y float64, width float64, height float64, line *Line) []*Vector2 {
	intersections := make([]*Vector2, 0)

	intersections = append(intersections, line.Intersection(NewLine(NewVector2(x, y), NewVector2(x+width, y))))
	intersections = append(intersections, line.Intersection(NewLine(NewVector2(x, y), NewVector2(x, y+height))))
	intersections = append(intersections, line.Intersection(NewLine(NewVector2(x+width, y), NewVector2(x+width, y+height))))
	intersections = append(intersections, line.Intersection(NewLine(NewVector2(x, y+height), NewVector2(x+width, y+height))))

	return intersections
}

func fill(width float64, height float64, shapes []*Shape) []*Polygon {

	numOfPoints := 0
	pointsPrefixSum := make([]int, len(shapes))
	for i, shape := range shapes {
		pointsPrefixSum[i] = numOfPoints
		numOfPoints += len(shape.points)
	}

	flatPoints := make([][2]float64, numOfPoints)

	for shapeIndex, shape := range shapes {
		for pointIndex, point := range shape.points {
			flatPoints[pointIndex+pointsPrefixSum[shapeIndex]][0] = point.X()
			flatPoints[pointIndex+pointsPrefixSum[shapeIndex]][1] = point.Y()
		}
	}

	segments := make([][2]int32, numOfPoints)

	for shapeIndex, shape := range shapes {
		for pointIndex := 0; pointIndex < len(shape.points); pointIndex++ {
			i := pointIndex + pointsPrefixSum[shapeIndex]
			segments[i][0] = int32(i)
			segments[i][1] = int32(((pointIndex + 1) % len(shape.points)) + pointsPrefixSum[shapeIndex])
		}
	}

	// Hole represented by a point lying inside it
	var holes = make([][2]float64, len(shapes))
	for i, _ := range shapes {
		holes[i][0] = 0
		holes[i][1] = 0.0
	}

	v, faces := triangle.ConstrainedDelaunay(flatPoints, segments, holes)

	betterPolys := make([]*Polygon, len(faces))
	for i, face := range faces {
		ourVerts := make([]Vector3, 0)
		ourVerts = append(ourVerts, *NewVector3(v[face[0]][0], 0, v[face[0]][1]))
		ourVerts = append(ourVerts, *NewVector3(v[face[1]][0], 0, v[face[1]][1]))
		ourVerts = append(ourVerts, *NewVector3(v[face[2]][0], 0, v[face[2]][1]))
		poly, _ := NewPolygon(ourVerts, ourVerts)
		betterPolys[i] = poly
	}
	return betterPolys
}

func carve(width float64, height float64, shapes []*Shape) []*Polygon {

	numOfPoints := 4
	pointsPrefixSum := make([]int, len(shapes))
	for i, shape := range shapes {
		pointsPrefixSum[i] = numOfPoints
		numOfPoints += len(shape.points)
	}

	flatPoints := make([][2]float64, numOfPoints)

	flatPoints[0] = [2]float64{0.0, 0.0}
	flatPoints[1] = [2]float64{0.0, height}
	flatPoints[2] = [2]float64{width, height}
	flatPoints[3] = [2]float64{width, 0.0}

	for shapeIndex, shape := range shapes {
		for pointIndex, point := range shape.points {
			flatPoints[pointIndex+pointsPrefixSum[shapeIndex]][0] = point.X()
			flatPoints[pointIndex+pointsPrefixSum[shapeIndex]][1] = point.Y()
		}
	}

	segments := make([][2]int32, numOfPoints)

	segments[0] = [2]int32{0, 1}
	segments[1] = [2]int32{1, 2}
	segments[2] = [2]int32{2, 3}
	segments[3] = [2]int32{3, 0}

	for shapeIndex, shape := range shapes {
		for pointIndex := 0; pointIndex < len(shape.points); pointIndex++ {
			i := pointIndex + pointsPrefixSum[shapeIndex]
			segments[i][0] = int32(i)
			segments[i][1] = int32(((pointIndex + 1) % len(shape.points)) + pointsPrefixSum[shapeIndex])
		}
	}

	// Hole represented by a point lying inside it
	var holes = make([][2]float64, len(shapes))
	for i, shape := range shapes {
		pointInShape := shape.RandomPointInShape()
		holes[i][0] = pointInShape.X()
		holes[i][1] = pointInShape.Y()
	}

	v, faces := triangle.ConstrainedDelaunay(flatPoints, segments, holes)

	betterPolys := make([]*Polygon, len(faces))
	for i, face := range faces {
		ourVerts := make([]Vector3, 3)
		ourVerts[0] = *NewVector3(v[face[0]][0], 0, v[face[0]][1])
		ourVerts[1] = *NewVector3(v[face[1]][0], 0, v[face[1]][1])
		ourVerts[2] = *NewVector3(v[face[2]][0], 0, v[face[2]][1])
		poly, _ := NewPolygon(ourVerts, ourVerts)
		betterPolys[i] = poly
	}
	return betterPolys
}

func main() {

	outerRadius := 1.2
	innerRadius := 1.0
	ringHeight := .8

	sides := 32
	polys := make([]*Polygon, sides*8)

	numTimesForTextureToRepeat := 8

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

	// model, err := NewModel(polys)
	// check(err)

	// f, err := os.Create("out.obj")
	// check(err)
	// defer f.Close()

	// w := bufio.NewWriter(f)
	// model.Save(w)
	// w.Flush()

	fontByteData, err := ioutil.ReadFile("./sample.ttf")
	check(err)
	parsedFont, err := truetype.Parse(fontByteData)
	check(err)

	// fontFace := truetype.NewFace(parsedFont, &truetype.Options{})

	textToEnscribe := "r"

	finalWord, err := NewModel([]*Polygon{})
	check(err)

	for _, char := range textToEnscribe {
		// log.Println(truetype.Index( - 97))

		glyph := truetype.GlyphBuf{}
		glyph.Load(parsedFont, 100, parsedFont.Index(char), font.HintingNone)

		letterPoints := make([]*Vector2, len(glyph.Points))
		for i, p := range glyph.Points {
			letterPoints[i] = NewVector2(float64(p.X), float64(p.Y))
		}
		shape := NewShape(letterPoints)
		shape.Scale(.1)

		left, right := shape.Split(5)

		if len(left) > 0 {
			lModel, _ := NewModel(carve(20.0, 20.0, left))
			finalWord = finalWord.Merge(lModel)
		}

		if len(right) > 0 {
			rModel, _ := NewModel(carve(20.0, 20.0, left))
			rModel = rModel.Translate(NewVector3(25, 0, 0))
			finalWord = finalWord.Merge(rModel)
		}

		//shape.Translate(NewVector2(7*float64(charIndex), 5))

	}

	f, err := os.Create("out.obj")
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	finalWord.Save(w)
	w.Flush()

}
