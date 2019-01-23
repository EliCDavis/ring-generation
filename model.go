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
