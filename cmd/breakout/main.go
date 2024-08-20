package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenSize   = 320
	paddleWidth  = 50
	paddleHeight = 6
	paddlePosY   = 260
	ballSpeed    = 130
	ballRadius   = 4
	ballStartY   = 160
	paddleSpeed  = 100
)

var (
	paddlePosX float32
	ballPos    rl.Vector2
	ballDir    rl.Vector2
	started    bool
)

func restart() {
	paddlePosX = screenSize/2 - paddleWidth/2
	ballPos = rl.NewVector2(screenSize/2, ballStartY)
}
func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(700, 700, "breakout")
	rl.SetTargetFPS(500)
	backgroundColor := rl.NewColor(150, 190, 220, 255)
	paddleColor := rl.NewColor(50, 150, 90, 255)

	restart()
	for !rl.WindowShouldClose() {
		var dt float32
		if !started {
			ballPos = rl.NewVector2(screenSize/2+float32(math.Cos(rl.GetTime())*screenSize/2.5), ballStartY)
			if rl.IsKeyPressed(rl.KeySpace) {
				paddleMiddle := rl.NewVector2(paddlePosX+paddleWidth/2, paddlePosY)
				ballToPaddle := rl.Vector2Subtract(paddleMiddle, ballPos)
				ballDir = rl.Vector2Normalize(ballToPaddle)
				started = true
			}
		} else {
			dt = rl.GetFrameTime()
		}
		ballPos = rl.Vector2Add(ballPos, rl.Vector2Scale(ballDir, ballSpeed*dt))
		var paddleMoveVelocity float32
		if rl.IsKeyDown(rl.KeyLeft) {
			paddleMoveVelocity -= paddleSpeed
		}
		if rl.IsKeyDown(rl.KeyRight) {
			paddleMoveVelocity += paddleSpeed
		}

		paddlePosX += paddleMoveVelocity * dt
		paddlePosX = rl.Clamp(paddlePosX, 0, screenSize-paddleWidth)
		rl.BeginDrawing()

		rl.ClearBackground(backgroundColor)
		camera := rl.Camera2D{
			Zoom: float32(rl.GetScreenHeight() / screenSize),
		}
		rl.BeginMode2D(camera)
		paddleRect := rl.NewRectangle(paddlePosX, paddlePosY, paddleWidth, paddleHeight)
		rl.DrawRectangleRec(paddleRect, paddleColor)
		rl.DrawCircleV(ballPos, ballRadius, rl.Red)

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
