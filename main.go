package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/EliCDavis/mesh"
	"github.com/EliCDavis/vector"
	"github.com/golang/freetype/truetype"
	"github.com/pradeep-pyro/triangle"
	"golang.org/x/image/font"
)

func check(e error) {
	if e != nil {
		log.Panicln("Error1")
		panic(e)
	}
}

func makeSquare(
	bottomLeft vector.Vector3,
	topLeft vector.Vector3,
	topRight vector.Vector3,
	bottomRight vector.Vector3,

	bottomLeftTexture vector.Vector2,
	topLeftTexture vector.Vector2,
	topRightTexture vector.Vector2,
	bottomRightTexture vector.Vector2,
) ([]mesh.Polygon, error) {
	polys := make([]mesh.Polygon, 2)

	poly, err := mesh.NewPolygonWithTexture(
		[]vector.Vector3{bottomLeft, topLeft, bottomRight},
		[]vector.Vector3{bottomLeft, topLeft, bottomRight},
		[]vector.Vector2{bottomLeftTexture, topLeftTexture, bottomRightTexture},
	)
	if err != nil {
		return nil, err
	}
	polys[0] = poly

	poly, err = mesh.NewPolygonWithTexture(
		[]vector.Vector3{topLeft, topRight, bottomRight},
		[]vector.Vector3{topLeft, topRight, bottomRight},
		[]vector.Vector2{topLeftTexture, topRightTexture, bottomRightTexture},
	)
	if err != nil {
		return nil, err
	}
	polys[1] = poly
	return polys, nil
}

func makeSquareFrom2D(bottomLeft, topRight *vector.Vector2) []mesh.Polygon {
	polys, _ := makeSquare(
		vector.NewVector3(bottomLeft.X(), 0, bottomLeft.Y()),
		vector.NewVector3(bottomLeft.X(), 0, topRight.Y()),
		vector.NewVector3(topRight.X(), 0, topRight.Y()),
		vector.NewVector3(topRight.X(), 0, bottomLeft.Y()),

		vector.NewVector2(bottomLeft.X(), bottomLeft.Y()),
		vector.NewVector2(bottomLeft.X(), topRight.Y()),
		vector.NewVector2(topRight.X(), topRight.Y()),
		vector.NewVector2(topRight.X(), bottomLeft.Y()),
	)
	return polys
}

// pointWithinBounds checks whether or not a point exists inside a rectangle
func pointWithinBounds(x float64, y float64, width float64, height float64, point vector.Vector2) bool {
	return point.X() >= x && point.X() < x+width && point.Y() >= y && point.Y() < y+height
}

func lineIntersectsRectangle(x float64, y float64, width float64, height float64, line mesh.Line) []vector.Vector2 {
	intersections := make([]vector.Vector2, 0)

	point, err := line.Intersection(mesh.NewLine(vector.NewVector2(x, y), vector.NewVector2(x+width, y)))
	if err == nil {
		intersections = append(intersections, point)
	}

	point, err = line.Intersection(mesh.NewLine(vector.NewVector2(x, y), vector.NewVector2(x, y+height)))
	if err == nil {
		intersections = append(intersections, point)
	}

	point, err = line.Intersection(mesh.NewLine(vector.NewVector2(x+width, y), vector.NewVector2(x+width, y+height)))
	if err == nil {
		intersections = append(intersections, point)
	}

	point, err = line.Intersection(mesh.NewLine(vector.NewVector2(x, y+height), vector.NewVector2(x+width, y+height)))
	if err == nil {
		intersections = append(intersections, point)
	}

	return intersections
}

