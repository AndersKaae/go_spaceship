package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Spaceship struct {
	spaceshipTexture rl.Texture2D
	textureWidth     int32
	textureHeight    int32
	centerX          int32
	x, y             int32
	speed            float32
	dustClouds       []DustCloud
	needDustClouds   bool
	launched         bool
	inSpace          bool
}

type DustCloud struct {
	x, y         int32
	speed        float32
	textureWidth int32
	centerX      int32
	alpha        float32 // Opacity of the dust cloud
}

func NewSpaceship(spaceshipTexture rl.Texture2D) *Spaceship {
	x := screenWidth/2 - spaceshipTexture.Width/2
	y := screenHeight - spaceshipTexture.Height

	return &Spaceship{
		textureWidth:   spaceshipTexture.Width,
		textureHeight:  spaceshipTexture.Height,
		x:              x,
		y:              y,
		centerX:        x + spaceshipTexture.Width/2, // <-- now x is known
		speed:          0,
		dustClouds:     []DustCloud{},
		launched:       false,
		needDustClouds: true,
		inSpace:        false,
	}
}

func (s *Spaceship) SpawnDustClouds(cloudTexture rl.Texture2D) {
	// Now s is a Spaceship, which has dustClouds slice

	if len(s.dustClouds) < 2 {
		x := s.centerX - 50
		leftCloud := DustCloud{
			x:            x,
			y:            s.y + 100,
			speed:        1,
			textureWidth: cloudTexture.Width,
			centerX:      x - cloudTexture.Width/2,
			alpha:        1.0, // Set initial opacity to 1.0
		}
		x = s.centerX + 50 - cloudTexture.Width
		rightCloud := DustCloud{
			x:            x,
			y:            s.y + 100,
			speed:        1,
			textureWidth: cloudTexture.Width,
			centerX:      x + cloudTexture.Width/2,
			alpha:        1.0, // Set initial opacity to 1.0
		}
		s.dustClouds = append(s.dustClouds, leftCloud, rightCloud)
	}
	s.needDustClouds = false
}

func (s *Spaceship) DrawDustClouds(texture rl.Texture2D) {
	for _, cloud := range s.dustClouds {
		rl.DrawTexture(texture, cloud.x, cloud.y, rl.Fade(rl.White, cloud.alpha))
	}
}

func (s *Spaceship) UpdateDustClouds() {
	for i := 0; i < len(s.dustClouds); i++ {
		// determine if the dust cloud is left of the spaceship
		if s.dustClouds[i].x < s.centerX {
			// move it to the left
			s.dustClouds[i].x -= int32(s.dustClouds[i].speed)
		} else {
			// move it to the right
			s.dustClouds[i].x += int32(s.dustClouds[i].speed)
		}
		// move the dust cloud down
		s.dustClouds[i].y += int32(s.dustClouds[i].speed) / 2
		// reduce the opacity of the dust cloud
		s.dustClouds[i].alpha -= 0.01
		// remove the dust cloud if its opacity is less than 0
		if s.dustClouds[i].alpha <= 0 {
			s.dustClouds = append(s.dustClouds[:i], s.dustClouds[i+1:]...)
			i-- // Adjust index after removal
		}
	}
}

func (s *Spaceship) Accelerate() {
	if s.speed < maxSpeed {
		s.speed += 1.05
	}
}
