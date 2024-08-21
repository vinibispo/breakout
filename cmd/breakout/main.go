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
	numBlocksX   = 10
	numBlocksY   = 8
	blockWidth   = 20
	blockHeight  = 10
)

var (
	paddlePosX float32
	ballPos    rl.Vector2
	ballDir    rl.Vector2
	started    bool
	blocks     [numBlocksX][numBlocksY]bool
	rowColors  = [numBlocksY]rl.Color{}
)

func restart() {
	paddlePosX = screenSize/2 - paddleWidth/2
	ballPos = rl.NewVector2(screenSize/2, ballStartY)
	started = false
	for i := 0; i < numBlocksX; i++ {
		for j := 0; j < numBlocksY; j++ {
			blocks[i][j] = true
		}
	}
}

func reflect(dir, normal rl.Vector2) rl.Vector2 {
	newDirection := rl.Vector2Reflect(dir, rl.Vector2Normalize(normal))
	return rl.Vector2Normalize(newDirection)
}
func main() {
	rowColors = [numBlocksY]rl.Color{
		rl.Red,
		rl.Red,
		rl.Orange,
		rl.Orange,
		rl.Green,
		rl.Green,
		rl.Yellow,
		rl.Yellow,
	}
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
		previousBallPos := ballPos
		ballPos = rl.Vector2Add(ballPos, rl.Vector2Scale(ballDir, ballSpeed*dt))
		if ballPos.X+ballRadius > screenSize {
			ballPos.X = screenSize - ballRadius
			ballDir = reflect(ballDir, rl.NewVector2(-1, 0))
		}
		if ballPos.X-ballRadius < 0 {
			ballPos.X = ballRadius
			ballDir = reflect(ballDir, rl.NewVector2(1, 0))
		}

		if ballPos.Y-ballRadius < 0 {
			ballPos.Y = ballRadius
			ballDir = reflect(ballDir, rl.NewVector2(0, 1))
		}

		if ballPos.Y > screenSize+ballRadius*6 {
			restart()
		}
		var paddleMoveVelocity float32
		if rl.IsKeyDown(rl.KeyLeft) {
			paddleMoveVelocity -= paddleSpeed
		}
		if rl.IsKeyDown(rl.KeyRight) {
			paddleMoveVelocity += paddleSpeed
		}

		paddlePosX += paddleMoveVelocity * dt
		paddlePosX = rl.Clamp(paddlePosX, 0, screenSize-paddleWidth)
		paddleRect := rl.NewRectangle(paddlePosX, paddlePosY, paddleWidth, paddleHeight)
		if rl.CheckCollisionCircleRec(ballPos, ballRadius, paddleRect) {
			collisionNormal := rl.NewVector2(0, 0)

			if previousBallPos.Y < (paddleRect.Y + paddleRect.Height) {
				collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(0, -1))
				ballPos.Y = paddleRect.Y - ballRadius
			}

			if previousBallPos.Y > (paddleRect.Y + paddleRect.Height) {
				collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(0, 1))
				ballPos.Y = paddleRect.Y + paddleRect.Height + ballRadius
			}

			if previousBallPos.X < paddleRect.X {
				collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(-1, 0))
			}

			if previousBallPos.X > paddleRect.X+paddleRect.Width {
				collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(1, 0))
			}

			if collisionNormal.X != 0 || collisionNormal.Y != 0 {
				ballDir = reflect(ballDir, collisionNormal)
			}
		}
		rl.BeginDrawing()

		rl.ClearBackground(backgroundColor)
		camera := rl.Camera2D{
			Zoom: float32(rl.GetScreenHeight() / screenSize),
		}
		rl.BeginMode2D(camera)
		rl.DrawRectangleRec(paddleRect, paddleColor)
		rl.DrawCircleV(ballPos, ballRadius, rl.Red)
		for x := 0; x < numBlocksX; x++ {
			for y := 0; y < numBlocksY; y++ {
				if !blocks[x][y] {
					continue
				}
				rect := rl.NewRectangle(
					float32(60+x*blockWidth),
					float32(40+y*blockHeight),
					blockWidth,
					blockHeight,
				)
				rl.DrawRectangleRec(rect, rowColors[y])
				topLeft := rl.NewVector2(rect.X, rect.Y)
				topRight := rl.NewVector2(rect.X+rect.Width, rect.Y)
				bottomLeft := rl.NewVector2(rect.X, rect.Y+rect.Height)
				bottomRight := rl.NewVector2(rect.X+rect.Width, rect.Y+rect.Height)
				rl.DrawLineEx(topLeft, topRight, 1, rl.NewColor(255, 255, 150, 100))
				rl.DrawLineEx(topLeft, bottomLeft, 1, rl.NewColor(255, 255, 150, 100))
				rl.DrawLineEx(topRight, bottomRight, 1, rl.NewColor(0, 0, 50, 100))
				rl.DrawLineEx(bottomLeft, bottomRight, 1, rl.NewColor(0, 0, 50, 100))
			}
		}

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
