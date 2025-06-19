package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"path/filepath"
	"strconv"
)

const (
	screenWidth      = int32(450)
	screenHeight     = int32(800)
	maxSpeed         = 20
	terminalVelocity = -20
	skyR             = 189
	skyG             = 245
	skyB             = 255
	starSpeed        = 1.0 // Speed at which stars move downwards
	maxBisser        = 100
	spaceFadeStart   = 1000 // Altitude range for space transition start
	spaceFadeEnd     = 4000 // Altitude range for space transition
	planetHeight     = 6000 // Height at which planets are drawn
)

var (
	frameCounter   int32 = 0
	drawnSpeed     float32
	skyColor       = Sky{skyR, skyG, skyB, 255}
	stars          []rl.Vector2
	altitude       int32 = 0
	skyObjectsList []SkyItem
	dustCloudsList []DustCloud
)

type Sky struct {
	r, g, b, a uint8
}

type textures struct {
	spaceship    rl.Texture2D
	spaceshipOff rl.Texture2D
	dustCloud    rl.Texture2D
	cloud        rl.Texture2D
	bissen       rl.Texture2D
	jupiter      rl.Texture2D
	mars         rl.Texture2D
	neptune      rl.Texture2D
	uranus       rl.Texture2D
	venus        rl.Texture2D
	luna         rl.Texture2D
	earth        rl.Texture2D
	saturn       rl.Texture2D
	sun          rl.Texture2D
}

func loadTextures() textures {
	return textures{
		spaceship:    createTextureFromImage("/graphics/spaceship.png", 0.05),
		spaceshipOff: createTextureFromImage("/graphics/spaceshipOff.png", 0.05),
		dustCloud:    createTextureFromImage("/graphics/dustCloud.png", 0.1),
		cloud:        createTextureFromImage("/graphics/cloud.png", 0.1),
		bissen:       createTextureFromImage("/graphics/bissen.png", 0.12),
		jupiter:      createTextureFromImage("/graphics/jupiterTexture.jpg", 1),
		mars:         createTextureFromImage("/graphics/marsTexture.jpg", 1),
		neptune:      createTextureFromImage("/graphics/neptuneTexture.jpg", 1),
		uranus:       createTextureFromImage("/graphics/uranusTexture.jpg", 1),
		venus:        createTextureFromImage("/graphics/venusTexture.jpg", 1),
		luna:         createTextureFromImage("/graphics/lunaTexture.jpg", 1),
		earth:        createTextureFromImage("/graphics/earthTexture.png", 1),
		saturn:       createTextureFromImage("/graphics/saturnTexture.png", 1),
		sun:          createTextureFromImage("/graphics/sunTexture.jpg", 1),
	}
}

type shaders struct {
	cicleMask rl.Shader
}

func loadShaders() shaders {
	return shaders{
		cicleMask: rl.LoadShader("", "shaders/circle_mask.fs"),
	}
}

func transitionSkyToSpace(sky Sky, altitude int32, playerSpaceship *Spaceship) Sky {

	// Normalize altitude between 0.0 and 1.0, starting from 500
	normalizedAltitude := float32(altitude-spaceFadeStart) / float32(spaceFadeEnd-spaceFadeStart)

	if normalizedAltitude > 1.0 {
		normalizedAltitude = 1.0
		playerSpaceship.inSpace = true
	}
	if normalizedAltitude < 0.0 {
		normalizedAltitude = 0.0
		playerSpaceship.inSpace = false
	}

	// Interpolate sky color
	sky.r = uint8(float32(skyR) * (1 - normalizedAltitude))
	sky.g = uint8(float32(skyG) * (1 - normalizedAltitude))
	sky.b = uint8(float32(skyB) * (1 - normalizedAltitude))

	return sky
}

func checkIfSpace(sky Sky) bool {
	var limit uint8 = 100
	if sky.r < limit && sky.g < limit && sky.b < limit {
		return true
	}
	return false
}

func findObjectInList(name string) *SkyItem {
	for i := range skyObjectsList {
		if skyObjectsList[i].name == name {
			return &skyObjectsList[i]
		}
	}
	return nil
}

func generateStars() {
	for i := 0; i < 100; i++ { // Adjust the number of stars as needed
		x := float32(rl.GetRandomValue(0, int32(screenWidth)))
		y := float32(rl.GetRandomValue(0, int32(screenHeight)))
		stars = append(stars, rl.Vector2{X: x, Y: y})
	}
}

func updateStars(speed float32) {
	adjustedSpeed := starSpeed
	if speed < 0 {
		adjustedSpeed = -starSpeed
	}
	for i := range stars {
		stars[i].Y += float32(adjustedSpeed)
		if stars[i].Y > float32(screenHeight) {
			stars[i].Y = 0
			stars[i].X = float32(rl.GetRandomValue(0, int32(screenWidth)))
		} else if stars[i].Y < 0 {
			stars[i].Y = float32(screenHeight)
			stars[i].X = float32(rl.GetRandomValue(0, int32(screenWidth)))
		}
	}
}

