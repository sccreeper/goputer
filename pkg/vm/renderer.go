package vm

import (
	_ "embed"
	"fmt"
	"math"
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

	// Bresenham's algorithm

	var x int = int(m.Registers[c.RVideoX0])
	var y int = int(m.Registers[c.RVideoY0])

	var dx int = int(math.Abs(float64(m.Registers[c.RVideoX1] - m.Registers[c.RVideoX0])))
	var dy int = -int(math.Abs(float64(m.Registers[c.RVideoY1] - m.Registers[c.RVideoY0])))
	var err int = dx + dy

	var sx int = 1
	if m.Registers[c.RVideoX0] >= m.Registers[c.RVideoX1] {
		sx = -1
	}

	var sy int = 1
	if m.Registers[c.RVideoY0] >= m.Registers[c.RVideoY1] {
		sy = -1
	}

	for {
		m.putPixel(x, y, m.getVideoColour())
		if x == int(m.Registers[c.RVideoX1]) && y == int(m.Registers[c.RVideoY1]) {

		}

		var e2 int = 2 * err

		if e2 >= dy {
			if x == int(m.Registers[c.RVideoX1]) {
				break
			}

			err += dy
			x += sx
		}

		if e2 <= dx {
			if y == int(m.Registers[c.RVideoY1]) {
				break
			}

			err += dx
			y += sy
		}

	}

}

func (m *VM) drawText() {

	if m.TextBuffer[0] == 0 || m.getVideoColour()[3] == 0 {
		return
	}

	var textOffsetX = m.Registers[c.RVideoX0]
	var textOffsetY = m.Registers[c.RVideoY0]

	var colour [4]byte = m.getVideoColour()

	for _, char := range m.TextBuffer {

		if char == 0 {
			break
		} else if char == '\n' {
			textOffsetX = m.Registers[c.RVideoX0]
			textOffsetY += FontCharHeight + 1
			continue
		} else if char < ' ' || char > '~' { // Printable ASCII range
			char = 127
		}

		char -= 32

		for y := 0; y < int(FontCharHeight); y++ {
			for x := 0; x < int(FontCharWidth); x++ {

				if fontData[char][(y*int(FontCharWidth))+x] == 255 {
					m.putPixel(int(textOffsetX)+x, int(textOffsetY)+y, colour)
				}

			}
		}

		textOffsetX += FontCharWidth + 1

	}

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

func (m *VM) putPixel(x int, y int, colour [4]byte) {
	if x >= int(VideoBufferWidth) || x < 0 || y >= int(VideoBufferHeight) || y < 0 {
		return
	}

	var pixelAddr int = (x * int(VideoBytesPerPixel)) + (y * int(VideoBufferWidth) * int(VideoBytesPerPixel))

	if colour[3] == 255 {

		m.MemArray[pixelAddr] = colour[0]
		m.MemArray[pixelAddr+1] = colour[1]
		m.MemArray[pixelAddr+2] = colour[2]

	} else {

		var blendedColour [3]byte = blendPixel(colour, [3]byte(m.MemArray[pixelAddr:pixelAddr+3]))
		m.MemArray[pixelAddr] = blendedColour[0]
		m.MemArray[pixelAddr+1] = blendedColour[1]
		m.MemArray[pixelAddr+2] = blendedColour[2]

	}
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
	// Printable ASCII range
	if char < ' ' || char > '~' {
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
