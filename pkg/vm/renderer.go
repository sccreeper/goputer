package vm

import (
	_ "embed"
	"fmt"
	c "sccreeper/goputer/pkg/constants"
)

//go:embed assets/font.bin
var fontBytes []byte
var fontData [FontNumCharacters][]byte = [FontNumCharacters][]byte{}

const (
	VideoBufferWidth       uint32 = 320
	VideoBufferHeight      uint32 = 240
	VideoBufferColourDepth uint32 = 24
	VideoBytesPerPixel     uint32 = VideoBufferColourDepth / 8
	VideoBufferSize        uint32 = VideoBufferWidth * VideoBufferHeight * VideoBytesPerPixel

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

	// Do the first line

	if m.Registers[c.RVideoX0] > m.Registers[c.RVideoX1] || m.Registers[c.RVideoY0] > m.Registers[c.RVideoY1] {
		// TODO add some sort of panic
		return
	}

	var colour = m.getVideoColour()

	var areaStart = (m.Registers[c.RVideoX0] * VideoBytesPerPixel) + (m.Registers[c.RVideoY0] * VideoBufferWidth * VideoBytesPerPixel)

	if colour[3] == 255 {

		// Do the first row

		for x := 0; x < int(m.Registers[c.RVideoX1]-m.Registers[c.RVideoX0]); x++ {
			var pixelAddr int = int(areaStart) + (x * int(VideoBytesPerPixel))

			m.MemArray[pixelAddr] = colour[0]
			m.MemArray[pixelAddr+1] = colour[1]
			m.MemArray[pixelAddr+2] = colour[2]
		}

		// Copy each row thereafter

		for y := 1; y < int(m.Registers[c.RVideoY1]-m.Registers[c.RVideoY0]); y++ {
			var pixelAddrStart = int(areaStart) + (y * int(VideoBufferWidth) * int(VideoBytesPerPixel))

			copy(
				m.MemArray[pixelAddrStart:pixelAddrStart+(int(VideoBufferWidth*VideoBytesPerPixel))],
				m.MemArray[areaStart:areaStart+(VideoBufferWidth*VideoBytesPerPixel)],
			)
		}

	} else {

		for y := 0; y < int(m.Registers[c.RVideoY1]-m.Registers[c.RVideoY0]); y++ {
			for x := 0; x < int(m.Registers[c.RVideoX1]-m.Registers[c.RVideoX0]); x++ {

				var pixelAddr int = int(areaStart) + (x * int(VideoBytesPerPixel)) + (y * int(VideoBufferWidth) * int(VideoBytesPerPixel))

				var blendedColour [3]byte = blendPixel(colour, [3]byte(m.MemArray[pixelAddr:pixelAddr+3]))
				m.MemArray[pixelAddr] = blendedColour[0]
				m.MemArray[pixelAddr+1] = blendedColour[1]
				m.MemArray[pixelAddr+2] = blendedColour[2]

			}
		}

	}

}

func (m *VM) drawLine() {

}

func (m *VM) drawText() {

}

func (m *VM) drawPolygon() {

}

func (m *VM) drawImage() {

}

func (m *VM) getVideoColour() (colour [4]byte) {

	colour = [4]byte{}

	colour[0] = byte(m.Registers[c.RVideoColour])
	colour[1] = byte(m.Registers[c.RVideoColour] >> 8)
	colour[2] = byte(m.Registers[c.RVideoColour] >> 16)
	colour[3] = byte(m.Registers[c.RVideoColour] >> 24)

	return colour

}

// Performs a sort of TM blend on two pixels
func blendPixel(src [4]byte, dest [3]byte) [3]byte {

	var invertedAlpha int = 255 - int(src[3])

	return [3]byte{
		byte((int(src[0])*int(src[3]) + int(dest[0])*invertedAlpha) / 255),
		byte((int(src[1])*int(src[3]) + int(dest[1])*invertedAlpha) / 255),
		byte((int(src[2])*int(src[3]) + int(dest[2])*invertedAlpha) / 255),
	}

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
