package rendering

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	IOUISize int = 40
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

func RenderIO(status []bool) {

	for index, v := range status {

		var c rl.Color

		//Colour depending whether or not IO is on or off.
		if v == true {
			c = ColourIOOn
		} else {
			c = ColourIOOff
		}

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

	}

}
