package main

import (
	"fmt"

	gui2 "github.com/bvisness/jamtech/raylib/raygui"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth = 800
const screenHeight = 450

var ballPosition = rl.Vector2{screenWidth / 2, screenHeight / 2}

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	fmt.Println(gui2.GetStyle(0, 0))

	for !rl.WindowShouldClose() {
		doFrame()
	}
}

//var colors = []string{"Maroon", "Blue", "Green"}
var colors = "Maroon;Blue;Green"
var rlColors = []rl.Color{rl.Maroon, rl.Blue, rl.Green}
var selectedColor = 0

var ballLabel = "Ball"

var textBoxActive = false
var dropdownOpen = false

func doFrame() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	ballPosition.X = gui.SliderBar(rl.Rectangle{600, 40, 120, 20}, ballPosition.X, 0, screenWidth)
	ballPosition.Y = gui.SliderBar(rl.Rectangle{600, 70, 120, 20}, ballPosition.Y, 0, screenHeight)

	var toggleTextBox bool
	if ballLabel, toggleTextBox = gui2.TextBox(rl.Rectangle{40, 40, 120, 20}, ballLabel, 100, textBoxActive); toggleTextBox {
		textBoxActive = !textBoxActive
	}

	//selectedColor = gui.ComboBox(rl.Rectangle{40, 40, 120, 20}, colors, selectedColor)
	selectedColor = gui2.ComboBox(rl.Rectangle{40, 70, 120, 20}, colors, selectedColor)
	if gui2.DropdownBox(rl.Rectangle{40, 100, 120, 20}, colors, &selectedColor, dropdownOpen) {
		dropdownOpen = !dropdownOpen
	}
	//ballLabel = gui.TextBox(rl.Rectangle{40, 70, 120, 20}, ballLabel)

	rl.ClearBackground(rl.RayWhite)
	rl.DrawCircleV(ballPosition, 50, rlColors[selectedColor])
	rl.DrawText(ballLabel, int32(ballPosition.X), int32(ballPosition.Y), 14, rl.DarkGray)
}
