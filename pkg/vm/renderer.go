package vm

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"iter"
	"math"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/gpimg"
	"sccreeper/goputer/pkg/util"
	"slices"
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

func (m *VM) drawArea() {

	// Do the first line

	if m.Registers[c.RVideoX0] > m.Registers[c.RVideoX1] || m.Registers[c.RVideoY0] > m.Registers[c.RVideoY1] {
		// TODO add some sort of panic
		return
	}

	var posX uint32 = util.Clamp(m.Registers[c.RVideoX0], 0, VideoBufferWidth)
	var posY uint32 = util.Clamp(m.Registers[c.RVideoY0], 0, VideoBufferHeight)
	var posX1 uint32 = util.Clamp(m.Registers[c.RVideoX1], 0, VideoBufferWidth)
	var posY1 uint32 = util.Clamp(m.Registers[c.RVideoY1], 0, VideoBufferHeight)

	var colour = m.getVideoColour()

	var areaStart = (uint32(posX) * VideoBytesPerPixel) + (uint32(posY) * VideoBufferWidth * VideoBytesPerPixel)

	if colour[3] == 255 {

		// Do the first row

		for x := 0; x < int(posX1-posX); x++ {
			var offset int = int(areaStart) + (x * int(VideoBytesPerPixel))

			m.MemArray[offset] = colour[0]
			m.MemArray[offset+1] = colour[1]
			m.MemArray[offset+2] = colour[2]
		}

		// Copy each row thereafter

		for y := 1; y < int(posY1-posY); y++ {
			var offset = int(areaStart) + (y * int(VideoBufferWidth) * int(VideoBytesPerPixel))

			copy(
				m.MemArray[offset:offset+(int(VideoBufferWidth*VideoBytesPerPixel))],
				m.MemArray[areaStart:areaStart+((uint32(posX1)-uint32(posX))*VideoBytesPerPixel)],
			)
		}

	} else {

		for y := 0; y < int(posY1-posY); y++ {
			for x := 0; x < int(posX1-posX); x++ {

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

	var colour [4]byte = m.getVideoColour()

	for v := range Bresenham(
		[2]int{int(m.Registers[c.RVideoX0]), int(m.Registers[c.RVideoY0])},
		[2]int{int(m.Registers[c.RVideoX1]), int(m.Registers[c.RVideoY1])},
	) {
		m.putPixel(v[0], v[1], colour)
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

	var numVertices int = int(m.DataBuffer[0])
	if numVertices <= 2 || 1+(numVertices*2) >= len(m.DataBuffer) {
		return
	}

	var vertices [][2]int = make([][2]int, 0, numVertices)
	var yMin int = math.MaxInt
	var yMax int = 0
	var xMin int = math.MaxInt
	var xMax int = 0

	for i := 1; i < numVertices*2; i += 2 {

		var adjustedX int = int(m.DataBuffer[i]) + int(m.Registers[c.RVideoX0])
		var adjustedY int = int(m.DataBuffer[i+1]) + int(m.Registers[c.RVideoY0])

		vertices = append(
			vertices, [2]int{
				adjustedX,
				adjustedY,
			},
		)

		if int(m.DataBuffer[i])+int(m.Registers[c.RVideoX0]) < xMin {
			xMin = adjustedX
		}

		if int(m.DataBuffer[i])+int(m.Registers[c.RVideoX0]) > xMax {
			xMax = adjustedX
		}

		if int(m.DataBuffer[i+1])+int(m.Registers[c.RVideoY0]) < yMin {
			yMin = adjustedY
		}

		if int(m.DataBuffer[i+1])+int(m.Registers[c.RVideoY0]) > yMax {
			yMax = adjustedY
		}

	}

	// Edge buckets
	// map[y int][]x int
	var edges map[int][]int = make(map[int][]int)

	// Calculate edges

	for i := 0; i < len(vertices); i++ {

		var next [2]int

		if i == len(vertices)-1 {
			next = vertices[0]
		} else {
			next = vertices[i+1]
		}

		for point := range Bresenham(vertices[i], next) {

			if _, keyExists := edges[point[1]]; keyExists {
				if !slices.Contains(edges[point[1]], point[0]) {
					edges[point[1]] = append(edges[point[1]], point[0])
				}
			} else {
				edges[point[1]] = make([]int, 0)
				edges[point[1]] = append(edges[point[1]], point[0])
			}

		}

	}

	// Finally fill in

	var colour [4]byte = m.getVideoColour()

	for y := yMin; y < yMax; y++ {
		var inShape bool = false

		if len(edges[y]) == 1 {
			m.putPixel(edges[y][0], y, colour)
			continue
		}

		for x := xMin; x < xMax; x++ {

			if inShape {
				m.putPixel(x, y, colour)
			}
			
			if slices.Contains(edges[y], x) && !slices.Contains(edges[y], x-1) {
				inShape = !inShape

				if inShape {
					m.putPixel(x, y, colour)
				}
			}
		
		}
	
	}

}

func (m *VM) drawImage() {

	var imgAddress uint32 = binary.LittleEndian.Uint32(m.DataBuffer[:4])

	var imgWidth int = int(binary.LittleEndian.Uint16(m.MemArray[imgAddress : imgAddress+2]))
	var imgHeight int = int(binary.LittleEndian.Uint16(m.MemArray[imgAddress+2 : imgAddress+4]))

	var imgFlags byte = m.MemArray[imgAddress+4]

	imgData, err := gpimg.Decode(m, int(m.Registers[c.RDataLength]), int(imgAddress))
	if err != nil {
		return
	}

	if imgFlags&gpimg.FlagBgOpaque != 0 {

		for y := 0; y < imgHeight; y++ {
			for x := 0; x < imgWidth; x++ {

				var offset int = (x * 4) + (y * imgWidth * 4)

				m.putPixel(
					x+int(m.Registers[c.RVideoX0]),
					y+int(m.Registers[c.RVideoY0]),
					[4]byte{
						imgData[offset],
						imgData[offset+1],
						imgData[offset+2],
						255,
					},
				)

			}
		}

	} else {

		for y := 0; y < imgHeight; y++ {
			for x := 0; x < imgWidth; x++ {

				var offset int = (x * 4) + (y * imgWidth * 4)

				m.putPixel(
					x+int(m.Registers[c.RVideoX0]),
					y+int(m.Registers[c.RVideoY0]),
					[4]byte(imgData[offset:offset+4]),
				)

			}
		}

	}

}

func (m *VM) clearVideo() {

	// Set first row

	var colour [4]byte = m.getVideoColour()

	for x := range int(VideoBufferWidth) {
		m.putPixel(x, 0, colour)
	}

	for y := 1; y < int(VideoBufferHeight); y++ {

		var offset int = y * int(VideoBufferWidth) * int(VideoBytesPerPixel)

		copy(
			m.MemArray[offset:offset+int(VideoBufferWidth)*int(VideoBytesPerPixel)],
			m.MemArray[:VideoBufferWidth*VideoBytesPerPixel],
		)

	}

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
		byte((int(src[0])*int(src[3]) + int(dest[0])*invertedAlpha) >> 8),
		byte((int(src[1])*int(src[3]) + int(dest[1])*invertedAlpha) >> 8),
		byte((int(src[2])*int(src[3]) + int(dest[2])*invertedAlpha) >> 8),
	}

}

func Bresenham(a [2]int, b [2]int) iter.Seq[[2]int] {

	var x int = int(a[0])
	var y int = int(a[1])

	var dx int = int(math.Abs(float64(b[0] - a[0])))
	var dy int = -int(math.Abs(float64(b[1] - a[1])))
	var err int = dx + dy

	var sx int = 1
	if a[0] >= b[0] {
		sx = -1
	}

	var sy int = 1
	if a[1] >= b[1] {
		sy = -1
	}

	return func(yield func([2]int) bool) {

		for {
			if !yield([2]int{x, y}) {
				return
			}

			if x == b[0] && y == b[1] {
				break
			}

			var e2 int = 2 * err

			if e2 >= dy {
				if x == b[0] {
					break
				}

				err += dy
				x += sx
			}

			if e2 <= dx {
				if y == b[1] {
					break
				}

				err += dx
				y += sy
			}

		}

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
