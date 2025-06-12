/*
Keyboard Layout Editor (KLE) implementation
github.com/ijprest/kle-serial
*/
package main

type Keyboard struct {
	meta KeyboardMetadata
	keys []Key
}

type KeyboardMetadata struct {
	author    string
	backcolor string
	//background { name string style string } | null
	name        string
	notes       string
	radii       string
	switchBrand string
	switchMount string
	switchType  string
}

type Key struct {
	color  string
	labels []string
	//textColor Array<string | undefined>
	//textSize Array<number | undefined>
	//default { textColor string textSize number }

	x      uint8
	y      uint8
	width  uint8
	height uint8

	x2      uint8
	y2      uint8
	width2  uint8
	height2 uint8

	rotation_x     uint8
	rotation_y     uint8
	rotation_angle uint8

	decal   bool
	ghost   bool
	stepped bool
	nub     bool

	profile string

	sm string // switch mount
	sb string // switch brand
	st string // switch type
}
