package colour

import (
	"encoding/binary"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ConvertColour(c uint32) rl.Color {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[:], c)

	return rl.NewColor(b[0], b[1], b[2], b[3])

}
