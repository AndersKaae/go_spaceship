package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Planet struct {
	name              string
	radius            float32
	texture           rl.Texture2D
	speed             float32
	textureOffset     float32
	yPos              float32 // Vertical position of the planet
	inSpace           bool
	renderTex         rl.RenderTexture2D
	initialized       bool
	lastRadius        float32 // Last radius used for rendering
	nameIndex         int
	potentialNames    []string
	potentialTextures []rl.Texture2D
}

func initSplashPlanet(planet *Planet, textures textures) *Planet {
	planetTextures := []rl.Texture2D{
		textures.jupiter,
		textures.mars,
		textures.neptune,
		textures.uranus,
		textures.venus,
		textures.luna,
		textures.earth,
		textures.saturn,
		textures.sun,
	}

	planetNames := []string{
		"Jupiter",
		"Mars",
		"Neptun",
		"Uranus",
		"Venus",
		"Månen",
		"Jorden",
		"Saturn",
		"Solen",
	}
	r := float32(1)

	planetIndex := int(rl.GetRandomValue(0, 8))

	splashPlanet := &Planet{
		radius:            0,
		texture:           planetTextures[planetIndex],
		name:              planetNames[planetIndex],
		speed:             0.4,
		yPos:              0 + r,
		inSpace:           false,
		nameIndex:         planetIndex,
		potentialNames:    planetNames,
		potentialTextures: planetTextures,
	}
	return splashPlanet
}

func (p *Planet) Draw(shader rl.Shader, altitude int32, speed float32) {
	if p.inSpace && altitude < planetHeight {
		return // Don't draw planets if altitude is too low
	}

	if p.yPos < float32(screenHeight/2)-p.radius {
		p.yPos += speed * 0.02 // Adjust vertical position based on speed change back to 0.02
	}

	renderSize := int32(p.radius * 2)
	if renderSize < 4 {
		return // Skip drawing if too small
	}

	if !p.initialized || p.lastRadius != p.radius {
		if p.initialized {
			rl.UnloadRenderTexture(p.renderTex)
		}
		p.renderTex = rl.LoadRenderTexture(renderSize, renderSize)
		p.initialized = true
		p.lastRadius = p.radius
	}

	// Simulate horizontal scrolling of texture (planet rotation)
	p.textureOffset += p.speed // Adjust this speed as needed
	if p.textureOffset >= float32(p.texture.Width) {
		p.textureOffset -= float32(p.texture.Width)
	}

	// Scale texture to fit the desired radius
	scale := p.radius / (float32(p.texture.Height) / 2)

	// --- DRAW TO OFFSCREEN TEXTURE ---
	rl.BeginTextureMode(p.renderTex)
	rl.ClearBackground(rl.Blank)

	// First visible portion of texture (right side)
	visibleWidth := float32(p.texture.Width) - p.textureOffset
	src1 := rl.NewRectangle(p.textureOffset, 0, visibleWidth, float32(p.texture.Height))
	dst1 := rl.NewRectangle(0, 0, visibleWidth*scale, float32(p.texture.Height)*scale)
	rl.DrawTexturePro(p.texture, src1, dst1, rl.NewVector2(0, 0), 0, rl.White)

	// Wrapped part (left side) if needed
	if p.textureOffset > 0 {
		src2 := rl.NewRectangle(0, 0, p.textureOffset, float32(p.texture.Height))
		dst2 := rl.NewRectangle(dst1.Width, 0, p.textureOffset*scale, float32(p.texture.Height)*scale)
		rl.DrawTexturePro(p.texture, src2, dst2, rl.NewVector2(0, 0), 0, rl.White)
	}

	rl.EndTextureMode()

	// --- APPLY SHADER AND DRAW TO SCREEN ---
	center := []float32{float32(p.renderTex.Texture.Width) / 2, float32(p.renderTex.Texture.Height) / 2}
	radiusUniform := []float32{p.radius}
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "center"), center, rl.ShaderUniformVec2)
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "radius"), radiusUniform, rl.ShaderUniformFloat)

	rl.BeginShaderMode(shader)

	screenPos := rl.NewVector2(
		float32(screenWidth)/2-p.radius,
		p.yPos,
	)

	// Note: flip vertically because textures from Framebuffer are upside down in raylib
	rl.DrawTextureRec(
		p.renderTex.Texture,
		rl.NewRectangle(0, 0, float32(p.renderTex.Texture.Width), -float32(p.renderTex.Texture.Height)),
		screenPos,
		rl.White,
	)

	rl.EndShaderMode()
}
