package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"strconv"
)

var (
	mplusNormalFont            font.Face
	screenWidth                float32 = 1000
	screenHeight               float32 = 800
	marginBetweenScreenAndArea float32 = 40.0
	areaBorderWidth            float32 = 1

	initialPaddleMarginFromScreen float32 = 75
	initialPaddleHeight           float32 = 100
	initialPaddleWidth            float32 = 10

	initialBallCenterX float32 = 500
	initialBallCenterY float32 = 400

	initialBallRadius float32 = 10
	initialBallAccX   float32 = 6
	initialBallAccY   float32 = 6
)

type Game struct {
	onPlaying                  bool
	screenWidth                float32
	screenHeight               float32
	marginBetweenScreenAndArea float32
	areaWidth                  float32
	areaBorderWidth            float32
	areaHeight                 float32
	areaBorderHeight           float32
	areaBorderColor            color.Color

	PlayerOne *Player
	PlayerTwo *Player
	Ball      *Ball

	StartKey ebiten.Key
}

func (g *Game) reset() {
	g.onPlaying = false
	initialBallCenter := &Coordinate{X: initialBallCenterX, Y: initialBallCenterY}
	g.Ball.Center = initialBallCenter
	g.onPlaying = false
}

func NewGame() *Game {

	initialBallCenter := &Coordinate{X: initialBallCenterX, Y: initialBallCenterY}

	playerOne := &Player{
		Paddle: &Paddle{
			Height:           initialPaddleHeight,
			Width:            initialPaddleWidth,
			MarginFromScreen: initialPaddleMarginFromScreen,
			Color:            color.White,
			acceleration:     8,
			position: &RectPosition{
				NorthLeft:  &Coordinate{X: initialPaddleMarginFromScreen, Y: screenHeight/2 - initialPaddleHeight/2},
				NorthRight: &Coordinate{X: initialPaddleMarginFromScreen + initialPaddleWidth, Y: screenHeight/2 - initialPaddleHeight/2},
				SouthLeft:  &Coordinate{X: initialPaddleMarginFromScreen, Y: screenHeight/2 - initialPaddleHeight/2 + initialPaddleHeight},
				SouthRight: &Coordinate{X: initialPaddleMarginFromScreen + initialPaddleWidth, Y: screenHeight/2 - initialPaddleHeight/2 + initialPaddleHeight},
			},
		},
		Score: &Score{
			position: &Coordinate{X: 30, Y: 30},
			value:    0,
			color:    color.White,
		},
		Name:    "Erkan",
		UpKey:   ebiten.KeyW,
		DownKey: ebiten.KeyS,
	}

	playerTwo := &Player{
		Paddle: &Paddle{
			Height:           initialPaddleHeight,
			Width:            initialPaddleWidth,
			MarginFromScreen: initialPaddleMarginFromScreen,
			Color:            color.White,
			acceleration:     8,
			position: &RectPosition{
				NorthLeft:  &Coordinate{X: screenWidth - initialPaddleMarginFromScreen - initialPaddleWidth, Y: screenHeight/2 - initialPaddleHeight/2},
				NorthRight: &Coordinate{X: screenWidth - initialPaddleMarginFromScreen, Y: screenHeight/2 - initialPaddleHeight/2},
				SouthLeft:  &Coordinate{X: screenWidth - initialPaddleMarginFromScreen - initialPaddleWidth, Y: (screenHeight/2 - initialPaddleHeight/2) + initialPaddleHeight},
				SouthRight: &Coordinate{X: screenWidth - initialPaddleMarginFromScreen, Y: (screenHeight/2 - initialPaddleHeight/2) + initialPaddleHeight},
			},
		},
		Score: &Score{
			position: &Coordinate{X: 950, Y: 30},
			value:    0,
			color:    color.White,
		},
		Name:    "Feyza",
		UpKey:   ebiten.KeyO,
		DownKey: ebiten.KeyL,
	}

	ball := &Ball{
		Center:            initialBallCenter,
		Radius:            initialBallRadius,
		ballAccelerationX: initialBallAccX,
		ballAccelerationY: initialBallAccY,
		Color:             color.White,
	}

	return &Game{
		onPlaying:                  false,
		screenWidth:                screenWidth,
		screenHeight:               screenHeight,
		marginBetweenScreenAndArea: marginBetweenScreenAndArea,
		areaBorderWidth:            areaBorderWidth,
		areaBorderHeight:           areaBorderWidth,
		areaWidth:                  screenWidth - marginBetweenScreenAndArea*2,
		areaHeight:                 screenHeight - marginBetweenScreenAndArea*2,
		PlayerOne:                  playerOne,
		PlayerTwo:                  playerTwo,
		Ball:                       ball,
		areaBorderColor:            color.White,
		StartKey:                   ebiten.KeySpace,
	}
}

