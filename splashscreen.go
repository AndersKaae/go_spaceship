package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func getAverageColor(tex rl.Texture2D) rl.Color {
	// Convert texture to image
	image := rl.LoadImageFromTexture(tex)
	defer rl.UnloadImage(image)

	// Get color data (one pixel per element)
	pixels := rl.LoadImageColors(image)
	defer rl.UnloadImageColors(pixels)

	var totalR, totalG, totalB, totalA int
	count := int(image.Width) * int(image.Height)

	for i := 0; i < count; i++ {
		totalR += int(pixels[i].R)
		totalG += int(pixels[i].G)
		totalB += int(pixels[i].B)
		totalA += int(pixels[i].A)
	}

	return rl.NewColor(
		uint8(totalR/count),
		uint8(totalG/count),
		uint8(totalB/count),
		uint8(totalA/count),
	)
}

func showSplashScreen(planet *Planet, shaders shaders, playerSpaceship *Spaceship, splashScreen *bool) {
	rl.BeginDrawing()
	renderStars()
	rl.ClearBackground(rl.Black)
	planet.Draw(shaders.cicleMask, altitude, playerSpaceship.speed)
	centerTekst("Otto i Rummet", screenHeight-700, rl.White)
	centerTekst(planet.name, screenHeight-180, planet.planetFontColor)
	centerTekst("Prese space key", screenHeight-100, rl.White)
	if rl.IsKeyPressed(rl.KeySpace) {
		*splashScreen = false
	}
	rl.EndDrawing()
}

func centerTekst(text string, y int32, color rl.Color) {
	fontSize := int32(40)
	textWidth := rl.MeasureText(text, fontSize)
	x := (screenWidth - textWidth) / 2
	rl.DrawText(text, x, y, fontSize, color)
}

func updateSplashscreen(planet *Planet, progress *float64) {
	planet.name = planet.potentialNames[planet.nameIndex]

	var changed bool

	if rl.IsKeyPressed(rl.KeyA) {
		planet.nameIndex--
		changed = true
		if planet.nameIndex < 0 {
			planet.nameIndex = len(planet.potentialTextures) - 1
		}
	}
	if rl.IsKeyPressed(rl.KeyD) {
		planet.nameIndex++
		changed = true
		if planet.nameIndex >= len(planet.potentialTextures) {
			planet.nameIndex = 0
		}
	}
	if changed {
		planet.texture = planet.potentialTextures[planet.nameIndex]
		planet.radius = 1.0
		planet.planetFontColor = getAverageColor(planet.texture)
		*progress = 0.0
		planet.initialized = false // force re-init render texture
	}

	// quite the game if 'Q' is pressed
	if rl.IsKeyPressed(rl.KeyQ) {
		rl.CloseWindow()
	}

	// Update radiusProgress toward 1
	if planet.radius < 199 {
		*progress = approachOne(*progress)
		if (*progress*200)-float64(planet.radius) > 1 { // Only update the sprite if the radius has changed significantly
			planet.radius = 200 * float32(*progress)
			planet.yPos = float32(screenHeight)/2 - planet.radius
		}
	}
}

func approachOne(current float64) float64 {
	approachRate := 0.01 // Adjust this value to control the speed of approach
	return current + (1-current)*approachRate
}