func renderStars() {
	for _, star := range stars {
		rl.DrawPixelV(star, rl.White)
	}
}

func createTextureFromImage(imagePath string, scale float32) rl.Texture2D {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current working directory. Exiting...")
		rl.CloseWindow()
		panic(err)
	}

	// Combine the working directory with the image path
	fullPath := filepath.Join(cwd, imagePath)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Println("File does not exist:", fullPath)
		rl.CloseWindow()
		panic("File does not exist: " + fullPath)
	}

	// Load the image
	image := rl.LoadImage(fullPath)

	// Check if the image was loaded correctly
	if image.Data == nil {
		fmt.Println("Failed to load image from file:", fullPath)
		rl.CloseWindow()
		panic("Failed to load image " + fullPath + " from file. Exiting...")
	}

	// Resize the image
	rl.ImageResize(image, int32(float32(image.Width)*scale), int32(float32(image.Height)*scale))

	// Convert the image to a texture
	texture := rl.LoadTextureFromImage(image)

	// Unload image from CPU memory after creating the texture
	rl.UnloadImage(image)

	return texture
}

func update(listOfBisser *[]bissen, textures textures, playerSpaceship *Spaceship, splashScreen *bool) {
	if rl.IsKeyPressed(rl.KeyR) && !playerSpaceship.launched {
		if len(*listOfBisser) < maxBisser {
			*listOfBisser = append(*listOfBisser, *NewBissen(textures))
		}
	}

	// First, update all bisser
	for i := range *listOfBisser {
		(*listOfBisser)[i].Update()
	}

	// Now, remove invisible bisser
	newList := (*listOfBisser)[:0]
	for _, b := range *listOfBisser {
		if b.visible {
			newList = append(newList, b)
		}
	}
	*listOfBisser = newList
}

func main() {
	// Initialize the window
	rl.InitWindow(screenWidth, screenHeight, "Otto i Rummet")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	textures := loadTextures()
	shaders := loadShaders()

	// Create the spaceship struct
	playerSpaceship := NewSpaceship(textures.spaceshipOff)

	// Instantiate the bissen
	listOfBisser := []bissen{}
	listOfBisser = append(listOfBisser, *NewBissen(textures))

	// Generate stars once
	generateStars()

	// Create planet
	r := float32(200)
	gamePlanet := &Planet{
		radius:  r,
		texture: textures.mars,
		speed:   0.4,
		yPos:    0 - (2 * r),
		inSpace: true,
	}

	r = float32(1)
	splashPlanet := &Planet{
		radius:  0,
		texture: textures.mars,
		speed:   0.4,
		yPos:    0 + r,
		inSpace: false,
	}

	// If we increase  the number of planets, we need to change the planetIndex
	planetIndex := int(rl.GetRandomValue(0, 8))

	ground := NewGround()

	splashScreen := true

	radiusProgress := 0.0

	for !rl.WindowShouldClose() {
		if splashScreen {
			updateSplashscreen(splashPlanet, &planetIndex, textures, &radiusProgress)
			showSplashScreen(splashPlanet, shaders, playerSpaceship, &splashScreen)
			continue
		}
		gamePlanet.texture = splashPlanet.texture
		updateStars(playerSpaceship.speed)
		update(&listOfBisser, textures, playerSpaceship, &splashScreen)
		renderScene(textures, shaders, playerSpaceship, listOfBisser, ground, gamePlanet)
	}
}

func showSplashScreen(planet *Planet, shaders shaders, playerSpaceship *Spaceship, splashScreen *bool) {
	rl.BeginDrawing()
	renderStars()
	rl.ClearBackground(rl.Black)
	planet.Draw(shaders.cicleMask, altitude, playerSpaceship.speed)
	centerTekst("Otto i Rummet", screenHeight-700)
	centerTekst(planet.name, screenHeight-200)
	centerTekst("Prese space key", screenHeight-100)
	if rl.IsKeyPressed(rl.KeySpace) {
		*splashScreen = false
	}
	rl.EndDrawing()
}

func centerTekst(text string, y int32) {
	fontSize := int32(40)
	textWidth := rl.MeasureText(text, fontSize)
	x := (screenWidth - textWidth) / 2
	rl.DrawText(text, x, y, fontSize, rl.White)
}