type RectPosition struct {
	NorthLeft  *Coordinate
	NorthRight *Coordinate
	SouthLeft  *Coordinate
	SouthRight *Coordinate
}

type Player struct {
	Paddle  *Paddle
	Score   *Score
	Name    string
	UpKey   ebiten.Key
	DownKey ebiten.Key
}

type Score struct {
	position *Coordinate
	value    int
	color    color.Color
}

type Coordinate struct {
	X float32
	Y float32
}

type Paddle struct {
	Height           float32
	Width            float32
	MarginFromScreen float32
	Color            color.Color
	position         *RectPosition
	acceleration     float32
}

type Ball struct {
	Center            *Coordinate
	Radius            float32
	ballAccelerationX float32
	ballAccelerationY float32
	Color             color.Color
}

func (g *Game) Update() error {
	playerOne := g.PlayerOne
	playerTwo := g.PlayerTwo
	ball := g.Ball

	if ebiten.IsKeyPressed(playerOne.UpKey) && playerOne.Paddle.position.NorthLeft.Y > g.marginBetweenScreenAndArea {
		playerOne.Paddle.position.NorthLeft.Y = playerOne.Paddle.position.NorthLeft.Y - playerOne.Paddle.acceleration
	}
	if ebiten.IsKeyPressed(playerOne.DownKey) && playerOne.Paddle.position.NorthLeft.Y < g.screenHeight-g.marginBetweenScreenAndArea-playerOne.Paddle.Height {
		playerOne.Paddle.position.NorthLeft.Y = playerOne.Paddle.position.NorthLeft.Y + playerOne.Paddle.acceleration
	}

	if ebiten.IsKeyPressed(playerTwo.UpKey) && playerTwo.Paddle.position.NorthRight.Y > g.marginBetweenScreenAndArea {
		playerTwo.Paddle.position.NorthRight.Y = playerTwo.Paddle.position.NorthRight.Y - playerTwo.Paddle.acceleration
	}
	if ebiten.IsKeyPressed(playerTwo.DownKey) && playerTwo.Paddle.position.NorthRight.Y < g.screenHeight-g.marginBetweenScreenAndArea-playerTwo.Paddle.Height {
		playerTwo.Paddle.position.NorthRight.Y = playerTwo.Paddle.position.NorthRight.Y + playerTwo.Paddle.acceleration
	}
	// paddle positions rearrange
	playerOne.Paddle.position.NorthRight.X = playerOne.Paddle.position.NorthLeft.X + playerOne.Paddle.Width
	playerOne.Paddle.position.SouthRight.X = playerOne.Paddle.position.NorthLeft.X + playerOne.Paddle.Width
	playerOne.Paddle.position.NorthRight.Y = playerOne.Paddle.position.NorthLeft.Y
	playerOne.Paddle.position.SouthRight.Y = playerOne.Paddle.position.NorthLeft.Y + playerOne.Paddle.Height

	playerTwo.Paddle.position.NorthLeft.X = playerTwo.Paddle.position.NorthRight.X - playerTwo.Paddle.Width
	playerTwo.Paddle.position.SouthLeft.X = playerTwo.Paddle.position.NorthRight.X - playerTwo.Paddle.Width
	playerTwo.Paddle.position.NorthLeft.Y = playerTwo.Paddle.position.NorthRight.Y
	playerTwo.Paddle.position.SouthLeft.Y = playerTwo.Paddle.position.NorthRight.Y + playerTwo.Paddle.Height

	if !g.onPlaying && ebiten.IsKeyPressed(g.StartKey) {
		g.onPlaying = true
	}

	if g.onPlaying {

		ball.Center.X = ball.Center.X - (ball.ballAccelerationX)
		ball.Center.Y = ball.Center.Y - (ball.ballAccelerationY)
		// up line || bottom line
		if ball.Center.Y-ball.Radius <= g.marginBetweenScreenAndArea || ball.Center.Y+ball.Radius >= g.screenHeight-g.marginBetweenScreenAndArea {
			ball.ballAccelerationY = ball.ballAccelerationY * -1
		} else if abs(ball.Center.X-playerOne.Paddle.position.NorthRight.X) <= ball.Radius &&
			abs(playerOne.Paddle.position.NorthRight.Y-ball.Radius) <= ball.Center.Y &&
			ball.Center.Y <= abs(playerOne.Paddle.position.SouthRight.Y+ball.Radius) {
			// first player
			if ball.Center.Y < playerOne.Paddle.position.NorthRight.Y || ball.Center.Y > playerOne.Paddle.position.SouthRight.Y {
				ball.ballAccelerationY = ball.ballAccelerationY * -1
			}
			ball.ballAccelerationX = ball.ballAccelerationX * -1
		} else if abs(ball.Center.X-playerTwo.Paddle.position.NorthLeft.X) <= ball.Radius &&
			abs(playerTwo.Paddle.position.NorthLeft.Y-ball.Radius) <= ball.Center.Y &&
			ball.Center.Y <= abs(playerTwo.Paddle.position.SouthLeft.Y+ball.Radius) {
			// second player
			if ball.Center.Y < playerTwo.Paddle.position.NorthLeft.Y || ball.Center.Y > playerTwo.Paddle.position.SouthLeft.Y {
				ball.ballAccelerationY = ball.ballAccelerationY * -1
			}
			ball.ballAccelerationX = ball.ballAccelerationX * -1
		} /* else if ball.Center.X-ball.Radius <= g.marginBetweenScreenAndArea || ball.Center.X+ball.Radius >= g.screenWidth-g.marginBetweenScreenAndArea {
			// left and right side
			ball.ballAccelerationX = ball.ballAccelerationX * -1
		}*/
		if ball.Center.X-ball.Radius <= g.marginBetweenScreenAndArea {
			playerTwo.Score.value++
			g.reset()
		} else if ball.Center.X+ball.Radius >= g.screenWidth-g.marginBetweenScreenAndArea {
			playerOne.Score.value++
			g.reset()
		}
	}
	return nil
}

