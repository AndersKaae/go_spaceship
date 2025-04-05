package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ground struct {
	groundColor  rl.Color
	groundWidth  int32
	groundHeight int32
}

func NewGround() *ground {
	groundWidth := int32(rl.GetScreenWidth())
	groundHeight := int32(rl.GetScreenHeight() - 70)
	groundColor := rl.NewColor(52, 235, 131, 255) // Custom color #34eb83

	return &ground{
		groundColor:  groundColor,
		groundWidth:  groundWidth,
		groundHeight: groundHeight,
	}
}

func (g *ground) draw(playerSpaceship Spaceship) {
	heightOverHalf := screenHeight/2 - playerSpaceship.textureHeight/2
	if playerSpaceship.y > heightOverHalf {
		rl.DrawRectangle(0, screenHeight-70, screenWidth, 70, g.groundColor)

	} else {
		rl.DrawRectangle(0, screenHeight-(70+(playerSpaceship.y-heightOverHalf)), screenWidth, 70, g.groundColor)
	}
}
