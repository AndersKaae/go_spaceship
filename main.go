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
}

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

func update(b *bissen) {
	if rl.IsKeyPressed(rl.KeyR) {
		b.visible = true
		b.x = screenWidth - 10
	}
}

func loadTextures() textures {
	return textures{
		spaceship:    createTextureFromImage("/graphics/spaceship.png", 0.05),
		spaceshipOff: createTextureFromImage("/graphics/spaceshipOff.png", 0.05),
		dustCloud:    createTextureFromImage("/graphics/dustCloud.png", 0.1),
		cloud:        createTextureFromImage("/graphics/cloud.png", 0.1),
		bissen:       createTextureFromImage("/graphics/bissen.png", 0.12),
	}
}

func transitionSkyToSpace(sky Sky, altitude int32, playerSpaceship *Spaceship) Sky {
	const spaceFadeStart = 500
	const spaceFadeEnd = 2000

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

func main() {
	// Initialize the window
	rl.InitWindow(screenWidth, screenHeight, "Bissen til Saturn")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	textures := loadTextures()

	// Create the spaceship struct
	playerSpaceship := NewSpaceship(textures.spaceshipOff)

	// Instantiate the bissen
	b := NewBissen(textures)

	// Generate stars once
	generateStars()

	// Create the ground
	ground := NewGround()

	for !rl.WindowShouldClose() {
		updateStars(playerSpaceship.speed)
		update(b)
		renderScene(textures, playerSpaceship, b, ground)
	}
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

func DrawDustCloud(x, y int32, dustCloud rl.Texture2D) {
	posX := x + 10 // or adjust as needed
	posY := y + 30 // put it slightly under the spaceship

	fmt.Println("Dust cloud position:", posX, posY)
	rl.DrawTexture(dustCloud, posX, posY, rl.White)
}

func DrawBissen(b *bissen, textures textures) {
	// Draw the bissen
	if b.visible {
		b.x -= 2
		rl.DrawTexture(textures.bissen, b.x, b.y, rl.White)
	}

	if b.x < screenWidth/2 {
		b.visible = false
	}
}

func renderScene(textures textures, playerSpaceship *Spaceship, b *bissen, ground *ground) {
	rl.BeginDrawing()

	backgroundColor := rl.NewColor(skyColor.r, skyColor.g, skyColor.b, skyColor.a)
	rl.ClearBackground(backgroundColor)

	// Count frames
	frameCounter++

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

	playerSpaceship.y -= int32(playerSpaceship.speed)

	altitude = altitude + int32(playerSpaceship.speed)
	drawSpeed(playerSpaceship.speed, frameCounter, skyColor, 10)
	drawSpaceshipHeight(altitude)

	// This draws the green ground
	ground.draw(*playerSpaceship)

	// This draws the bissen
	DrawBissen(b, textures)

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

	if checkIfSpace(skyColor) {
		renderStars()
	}

	rl.EndDrawing()
}
