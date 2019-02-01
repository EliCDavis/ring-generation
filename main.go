package main

import (
	"bufio"
	"math"
	"os"
	"sort"
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

func subCarve(x float64, y float64, width float64, height float64, bounds []Vector2) []*Polygon {

	// If there exists more than one point in our region, subdivide
	var pointInRegion *Vector2
	for _, point := range bounds {
		if pointWithinBounds(x, y, width, height, point) {
			if pointInRegion != nil {
				return subDivide(x, y, width, height, bounds)
			}
			pointInRegion = &point
		}
	}

	// If there exists more than two lines in our region, subdivide
	allIntersections := make([]*Vector2, 0)
	for i := 1; i < len(bounds); i++ {
		line := NewLine(&bounds[i-1], &bounds[i])
		intersctions := lineIntersectsRectangle(x, y, width, height, line)
		if len(intersctions) > 0 {
			allIntersections = append(allIntersections, intersctions...)
			if len(allIntersections) > 2 {
				return subDivide(x, y, width, height, bounds)
			}
		}
	}

	// If no intersections or points, determine whether or not if this region is inside or outside the bounds
	if pointInRegion == nil && len(allIntersections) == 0 {
		if isInside(bounds, NewVector2(x, y)) == false {
			return makeSquareFrom2D(NewVector2(x, y), NewVector2(x+width, y+height))
		}
		return nil
	}

	// Else determine what corners of the rectangle we have access too
	pointsToWorkWith := make([]*Vector2, 0)
	if pointInRegion != nil {
		pointsToWorkWith = append(pointsToWorkWith, pointInRegion)
	}
	if len(allIntersections) > 0 {
		pointsToWorkWith = append(pointsToWorkWith, allIntersections...)
	}
	if isInside(bounds, NewVector2(x, y)) == false {
		pointsToWorkWith = append(pointsToWorkWith, NewVector2(x, y))
	}
	if isInside(bounds, NewVector2(x+width, y)) == false {
		pointsToWorkWith = append(pointsToWorkWith, NewVector2(x+width, y))
	}
	if isInside(bounds, NewVector2(x, y+height)) == false {
		pointsToWorkWith = append(pointsToWorkWith, NewVector2(x, y+height))
	}
	if isInside(bounds, NewVector2(x+width, y+height)) == false {
		pointsToWorkWith = append(pointsToWorkWith, NewVector2(x+width, y+height))
	}

	resultingShape := NewShape(pointsToWorkWith)
	sort.Sort(resultingShape)
	switch len(pointsToWorkWith) {
	case 3:
		return []*Polygon{NewPolygonFromShape(resultingShape)}

	case 4:
		return []*Polygon{
			NewPolygonFromFlatPoints([]*Vector2{resultingShape.points[0], resultingShape.points[1], resultingShape.points[2]}),
			NewPolygonFromFlatPoints([]*Vector2{resultingShape.points[0], resultingShape.points[2], resultingShape.points[3]}),
		}

	case 5:
		closestPoint := resultingShape.PointClosestToCenter()
		resultingPolygons := make([]*Polygon, 0)
		for i := 1; i < len(resultingShape.points)+1; i++ {
			n := i % len(resultingShape.points)
			if n != closestPoint && i-1 != closestPoint {
				if n > i {
					resultingPolygons = append(
						resultingPolygons,
						NewPolygonFromFlatPoints([]*Vector2{resultingShape.points[i-1], resultingShape.points[n], resultingShape.points[closestPoint]}))
				} else {
					resultingPolygons = append(
						resultingPolygons,
						NewPolygonFromFlatPoints([]*Vector2{resultingShape.points[n], resultingShape.points[closestPoint], resultingShape.points[i-1]}))
				}
			}
		}
		return resultingPolygons

	case 6:
		return subDivide(x, y, width, height, bounds)
	}
	panic("THERE ARE NOT THE RIGHT NUMBER OF POINTS: " + string(len(pointsToWorkWith)))
}

func subDivide(x float64, y float64, width float64, height float64, bounds []Vector2) []*Polygon {
	widthHalved := width / 2.0
	heightHalved := height / 2.0
	polys := make([]*Polygon, 0)
	polys = append(polys, subCarve(x, y, widthHalved, heightHalved, bounds)...)
	polys = append(polys, subCarve(x+widthHalved, y, widthHalved, heightHalved, bounds)...)
	polys = append(polys, subCarve(x, y+heightHalved, widthHalved, heightHalved, bounds)...)
	polys = append(polys, subCarve(x+widthHalved, y+heightHalved, widthHalved, heightHalved, bounds)...)
	return polys
}

func carve(width float64, height float64, bounds []Vector2) []*Polygon {
	return subDivide(0, 0, width, height, bounds)
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

	ahhhhh := carve(1.0, 1.0, []Vector2{
		Vector2{0.1, 0.1},
		Vector2{0.1, 0.9},
		Vector2{0.9, 0.9},
		Vector2{0.9, 0.1},
	})

	model, err := NewModel(ahhhhh)
	check(err)

	f, err := os.Create("out.obj")
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	model.Save(w)
	w.Flush()

}
