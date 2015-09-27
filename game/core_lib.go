package main

type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func NewVector(x, y float64) Vector2D {
	rect := Vector2D{X: x, Y: y}

	return rect
}
