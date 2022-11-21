package rendering

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ColourIOOn rl.Color = rl.Color{
	R: 251,
	G: 250,
	B: 169,
	A: 255,
}

var ColourIOOff rl.Color = rl.Color{
	R: 46,
	G: 46,
	B: 46,
	A: 255,
}

var SwitchToggleBackground rl.Color = rl.Color{
	R: 104,
	G: 104,
	B: 104,
	A: 255,
}

func RenderIO(status []bool, switches []IOSwitch) {

	rl.ClearBackground(rl.DarkGray)

	for index, v := range status {

		var c rl.Color

		//Colour depending whether or not IO is on or off.
		if v == true {
			c = ColourIOOn
		} else {
			c = ColourIOOff
		}

		if index < 8 {
			rl.DrawRectangleRounded(rl.Rectangle{
				Width:  float32(IOUISize),
				Height: float32(IOUISize),
				X:      float32(index * IOUISize),
				Y:      0,
			},
				0.5,
				16,
				c,
			)

			rl.DrawRectangleRoundedLines(rl.Rectangle{
				Width:  float32(IOUISize),
				Height: float32(IOUISize),
				X:      float32(index * IOUISize),
				Y:      0,
			},
				0.5,
				16,
				1,
				rl.Black,
			)
		} else {
			for _, v := range switches {

				v.RenderSwitch()

			}
		}

	}

}

type IOSwitch struct {
	Toggled bool
	ID      uint32
	X       float32
	Y       float32
}

func (s *IOSwitch) RenderSwitch() {

	//Draw background

	rl.DrawRectangleRounded(rl.Rectangle{
		Width:  float32(IOUISize),
		Height: float32(IOUISize),
		X:      s.X,
		Y:      s.Y,
	},
		0.5,
		16,
		ColourIOOff,
	)

	rl.DrawRectangleRoundedLines(rl.Rectangle{
		Width:  float32(IOUISize),
		Height: float32(IOUISize),
		X:      s.X,
		Y:      s.Y,
	},
		0.5,
		16,
		1,
		rl.Black,
	)

	//Draw actual switches

	rl.DrawRectangle(int32(s.X)+5, int32(s.Y)+5, int32(IOUISize)-5, (int32(IOUISize)/2)-5, SwitchToggleBackground)

	if s.Toggled {

		rl.DrawRectangle(
			int32(s.X)+int32(IOUISize)/2,
			int32(s.Y)+5,
			int32(IOUISize)/2,
			(int32(IOUISize)/2)-5,
			rl.White,
		)

		rl.DrawText("On", int32(s.X)+5, int32(s.Y)+5+(int32(IOUISize)/2), 8, rl.White)

	} else {

		rl.DrawRectangle(
			int32(s.X)+5,
			int32(s.Y)+5,
			int32(IOUISize)/2,
			(int32(IOUISize)/2)-5,
			rl.White,
		)

		rl.DrawText("Off", int32(s.X)+5, int32(s.Y)+5+(int32(IOUISize)/2), 8, rl.White)
	}

}

// Called when a mouse click has been detected
func (s *IOSwitch) Update(p rl.Vector2) bool {

	if rl.CheckCollisionPointRec(
		p,
		rl.Rectangle{
			X:      s.X,
			Y:      s.Y + float32(DebugUISize),
			Width:  float32(IOUISize),
			Height: float32(IOUISize),
		}) {

		s.Toggled = !s.Toggled

		return true

	}

	return false

}