func updateSplashscreen(planet *Planet, planetIndex *int, textures textures, progress *float64) {
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

	planet.name = planetNames[*planetIndex]

	var changed bool

	if rl.IsKeyPressed(rl.KeyA) {
		*planetIndex--
		changed = true
		if *planetIndex < 0 {
			*planetIndex = len(planetTextures) - 1
		}
	}
	if rl.IsKeyPressed(rl.KeyD) {
		*planetIndex++
		changed = true
		if *planetIndex >= len(planetTextures) {
			*planetIndex = 0
		}
	}
	if changed {
		planet.texture = planetTextures[*planetIndex]
		planet.radius = 1.0
		*progress = 0.0
		planet.initialized = false // force re-init render texture
	}

	// Update radiusProgress toward 1
	if planet.radius < 199 {
		*progress = approachOne(*progress)
		planet.radius = 200 * float32(*progress)
		planet.yPos = float32(screenHeight)/2 - planet.radius
		println("Progress:", planet.radius, "Planet Index:", *planetIndex)
	}
}

func approachOne(current float64) float64 {
	approachRate := 0.01 // Adjust this value to control the speed of approach
	return current + (1-current)*approachRate
}

func drawSpeed(speed float32, frameCounter int32, skyColor Sky, textY int32) {
	// Only update the speed every second
	if frameCounter%5 == 0 {
		drawnSpeed = speed
	}
	if drawnSpeed < 0 {
		drawnSpeed = drawnSpeed * -1
	}
	speedStr := strconv.Itoa(int(drawnSpeed))
	fontSize := int32(20)
	textWidth := rl.MeasureText(speedStr, fontSize)
	textX := screenWidth - textWidth - 10 // 10 pixels from the right edge
	if checkIfSpace(skyColor) == false {
		rl.DrawText(speedStr, textX, textY, fontSize, rl.Black)
	} else {
		rl.DrawText(speedStr, textX, textY, fontSize, rl.White)
	}
}

func drawSpaceshipHeight(altitude int32) {
	fontSize := int32(20)
	heightStr := strconv.Itoa(int(altitude))
	// 10 pixes from the left edge
	textX := int32(10)
	if checkIfSpace(skyColor) == false {
		rl.DrawText(heightStr, textX, 10, fontSize, rl.Black)
	} else {
		rl.DrawText(heightStr, textX, 10, fontSize, rl.White)
	}
}

func renderScene(textures textures, shaders shaders, playerSpaceship *Spaceship, listOfBisser []bissen, ground *ground, jupiter *Planet) {
	rl.BeginDrawing()

	backgroundColor := rl.NewColor(skyColor.r, skyColor.g, skyColor.b, skyColor.a)
	rl.ClearBackground(backgroundColor)

	// Count frames
	frameCounter++

	if checkIfSpace(skyColor) {
		renderStars()
	}

	spaceShipToDraw := textures.spaceshipOff
	if rl.IsKeyDown(rl.KeySpace) {
		spaceShipToDraw = textures.spaceship
		playerSpaceship.launched = true
		playerSpaceship.Accelerate()
	}

	if playerSpaceship.speed > terminalVelocity {
		playerSpaceship.speed -= 1
	}

	// Detect collision with the ground
	if playerSpaceship.y >= screenHeight-textures.spaceship.Height && playerSpaceship.speed < 0 {
		playerSpaceship.y = screenHeight - textures.spaceship.Height
		playerSpaceship.speed = 0
		playerSpaceship.launched = false
		playerSpaceship.needDustClouds = true
		altitude = 0
	}

	// This draws the green ground
	ground.draw(*playerSpaceship)

	// This draws the bissen
	for i := range listOfBisser {
		listOfBisser[i].Draw(textures)
	}

	jupiter.Draw(shaders.cicleMask, altitude, playerSpaceship.speed)

	cloud1 := findObjectInList("cloud1")
	if cloud1 == nil {
		cloud1 = NewSkyItem("cloud1", rl.Vector2{X: 50, Y: 50}, textures.cloud)
	}

	heightOverHalf := screenHeight/2 - textures.spaceship.Height/2

	if playerSpaceship.y > heightOverHalf {
		rl.DrawTexture(spaceShipToDraw, playerSpaceship.x, playerSpaceship.y, rl.White)
		cloud1.move(0)

	} else {
		rl.DrawTexture(spaceShipToDraw, playerSpaceship.x, screenHeight/2-textures.spaceship.Height/2, rl.White)
		cloud1.move(int32(playerSpaceship.speed))
	}

	skyColor = transitionSkyToSpace(skyColor, altitude, playerSpaceship)

	if playerSpaceship.launched && playerSpaceship.needDustClouds {
		playerSpaceship.SpawnDustClouds(textures.dustCloud)
	}

	playerSpaceship.UpdateDustClouds()
	playerSpaceship.DrawDustClouds(textures.dustCloud)

	if !playerSpaceship.inSpace {
		cloud1.draw()
	}

	playerSpaceship.y -= int32(playerSpaceship.speed)

	altitude = altitude + int32(playerSpaceship.speed)
	drawSpeed(playerSpaceship.speed, frameCounter, skyColor, 10)
	drawSpaceshipHeight(altitude)

	rl.EndDrawing()
}
