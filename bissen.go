package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type bissen struct {
	x, y    int32
	visible bool
}

// instantiate the bissen
func NewBissen(textures textures) *bissen {
	return &bissen{
		x:       screenWidth - 10,
		y:       screenHeight - textures.bissen.Height - 20,
		visible: true,
	}
}

func (b *bissen) Update(listOfBisser []bissen, textures textures) {
	if rl.IsKeyPressed(rl.KeyR) {
		listOfBisser = append(listOfBisser, *NewBissen(textures))
		fmt.Print(len(listOfBisser))
	}
}

func (b *bissen) Draw(textures textures) {
	if b.visible {
		b.x -= 2
		rl.DrawTexture(textures.bissen, b.x, b.y, rl.White)
	}

	if b.x < screenWidth/2 {
		b.visible = false
	}
}
