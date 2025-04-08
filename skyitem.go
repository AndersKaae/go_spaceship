package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SkyItem struct {
	name     string
	position rl.Vector2
	texture  rl.Texture2D
	speed    float32
}

func NewSkyItem(name string, position rl.Vector2, cloud rl.Texture2D) *SkyItem {
	skyItem := SkyItem{
		name:     name,
		position: position,
		texture:  cloud,
		speed:    0.5,
	}
	skyObjectsList = append(skyObjectsList, skyItem)
	return &skyItem
}

func (s *SkyItem) draw() {
	rl.DrawTexture(s.texture, int32(s.position.X), int32(s.position.Y), rl.White)
}

func (s *SkyItem) move(speed int32) {
	s.position.X += s.speed
	if s.position.X > float32(screenWidth) {
		s.position.X = -float32(s.texture.Width)
	} else if s.position.X < -float32(s.texture.Width) {
		s.position.X = float32(screenWidth)
	}

	s.position.Y += float32(speed)
	if s.position.Y > float32(screenHeight) {
		s.position.Y = -float32(s.texture.Height)
		s.position.X = float32(rl.GetRandomValue(0, int32(screenWidth)))
	} else if s.position.Y < -float32(s.texture.Height) {
		s.position.Y = float32(screenHeight)
		s.position.X = float32(rl.GetRandomValue(0, int32(screenWidth)))
	}
}
