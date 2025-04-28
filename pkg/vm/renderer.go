package vm

import (
	_ "embed"
	"fmt"
)

//go:embed assets/font.bin
var fontBytes []byte
var fontData [FontNumCharacters][]byte = [FontNumCharacters][]byte{}

const (
	VideoBufferWidth       uint32 = 320
	VideoBufferHeight      uint32 = 240
	VideoBufferColourDepth uint32 = 8
	VideoBufferSize        uint32 = VideoBufferWidth * VideoBufferHeight * (VideoBufferColourDepth / 8)

	FontNumCharacters uint32 = 96
	FontCharWidth     uint32 = 5
	FontCharHeight    uint32 = 7
	FontCharGap       uint32 = 1
)

func init() {

	// Populate character array

	for i := 0; i < int(FontNumCharacters); i++ {
		fontData[i] = make([]byte, 0, FontCharWidth*FontCharHeight)

		dataStart := i * int(FontCharWidth) * 4

		for y := 0; y < int(FontCharHeight); y++ {
			for x := 0; x < int(FontCharWidth); x++ {

				pixelIndex := dataStart + (y * int(FontCharWidth*FontNumCharacters) * 4) + (x * 4)
				fontData[i] = append(fontData[i], fontBytes[pixelIndex])

			}
		}
	}

}

func (m *VM) drawSquare() {

}

func (m *VM) drawLine() {

}

func (m *VM) drawText() {

}

func (m *VM) drawPolygon() {

}

func (m *VM) drawImage() {

}

func PrintChar(char int) {
	if char >= int(FontNumCharacters) || char < 32 {
		char = 127
	}

	char -= 32

	for i := 0; i < len(fontData[char]); i += int(FontCharWidth) {
		var lineString string

		for _, v := range fontData[char][i : i+int(FontCharWidth)] {
			if v == 255 {
				lineString += "#"
			} else {
				lineString += " "
			}
		}

		fmt.Println(lineString)

	}
}
