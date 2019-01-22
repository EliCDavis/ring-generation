package main

import (
	"errors"
	"fmt"
	"io"
)

// Polygon represents a single polygon made up of multiple points
type Polygon struct {
	vertices []Vector3
	normals  []Vector3
}

// NewPolygon creates a new polygon
func NewPolygon(vertices []Vector3, normals []Vector3) (*Polygon, error) {
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

	return &Polygon{vertices, normals}, nil
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

		face += fmt.Sprintf("%d//%d ", pointIndex+pointOffset, pointIndex+pointOffset)
	}

	_, err := w.Write([]byte(face + "\n"))
	if err != nil {
		return 0, err
	}

	return pointOffset + len(p.vertices), nil
}
