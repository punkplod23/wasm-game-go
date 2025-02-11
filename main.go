package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
)

const (
	rectangleSize = 1
	startGame     = false
	endGame       = false
	snakeSize     = 20
	snakeSpeed    = 21
)

func clearCanvas(ctx js.Value, canvas js.Value) {
	// Clear the canvas
	ctx.Call("clearRect", 0, 0, canvas.Get("width").Int(), canvas.Get("height").Int())
}
func checkCollision(canvas js.Value) bool {
	head := snake[0]
	headX := head.Get("x").Int()
	headY := head.Get("y").Int()

	// Check collision with walls (color black)
	if headX < 0 || headY < 0 || headX >= canvas.Get("width").Int() || headY >= canvas.Get("height").Int() {
		return true
	}

	// Check collision with itself
	for _, segment := range snake[1:] {
		if headX == segment.Get("x").Int() && headY == segment.Get("y").Int() {
			return true
		}
	}

	return false
}

func drawGrid(ctx js.Value, canvas js.Value) {
	clearCanvas(ctx, canvas)
	//Create a Map
	canvasWidth := canvas.Get("width").Int()
	canvasHeight := canvas.Get("height").Int()
	rows := canvasHeight / rectangleSize
	cols := canvasWidth / rectangleSize

	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
		for j := range grid[i] {
			if i == 0 || i == rows-1 || j == 0 || j == cols-1 {
				grid[i][j] = 1
			} else {
				grid[i][j] = 0
			}
		}
	}

	size := rectangleSize

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 1 {
				draw(ctx, "black", size, j*rectangleSize, i*rectangleSize)
			}
			if grid[i][j] == 0 {
				draw(ctx, "white", size, j*rectangleSize, i*rectangleSize)
			}
		}
	}
}

var (
	snake     = []js.Value{}
	direction = "right"
)
var food js.Value

func placeFood(canvas js.Value) {
	canvasWidth := canvas.Get("width").Int()
	canvasHeight := canvas.Get("height").Int()
	foodX := (rand.Intn(canvasWidth/rectangleSize) * rectangleSize)
	foodY := (rand.Intn(canvasHeight/rectangleSize) * rectangleSize)
	food = js.ValueOf(map[string]interface{}{
		"x": foodX,
		"y": foodY,
	})
}

func drawFood(ctx js.Value) {
	draw(ctx, "red", snakeSize, food.Get("x").Int(), food.Get("y").Int())
}

func moveSnake(ctx js.Value, canvas js.Value) {
	head := snake[0]
	newHead := js.ValueOf(map[string]interface{}{
		"x": head.Get("x").Int(),
		"y": head.Get("y").Int(),
	})

	switch direction {
	case "right":
		newHead.Set("x", newHead.Get("x").Int()+snakeSpeed)
	case "left":
		newHead.Set("x", newHead.Get("x").Int()-snakeSpeed)
	case "up":
		newHead.Set("y", newHead.Get("y").Int()-snakeSpeed)
	case "down":
		newHead.Set("y", newHead.Get("y").Int()+snakeSpeed)
	}

	snake = append([]js.Value{newHead}, snake[:len(snake)-1]...)
	drawSnake(ctx)
}

func drawSnake(ctx js.Value) {
	for _, segment := range snake {
		draw(ctx, "green", snakeSize, segment.Get("x").Int(), segment.Get("y").Int())
	}
}

func handleKeyPress(event js.Value) {
	key := event.Get("key").String()
	switch key {
	case "ArrowUp":
		if direction != "down" {
			direction = "up"
		}
	case "ArrowDown":
		if direction != "up" {
			direction = "down"
		}
	case "ArrowLeft":
		if direction != "right" {
			direction = "left"
		}
	case "ArrowRight":
		if direction != "left" {
			direction = "right"
		}
	}
}

func gameLoop(ctx js.Value, canvas js.Value) {
	drawGrid(ctx, canvas)
	moveSnake(ctx, canvas)
	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		gameLoop(ctx, canvas)
		if checkCollision(canvas) {
			fmt.Println("Game Over!")
			titleScreen(ctx, canvas)
		}
		return nil
	}), 100)
}

func titleScreen(ctx js.Value, canvas js.Value) {
	clearCanvas(ctx, canvas)
	//ctx.Set("imageSmoothingEnabled", true)
	ctx.Set("imageSmoothingEnabled", false)
	// Set fill color to blue
	ctx.Set("fillStyle", "blue")

	// Draw a rectangle
	ctx.Call("fillRect", 0, 0, canvas.Get("width").Int(), canvas.Get("height").Int())

	//placeFood(canvas)
	//drawFood(ctx)

	// Set font and alignment for the title
	ctx.Set("font", "30px Arial")
	ctx.Set("textAlign", "center")
	ctx.Set("fillStyle", "white")
	ctx.Set("imageSmoothingEnabled", true)
	// Draw the title text
	ctx.Call("fillText", "Snake", canvas.Get("width").Int()/2, canvas.Get("height").Int()/2)

	// Add clickable text to start the game
	ctx.Set("font", "20px Arial")
	ctx.Set("fillStyle", "white")
	ctx.Call("fillText", "Click to Start", canvas.Get("width").Int()/2, canvas.Get("height").Int()/2+40)

	// Add event listener for click event to start the game
	canvas.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Starting Game!")
		// Initial position of the snake
		// Set fill color to green for the snake
		ctx.Set("fillStyle", "green")
		snake = []js.Value{
			js.ValueOf(map[string]interface{}{"x": 50, "y": 50}),
			js.ValueOf(map[string]interface{}{"x": 50 - rectangleSize, "y": 50}),
			js.ValueOf(map[string]interface{}{"x": 50 - 2*rectangleSize, "y": 50}),
		}

		// Draw the snake
		drawSnake(ctx)
		document := js.Global().Get("document")
		document.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			handleKeyPress(args[0])
			return nil
		}))
		gameLoop(ctx, canvas)
		return nil
	}))
}

func draw(ctx js.Value, colour string, size int, x int, y int) {

	ctx.Set("fillStyle", colour)

	// Set fill color to blue
	ctx.Set("fillSize", 1)

	// Draw a rectangle
	ctx.Call("fillRect", x, y, size, size)
}

func main() {
	fmt.Println("Hello, World!")
	// Get the document object
	document := js.Global().Get("document")

	// Get the canvas element by ID
	canvas := document.Call("getElementById", "wasmcanvas")

	// Get the 2D drawing context
	ctx := canvas.Call("getContext", "2d")
	titleScreen(ctx, canvas)
	select {}

}
