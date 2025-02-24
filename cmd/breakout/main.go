package main

import (
	"fmt"
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
	paddlePosX         float32
	ballPos            rl.Vector2
	ballDir            rl.Vector2
	started            bool
	gameOver           bool
	blocks             [numBlocksX][numBlocksY]bool
	rowColors          [numBlocksY]rl.Color
	score              int
	blockScoreColor    map[rl.Color]int
	accumulatedTime    float32
	previousBallPos    rl.Vector2
	previousPaddlePosX float32
)

func calcBlockRect(x, y int) rl.Rectangle {
	return rl.NewRectangle(
		float32(60+x*blockWidth),
		float32(40+y*blockHeight),
		blockWidth,
		blockHeight,
	)
}

func restart() {
	paddlePosX = screenSize/2 - paddleWidth/2
	previousPaddlePosX = paddlePosX
	ballPos = rl.NewVector2(screenSize/2, ballStartY)
	previousBallPos = ballPos
	started = false
	gameOver = false
	score = 0
	for i := 0; i < numBlocksX; i++ {
		for j := 0; j < numBlocksY; j++ {
			blocks[i][j] = true
		}
	}
	blockScoreColor = map[rl.Color]int{
		rl.Yellow: 2,
		rl.Green:  4,
		rl.Orange: 6,
		rl.Red:    8,
	}
}

func reflect(dir, normal rl.Vector2) rl.Vector2 {
	newDirection := rl.Vector2Reflect(dir, rl.Vector2Normalize(normal))
	return rl.Vector2Normalize(newDirection)
}

func blockExists(x, y int) bool {
	if x < 0 || x >= numBlocksX {
		return false
	}
	if y < 0 || y >= numBlocksY {
		return false
	}
	return blocks[x][y]
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
	rl.InitWindow(750, 750, "breakout")
	rl.InitAudioDevice()
	rl.SetTargetFPS(500)
	backgroundColor := rl.NewColor(150, 190, 220, 255)
	ballTexture := rl.LoadTexture("assets/ball.png")
	paddleTexture := rl.LoadTexture("assets/paddle.png")
	hitBlockSound := rl.LoadSound("assets/hit_block.wav")
	hitPaddleSound := rl.LoadSound("assets/hit_paddle.wav")
	gameOverSound := rl.LoadSound("assets/game_over.wav")

	restart()
	for !rl.WindowShouldClose() {
		const dt = 1.0 / 60.0
		switch {
		case !started:
			ballPos = rl.NewVector2(screenSize/2+float32(math.Cos(rl.GetTime())*screenSize/2.5), ballStartY)
			previousBallPos = ballPos
			if rl.IsKeyPressed(rl.KeySpace) {
				paddleMiddle := rl.NewVector2(paddlePosX+paddleWidth/2, paddlePosY)
				ballToPaddle := rl.Vector2Subtract(paddleMiddle, ballPos)
				ballDir = rl.Vector2Normalize(ballToPaddle)
				started = true
			}
		case gameOver:
			if rl.IsKeyPressed(rl.KeySpace) {
				restart()
			}
		default:
			accumulatedTime += rl.GetFrameTime()
		}

		for accumulatedTime >= dt {
			previousBallPos = ballPos
			previousPaddlePosX = paddlePosX
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

			if !gameOver && ballPos.Y > screenSize+ballRadius*6 {
				gameOver = true
				rl.PlaySound(gameOverSound)
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
				rl.PlaySound(hitPaddleSound)
			}

		out:
			for x := 0; x < numBlocksX; x++ {
				for y := 0; y < numBlocksY; y++ {
					if !blocks[x][y] {
						continue
					}
					blockRect := calcBlockRect(x, y)
					if rl.CheckCollisionCircleRec(ballPos, ballRadius, blockRect) {
						collisionNormal := rl.NewVector2(0, 0)
						if previousBallPos.Y < blockRect.Y {
							collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(0, -1))
						}
						if previousBallPos.Y > blockRect.Y+blockRect.Height {
							collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(0, 1))
						}

						if previousBallPos.X < blockRect.X {
							collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(-1, 0))
						}
						if previousBallPos.X > blockRect.X+blockRect.Width {
							collisionNormal = rl.Vector2Add(collisionNormal, rl.NewVector2(1, 0))
						}

						if blockExists(x+int(collisionNormal.X), y) {
							collisionNormal.X = 0
						}

						if blockExists(x, y+int(collisionNormal.Y)) {
							collisionNormal.Y = 0
						}

						if collisionNormal.X != 0 || collisionNormal.Y != 0 {
							ballDir = reflect(ballDir, collisionNormal)
						}
						rowColor := rowColors[y]
						score += blockScoreColor[rowColor]
						rl.SetSoundPitch(hitBlockSound, float32(rl.GetRandomValue(8, 12))/10)
						rl.PlaySound(hitBlockSound)
						blocks[x][y] = false
						break out
					}
				}
			}
			accumulatedTime -= dt
		}
		blendFactor := accumulatedTime / dt
		ballRenderPos := rl.Vector2Lerp(previousBallPos, ballPos, blendFactor)
		paddleRenderPosX := rl.Lerp(previousPaddlePosX, paddlePosX, blendFactor)

		rl.BeginDrawing()

		rl.ClearBackground(backgroundColor)
		camera := rl.Camera2D{
			Zoom: float32(rl.GetScreenHeight() / screenSize),
		}
		rl.BeginMode2D(camera)
		rl.DrawTextureV(paddleTexture, rl.NewVector2(paddleRenderPosX, paddlePosY), rl.White)
		rl.DrawTextureV(ballTexture, rl.Vector2Subtract(ballRenderPos, rl.NewVector2(ballRadius, ballRadius)), rl.White)

		for x := 0; x < numBlocksX; x++ {
			for y := 0; y < numBlocksY; y++ {
				if !blocks[x][y] {
					continue
				}
				rect := calcBlockRect(x, y)
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
		rl.DrawText(fmt.Sprint(score), 5, 5, 10, rl.White)

		if !started {
			startText := "Press Space to Start"
			startTextWidth := rl.MeasureText(startText, 15)
			rl.DrawText(startText, screenSize/2-startTextWidth/2, ballStartY-30, 15, rl.White)
		}

		if gameOver {
			gameOverText := fmt.Sprintf("Game Over\nScore: %d", score)
			gameOverTextWidth := rl.MeasureText(gameOverText, 15)
			rl.DrawText(gameOverText, screenSize/2-gameOverTextWidth/2, ballStartY-30, 15, rl.White)
		}

		rl.EndMode2D()
		rl.EndDrawing()
	}
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
