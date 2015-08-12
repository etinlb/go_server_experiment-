package main

import "fmt"

var gameId = 1

// TODO: Learn go better so these and the messages structs could be combined
// Might have to structure the json data begin sent differently
type GameObject interface {
	Update()
}

type Mover interface {
	Move(xVel, yVel int)
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type BaseGameObjData struct {
	Id string `json:"id"`
}

type MovableObject struct {
	BaseGameObjData
	Rect
	XVel int `json:"xVel"`
	YVel int `json:"yVel"`
}

func (m *MovableObject) Update() {
	// yeah
	fmt.Printf("%+v IN Update\n", m)
	m.X += m.XVel
	m.Y += m.YVel
	fmt.Printf("%+v After moving\n", m)

}

func (m *MovableObject) Move(xVel, yVel int) {
	// yeah
	fmt.Printf("%+v IN MOVe\n", m)
	m.XVel = xVel
	m.YVel = yVel
	fmt.Printf("%+v After Moving\n", m)
}

func NewGameObject(x, y int, id string) MovableObject {
	// rect := MakeRect(x, y)
	gameObject := MovableObject{BaseGameObjData: BaseGameObjData{Id: id}, Rect: Rect{X: x, Y: y}, XVel: 0, YVel: 0}

	return gameObject
}

func MakeRect(x, y int) Rect {
	rect := Rect{X: x, Y: y}

	return rect
}
