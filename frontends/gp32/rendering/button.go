package rendering

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Button struct {
	Text string
	PosX int
	PosY int
	Width int
	Height int
	Callback func ()
	Colour color.RGBA
}

func NewButton(text string, x int, y int, width int, height int, callback func ()) Button {
	return Button{
		Text: text,
		PosX: x,
		PosY: y,
		Width: width,
		Height: height,
		Callback: callback,
	}
}

func (b *Button) Draw() {
	
	rl.DrawRectangle(int32(b.PosX), int32(b.PosX), int32(b.Width), int32(b.Height), rl.LightGray)
	rl.DrawText(b.Text, int32(b.PosX), int32(b.PosY), 14, rl.Black)

}

func (b *Button) Update(p rl.Vector2) {
	
	if rl.CheckCollisionPointRec(
		p,
		rl.Rectangle{
			X:      float32(b.PosX),
			Y:      float32(b.PosY),
			Width:  float32(b.PosX+b.Width),
			Height: float32(b.PosY+b.Height),
		}) {

		b.Callback()

	}

}