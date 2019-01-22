package main

import (
	"errors"
	"io"
)

// Vertex is a position and a normal
type Vertex struct {
	position *Vector3
	normal   *Vector3
}

// NewVertex builds a new vertex
func NewVertex(position *Vector3, normal *Vector3) (*Vertex, error) {
	if position == nil {
		return nil, errors.New("Vertex must have a position")
	}

	if normal == nil {
		return nil, errors.New("Vertex must have a normal")
	}

	return &Vertex{position, normal}, nil
}

// Save Writes a model to obj format
func (m Vertex) Save(w io.Writer) error {

	return nil
}
