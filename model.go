package main

import (
	"errors"
	"io"
)

// Model is built with a collection of polygons
type Model struct {
	faces []*Polygon
}

// NewModel builds a new model
func NewModel(faces []*Polygon) (*Model, error) {
	if faces == nil {
		return nil, errors.New("Can not have nil faces")
	}
	return &Model{faces}, nil
}

func (m Model) Merge(other *Model) *Model {
	return &Model{append(m.faces, other.faces...)}
}

func (m Model) Translate(movement *Vector3) *Model {
	newFaces := make([]*Polygon, len(m.faces))
	for f := range m.faces {
		newFaces[f] = m.faces[f].Translate(movement)
	}
	return &Model{newFaces}
}

func (p Model) Rotate(amount *Vector3, pivot *Vector3) *Model {
	newVertices := make([]*Polygon, len(p.faces))
	for pIndex, point := range p.faces {
		newVertices[pIndex] = point.Rotate(amount, pivot)
	}
	return &Model{newVertices}
}

// Save Writes a model to obj format
func (m Model) Save(w io.Writer) error {

	w.Write([]byte("mtllib master.mtl\n"))
	w.Write([]byte("usemtl wood\n"))

	offset := 1
	var err error
	for _, face := range m.faces {
		offset, err = face.Save(w, offset)
		if err != nil {
			return err
		}
	}

	return nil
}
