package main

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/EliCDavis/vector"
)

// Polygon represents a single polygon made up of multiple points
type Polygon struct {
	vertices []*vector.Vector3
	normals  []*vector.Vector3
	uv       []*vector.Vector2
}

// NewPolygon creates a new polygon
func NewPolygon(vertices []*vector.Vector3, normals []*vector.Vector3) (*Polygon, error) {
	if vertices == nil {
		return nil, errors.New("Must provide vertices")
	}

	if normals == nil {
		return nil, errors.New("Must provide normals")
	}

	if len(vertices) < 3 {
		return nil, errors.New("Polygon must have 3 or more points")
	}

	if len(normals) < 3 {
		return nil, errors.New("Polygon must have 3 or more normals")
	}

	if len(normals) != len(vertices) {
		return nil, errors.New("The number of vertices and normals must match")
	}

	return &Polygon{vertices, normals, nil}, nil
}

func (p Polygon) Translate(mv *vector.Vector3) *Polygon {
	newVertices := make([]*vector.Vector3, len(p.vertices))
	for pIndex := range p.vertices {
		newVertices[pIndex] = p.vertices[pIndex].Add(mv)
	}
	return &Polygon{newVertices, newVertices, p.uv}
}

func (p Polygon) Rotate(amount *vector.Vector3, pivot *vector.Vector3) *Polygon {
	newVertices := make([]*vector.Vector3, len(p.vertices))

	for pIndex, point := range p.vertices {

		// https://play.golang.org/p/qWUotd3Lb56
		final := point.Sub(pivot)

		// Pretty sure is correct
		zLength := math.Sqrt(math.Pow(final.X(), 2.0) + math.Pow(final.Y(), 2.0))
		if zLength > 0 {
			zRot := math.Atan2(final.Y(), final.X()) + amount.Z()
			final = vector.NewVector3(
				math.Cos(zRot)*zLength,
				math.Sin(zRot)*zLength,
				final.Z(),
			)
		}

		// Not sure
		// yLength := math.Sqrt(math.Pow(final.x, 2.0) + math.Pow(final.z, 2.0))
		// if yLength > 0 {
		// 	yRot := math.Atan(final.z/final.x) + amount.y
		// 	final = NewVector3(
		// 		math.Cos(yRot)*yLength,
		// 		final.y,
		// 		math.Sin(yRot)*yLength,
		// 	)
		// }

		// Not sure
		// xLength := math.Sqrt(math.Pow(final.z, 2.0) + math.Pow(final.y, 2.0))
		// if xLength > 0 {
		// 	xRot := math.Atan(final.z/final.y) + amount.x
		// 	final = NewVector3(
		// 		final.x,
		// 		math.Cos(xRot)*xLength,
		// 		math.Sin(xRot)*xLength,
		// 	)
		// }

		newVertices[pIndex] = final.Add(pivot)
	}
	return &Polygon{newVertices, newVertices, p.uv}
}

// NewPolygonFromShape creates a 3D polygon from a 2D shape
func NewPolygonFromShape(shape *Shape) *Polygon {
	vertices := make([]*vector.Vector3, len(shape.points))
	for i, point := range shape.points {
		vertices[i] = vector.NewVector3(point.X(), 0, point.Y())
	}
	poly, _ := NewPolygon(vertices, vertices)
	return poly
}

// NewPolygonFromFlatPoints creates a polygon from 2d points
func NewPolygonFromFlatPoints(points []*vector.Vector2) *Polygon {
	vertices := make([]*vector.Vector3, len(points))
	for i, point := range points {
		vertices[i] = vector.NewVector3(point.X(), 0, point.Y())
	}
	poly, _ := NewPolygon(vertices, vertices)
	return poly
}

// NewPolygonWithTexture creates a polygon with uv coordinates
func NewPolygonWithTexture(vertices []*vector.Vector3, normals []*vector.Vector3, texture []*vector.Vector2) (*Polygon, error) {
	poly, err := NewPolygon(vertices, normals)
	if err != nil {
		return nil, err
	}

	if texture == nil {
		return nil, errors.New("Must provide texture")
	}

	if len(texture) != len(vertices) {
		return nil, errors.New("Texture length must match vertices")
	}

	poly.uv = texture
	return poly, nil
}

// Save Writes a polygon to obj format and returns the number of
func (p Polygon) Save(w io.Writer, pointOffset int) (int, error) {

	face := "f "

	for pointIndex := 0; pointIndex < len(p.vertices); pointIndex++ {
		_, err := w.Write([]byte(fmt.Sprintf("v %f %f %f\n", p.vertices[pointIndex].X(), p.vertices[pointIndex].Y(), p.vertices[pointIndex].Z())))
		if err != nil {
			return 0, err
		}

		_, err = w.Write([]byte(fmt.Sprintf("vn %f %f %f\n", p.normals[pointIndex].X(), p.normals[pointIndex].Y(), p.normals[pointIndex].Z())))
		if err != nil {
			return 0, err
		}

		if p.uv != nil {
			_, err = w.Write([]byte(fmt.Sprintf("vt %f %f \n", p.uv[pointIndex].X(), p.uv[pointIndex].Y())))
			if err != nil {
				return 0, err
			}
			face += fmt.Sprintf("%d/%d ", pointIndex+pointOffset, pointIndex+pointOffset)
		} else {
			face += fmt.Sprintf("%d ", pointIndex+pointOffset)
		}

	}

	_, err := w.Write([]byte(face + "\n"))
	if err != nil {
		return 0, err
	}

	return pointOffset + len(p.vertices), nil
}