func fill(width float64, height float64, shapes []mesh.Shape) ([]mesh.Polygon, error) {

	for _, shape := range shapes {
		if len(shape.GetPoints()) < 3 {
			return nil, errors.New("Can't make a polygon with less than 3 points")
		}
	}

	numOfPoints := 0
	pointsPrefixSum := make([]int, len(shapes))
	for i, shape := range shapes {
		pointsPrefixSum[i] = numOfPoints
		numOfPoints += len(shape.GetPoints())
	}

	flatPoints := make([][2]float64, numOfPoints)

	for shapeIndex, shape := range shapes {
		for pointIndex, point := range shape.GetPoints() {
			flatPoints[pointIndex+pointsPrefixSum[shapeIndex]][0] = point.X()
			flatPoints[pointIndex+pointsPrefixSum[shapeIndex]][1] = point.Y()
		}
	}

	segments := make([][2]int32, numOfPoints)

	for shapeIndex, shape := range shapes {
		for pointIndex := 0; pointIndex < len(shape.GetPoints()); pointIndex++ {
			i := pointIndex + pointsPrefixSum[shapeIndex]
			segments[i][0] = int32(i)
			segments[i][1] = int32(((pointIndex + 1) % len(shape.GetPoints())) + pointsPrefixSum[shapeIndex])
		}
	}

	// Hole represented by a point lying inside it
	var holes = make([][2]float64, len(shapes))
	for i := range shapes {
		holes[i][0] = 0
		holes[i][1] = 0.0
	}

	v, faces := triangle.ConstrainedDelaunay(flatPoints, segments, holes)

	betterPolys := make([]mesh.Polygon, len(faces))
	for i, face := range faces {
		ourVerts := make([]vector.Vector3, 0)
		ourVerts = append(ourVerts, vector.NewVector3(v[face[0]][0], 0, v[face[0]][1]))
		ourVerts = append(ourVerts, vector.NewVector3(v[face[1]][0], 0, v[face[1]][1]))
		ourVerts = append(ourVerts, vector.NewVector3(v[face[2]][0], 0, v[face[2]][1]))
		poly, _ := mesh.NewPolygon(ourVerts, ourVerts)
		betterPolys[i] = poly
	}
	return betterPolys, nil
}

