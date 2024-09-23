package main

import (
	"strconv"
  "os"
	"path/filepath"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth      = int32(450)
	screenHeight     = int32(800)
	scaleFactor      = 0.05 // Scale down the spaceship to 5% of its original size
	maxSpeed         = 20
	terminalVelocity = -20 
	skyR             = 189
	skyG             = 245
	skyB             = 255
	starSpeed        = 1.0 // Speed at which stars move downwards
)

type Spaceship struct {
	x, y  int32
	speed float32
}

type SkyItem struct {
	name     string
	position rl.Vector2
	texture  rl.Texture2D
	speed    float32
}

func NewSkyItem(name string, position rl.Vector2, texturePath string) *SkyItem {
	skyItem := SkyItem{
		name:     name,
		position: position,
		texture:  createTextureFromImage(texturePath, 0.1),
		speed:    0.5,
	}
	skyObjectsList = append(skyObjectsList, skyItem)
	return &skyItem
}

func NewSpaceship(spaceship rl.Texture2D) *Spaceship {
	return &Spaceship{
		x:     screenWidth/2 - spaceship.Width/2,
		y:     screenHeight - spaceship.Height,
		speed: 0,
	}
}

func (s *Spaceship) Accelerate() {
	if s.speed < maxSpeed {
		s.speed += 1.05 
	}
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

type Sky struct {
	r, g, b, a uint8
}

func darkenSky(sky Sky) Sky {
	if sky.r > 0 {
		sky.r -= 1
	}
	if sky.g > 0 {
		sky.g -= 1
	}
	if sky.b > 0 {
		sky.b -= 1
	}
	return sky
}

func lightenSky(sky Sky) Sky {
	if sky.r < skyR {
		sky.r += 1
	}
	if sky.g < skyG {
		sky.g += 1
	}
	if sky.b < skyB {
		sky.b += 1
	}
	return sky
}

func checkIfSpace(sky Sky) bool {
	var limit uint8 = 100 
	if sky.r < limit && sky.g < limit && sky.b < limit {
		return true
	}
	return false
}

var (
	frameCounter int32 = 0
	drawnSpeed   float32
	skyColor = Sky{skyR, skyG, skyB, 255}
	stars        []rl.Vector2
	skyObjectsList   []SkyItem
)

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

	// Load the spaceship image
	spaceship := createTextureFromImage("/graphics/spaceship.png", 0.05)
	spaceshipOff := createTextureFromImage("/graphics/spaceshipOff.png", 0.05)

  // Create the spaceship struct
	playerSpaceship := NewSpaceship(spaceshipOff)

	// Generate stars once
	generateStars()

	for !rl.WindowShouldClose() {
		updateStars(playerSpaceship.speed)
		renderScene(spaceship, spaceshipOff, playerSpaceship)
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

func DrawDustCloud(x, y int32) {
	var dustCloud rl.Texture2D
	var posX, posY int32
	if dustCloud.Height == 0 {
		posX = screenWidth/2 + 20 
		posY = screenHeight - 50
		dustCloud = createTextureFromImage("/graphics/dustCloud.png", 0.10)
	}
	rl.DrawTexture(dustCloud, posX, posY, rl.White)
}

func renderScene(spaceship, spaceshipOff rl.Texture2D, playerSpaceship *Spaceship) {
	rl.BeginDrawing()

	backgroundColor := rl.NewColor(skyColor.r, skyColor.g, skyColor.b, skyColor.a)
	rl.ClearBackground(backgroundColor)

	// Count frames
	frameCounter++

	spaceShipToDraw := spaceshipOff
	if rl.IsKeyDown(rl.KeySpace) {
		spaceShipToDraw = spaceship
		playerSpaceship.Accelerate()
	}

	if playerSpaceship.speed > terminalVelocity {
		playerSpaceship.speed -= 1 
	}

	// Detect collision with the ground
	if playerSpaceship.y >= screenHeight-spaceship.Height && playerSpaceship.speed < 0 {
		playerSpaceship.y = screenHeight - spaceship.Height
		playerSpaceship.speed = 0
	}

	playerSpaceship.y -= int32(playerSpaceship.speed)

	drawSpeed(playerSpaceship.speed, frameCounter, skyColor, 10)

	customGreen := rl.NewColor(52, 235, 131, 255) // Custom color #34eb83

	heightOverHalf := screenHeight/2 - spaceship.Height/2

	cloud1 := findObjectInList("cloud1")
	if cloud1 == nil {
		cloud1 = NewSkyItem("cloud1", rl.Vector2{X: 50, Y: 50}, "/graphics/cloud.png")
	}
	if playerSpaceship.y > heightOverHalf {
		rl.DrawRectangle(0, screenHeight-70, screenWidth, 70, customGreen)
		rl.DrawTexture(spaceShipToDraw, playerSpaceship.x, playerSpaceship.y, rl.White)

		cloud1.move(0)

		// Detect if spaceship on the ground but starting to fly
		if playerSpaceship.y == screenHeight-spaceship.Height && playerSpaceship.speed > 0 {
			DrawDustCloud(playerSpaceship.x, playerSpaceship.y)
		}
	} else {
		rl.DrawRectangle(0, screenHeight-(70+(playerSpaceship.y-heightOverHalf)), screenWidth, 70, customGreen)
		rl.DrawTexture(spaceShipToDraw, playerSpaceship.x, screenHeight/2-spaceship.Height/2, rl.White)
		cloud1.move(int32(playerSpaceship.speed))

		if playerSpaceship.speed > 0 {
			skyColor = darkenSky(skyColor)
		}
		if playerSpaceship.speed < 0 {
			skyColor = lightenSky(skyColor)
		}
	}

	cloud1.draw()

	if checkIfSpace(skyColor) {
		renderStars()
	}

	rl.EndDrawing()
}