func abs(val float32) float32 {
	if val < 0 {
		return val * -1
	}
	return val
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, g.marginBetweenScreenAndArea, g.marginBetweenScreenAndArea, g.areaWidth, g.areaBorderWidth, g.areaBorderColor, false)
	vector.DrawFilledRect(screen, g.marginBetweenScreenAndArea, g.screenHeight-g.marginBetweenScreenAndArea, g.areaWidth, g.areaBorderWidth, g.areaBorderColor, false)
	vector.DrawFilledRect(screen, g.marginBetweenScreenAndArea, g.marginBetweenScreenAndArea, g.areaBorderWidth, g.areaHeight, g.areaBorderColor, false)
	vector.DrawFilledRect(screen, g.screenWidth-g.marginBetweenScreenAndArea, g.marginBetweenScreenAndArea, g.areaBorderWidth, g.areaHeight, g.areaBorderColor, false)

	vector.DrawFilledRect(screen, g.PlayerOne.Paddle.position.NorthLeft.X, g.PlayerOne.Paddle.position.NorthLeft.Y, g.PlayerOne.Paddle.Width, g.PlayerOne.Paddle.Height, g.PlayerOne.Paddle.Color, false)
	vector.DrawFilledRect(screen, g.PlayerTwo.Paddle.position.NorthLeft.X, g.PlayerTwo.Paddle.position.NorthLeft.Y, g.PlayerTwo.Paddle.Width, g.PlayerTwo.Paddle.Height, g.PlayerTwo.Paddle.Color, false)
	vector.DrawFilledCircle(screen, g.Ball.Center.X, g.Ball.Center.Y, g.Ball.Radius, g.Ball.Color, false)

	text.Draw(screen, strconv.Itoa(g.PlayerOne.Score.value), mplusNormalFont, int(g.PlayerOne.Score.position.X), int(g.PlayerOne.Score.position.Y), g.PlayerOne.Score.color)
	text.Draw(screen, strconv.Itoa(g.PlayerTwo.Score.value), mplusNormalFont, int(g.PlayerTwo.Score.position.X), int(g.PlayerTwo.Score.position.Y), g.PlayerTwo.Score.color)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.screenWidth), int(g.screenHeight)
}

func main() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(int(screenWidth), int(screenHeight))
	ebiten.SetWindowTitle("Pong")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