func carve(width float64, height float64, shapes []mesh.Shape) ([]mesh.Polygon, error) {

	for _, shape := range shapes {
		if len(shape.GetPoints()) < 3 {
			return nil, errors.New("Can't make a polygon with less than 3 points")
		}
	}

	numOfPoints := 4
	pointsPrefixSum := make([]int, len(shapes))
	for i, shape := range shapes {
		pointsPrefixSum[i] = numOfPoints
		numOfPoints += len(shape.GetPoints())
	}

	flatPoints := make([][2]float64, numOfPoints)

	flatPoints[0] = [2]float64{0.0, 0.0}
	flatPoints[1] = [2]float64{0.0, height}
	flatPoints[2] = [2]float64{width, height}
	flatPoints[3] = [2]float64{width, 0.0}

	for shapeIndex, shape := range shapes {
		for pointIndex, point := range shape.GetPoints() {
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
		for pointIndex := 0; pointIndex < len(shape.GetPoints()); pointIndex++ {
			i := pointIndex + pointsPrefixSum[shapeIndex]
			segments[i][0] = int32(i)
			segments[i][1] = int32(((pointIndex + 1) % len(shape.GetPoints())) + pointsPrefixSum[shapeIndex])
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

	betterPolys := make([]mesh.Polygon, len(faces))
	for i, face := range faces {
		ourVerts := make([]vector.Vector3, 3)
		ourVerts[0] = vector.NewVector3(v[face[0]][0], 0, v[face[0]][1])
		ourVerts[1] = vector.NewVector3(v[face[1]][0], 0, v[face[1]][1])
		ourVerts[2] = vector.NewVector3(v[face[2]][0], 0, v[face[2]][1])
		poly, _ := mesh.NewPolygon(ourVerts, ourVerts)
		betterPolys[i] = poly
	}
	return betterPolys, nil
}

func main() {

	outerRadius := 1.2
	innerRadius := 1.0
	ringHeight := .8

	sides := 32
	polys := make([]mesh.Polygon, sides*8)

	numTimesForTextureToRepeat := 8

	angleIncrement := (1.0 / float64(sides)) * 2.0 * math.Pi
	for sideIndex := 0; sideIndex < sides; sideIndex++ {
		angle := angleIncrement * float64(sideIndex)
		angleNext := angleIncrement * (float64(sideIndex) + 1)

		bottomLeftUV := vector.NewVector2(math.Min(float64(sideIndex)/float64(sides/numTimesForTextureToRepeat), 1), 0)
		topLeftUV := vector.NewVector2(math.Min(float64(sideIndex)/float64(sides/numTimesForTextureToRepeat), 1), 1.0)
		topRightUV := vector.NewVector2(math.Min(float64(sideIndex+1)/float64(sides/numTimesForTextureToRepeat), 1), 1.0)
		bottomRightUV := vector.NewVector2(math.Min(float64(sideIndex+1)/float64(sides/numTimesForTextureToRepeat), 1), 0)

		// outer
		square, err := makeSquare(
			vector.NewVector3(math.Cos(angle)*outerRadius, 0, math.Sin(angle)*outerRadius),
			vector.NewVector3(math.Cos(angle)*outerRadius, ringHeight, math.Sin(angle)*outerRadius),
			vector.NewVector3(math.Cos(angleNext)*outerRadius, ringHeight, math.Sin(angleNext)*outerRadius),
			vector.NewVector3(math.Cos(angleNext)*outerRadius, 0, math.Sin(angleNext)*outerRadius),
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
			vector.NewVector3(math.Cos(angleNext)*innerRadius, 0, math.Sin(angleNext)*innerRadius),
			vector.NewVector3(math.Cos(angleNext)*innerRadius, ringHeight, math.Sin(angleNext)*innerRadius),
			vector.NewVector3(math.Cos(angle)*innerRadius, ringHeight, math.Sin(angle)*innerRadius),
			vector.NewVector3(math.Cos(angle)*innerRadius, 0, math.Sin(angle)*innerRadius),
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
			vector.NewVector3(math.Cos(angle)*outerRadius, ringHeight, math.Sin(angle)*outerRadius),
			vector.NewVector3(math.Cos(angle)*innerRadius, ringHeight, math.Sin(angle)*innerRadius),
			vector.NewVector3(math.Cos(angleNext)*innerRadius, ringHeight, math.Sin(angleNext)*innerRadius),
			vector.NewVector3(math.Cos(angleNext)*outerRadius, ringHeight, math.Sin(angleNext)*outerRadius),
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
			vector.NewVector3(math.Cos(angle)*innerRadius, 0, math.Sin(angle)*innerRadius),
			vector.NewVector3(math.Cos(angle)*outerRadius, 0, math.Sin(angle)*outerRadius),
			vector.NewVector3(math.Cos(angleNext)*outerRadius, 0, math.Sin(angleNext)*outerRadius),
			vector.NewVector3(math.Cos(angleNext)*innerRadius, 0, math.Sin(angleNext)*innerRadius),
			bottomLeftUV,
			topLeftUV,
			topRightUV,
			bottomRightUV,
		)

		check(err)
		polys[(sideIndex*8)+6] = square[0]
		polys[(sideIndex*8)+7] = square[1]
	}

	fontByteData, err := ioutil.ReadFile("./sample.ttf")
	check(err)
	parsedFont, err := truetype.Parse(fontByteData)
	check(err)

	textToEnscribe := "Twitter"

	finalWord, err := mesh.NewModel([]mesh.Polygon{})
	check(err)

	for _, char := range textToEnscribe {
		log.Println(string(char))

		glyph := truetype.GlyphBuf{}
		glyph.Load(parsedFont, 100, parsedFont.Index(char), font.HintingNone)

		letterPoints := make([]vector.Vector2, len(glyph.Points))
		for i, p := range glyph.Points {
			letterPoints[i] = vector.NewVector2(float64(p.X), float64(p.Y)+10)
		}

		shape, err := mesh.NewShape(letterPoints)
		if err != nil {
			continue
		}

		shape.Scale(.1)

		bottomLeftBounds, topRightBounds := shape.GetBounds()
		width := (topRightBounds.X() - bottomLeftBounds.X())

		left, right := shape.Split((width / 2) + bottomLeftBounds.X())

		for r := range right {
			right[r].Translate(vector.NewVector2(-(width / 2), 0))
		}

		if len(left) > 0 {
			carved, err := carve(width/2, 10.0, left)
			check(err)
			lModel, err := mesh.NewModel(carved)
			check(err)
			finalWord = finalWord.Merge(lModel) //.Rotate(NewVector3(0, 0, -.2), NewVector3(0, 0, 0)).Translate(NewVector3(-width/2, 1, 0))
		}

		if len(right) > 0 {
			carved, err := carve((width/2)+0.5, 10.0, right)
			check(err)
			rModel, err := mesh.NewModel(carved)
			check(err)
			rModel = rModel.Translate(vector.NewVector3((width/2)+0, 0, 0))
			finalWord = finalWord.Merge(rModel) //.Rotate(NewVector3(0, 0, -.2), NewVector3(0, 0, 0)).Translate(NewVector3(-width/2, 1, 0))
		}

		finalWord = finalWord.Rotate(vector.NewVector3(0, 0, -.2), vector.NewVector3(0, 0, 0)).Translate(vector.NewVector3(-width, 1, 0))
	}

	log.Println("completed")

	f, err := os.Create("out.obj")
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	finalWord.Save(w)
	w.Flush()

}
