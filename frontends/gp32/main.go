// Default frontend for goputer
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sccreeper/goputer/frontends/gp32/colour"
	"sccreeper/goputer/frontends/gp32/rendering"
	"sccreeper/goputer/frontends/gp32/sound"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/vm"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Name string = "GP32"
var Description string = "Default graphical front end"
var Authour string = "Oscar Peace (sccreeper)"
var Repository string = "https://github.com/sccreeper/goputer"

//To avoid double firing interrupts

type PreviousMousePos struct {
	MouseX uint32
	MouseY uint32
	Button uint32
}

func Run(program []byte, args []string) {

	//Init

	rand.Seed(time.Now().UnixNano())

	log.Println("GP32 frontend starting...")
	fmt.Println()

	rl.InitWindow(640, 480+int32(rendering.TotalYOffset), fmt.Sprintf("gp32 - %s", args[0]))
	rl.SetTargetFPS(128)

	var gp32 vm.VM
	var gp32_chan chan c.Interrupt = make(chan c.Interrupt)
	var gp32_subbed_chan chan c.Interrupt = make(chan c.Interrupt)
	var text_string string = ""

	var IO_status [16]bool = [16]bool{}
	var IOToggleSwitches [8]rendering.IOSwitch = [8]rendering.IOSwitch{}

	for index := range IOToggleSwitches {

		IOToggleSwitches[index] = rendering.IOSwitch{
			Toggled: false,
			ID:      uint32(index) + 8,
			X:       float32((index * rendering.IOUISize) + (8 * rendering.IOUISize)),
			Y:       0,
		}

	}

	var VideoRenderTexture rl.RenderTexture2D = rl.LoadRenderTexture(640, 480)
	var IOStatusRenderTexture rl.RenderTexture2D = rl.LoadRenderTexture(640, int32(rendering.IOUISize))
	var VMStatusRenderTexture rl.RenderTexture2D = rl.LoadRenderTexture(640, int32(rendering.DebugUISize))

	//Clear backgrounds of both textures

	rl.BeginTextureMode(VideoRenderTexture)
	rl.ClearBackground(rl.Black)
	rl.EndTextureMode()

	rl.BeginTextureMode(IOStatusRenderTexture)
	rl.ClearBackground(rl.Black)
	rl.EndTextureMode()

	rl.BeginTextureMode(VMStatusRenderTexture)
	rl.ClearBackground(rl.Black)
	rl.EndTextureMode()

	rendering.InitVMDebug()

	//Set mouse to arbitrary number so inputs register

	var previous_mouse PreviousMousePos = PreviousMousePos{
		Button: 69,
	}

	vm.InitVM(&gp32, program, gp32_chan, gp32_subbed_chan, false)

	go gp32.Run()

	sound.SoundInit()

	//Start rendering

	for !rl.WindowShouldClose() {

		//Render IO

		for i := 0; i < 8; i++ {

			if gp32.Registers[i+int(c.RIO00)] != 0 {
				IO_status[i] = true
			} else {
				IO_status[i] = false
			}

		}

		rl.BeginTextureMode(IOStatusRenderTexture)
		rendering.RenderIO(IO_status[:], IOToggleSwitches[:])
		rl.EndTextureMode()

		rl.BeginTextureMode(VMStatusRenderTexture)
		rendering.RenderVMDebug(&gp32)
		rl.EndTextureMode()

		//Check if finished and then exit program loop

		if gp32.Finished {
			break
		}

		//Handle interrupts

		rl.BeginTextureMode(VideoRenderTexture)

		select {
		case x := <-gp32_chan:
			switch x {
			//Video interrupts
			//TODO: Change colour to use colour in vc register
			case c.IntVideoText:
				if gp32.TextBuffer[0] == 0 {
					text_string = ""
				} else {
					str_temp := string(gp32.TextBuffer[:])
					str_temp = strings.ReplaceAll(str_temp, "\x00", "")
					text_string += strings.ReplaceAll(str_temp, `\n`, "\n")
					rl.DrawText(text_string, 0, 0, 16, colour.ConvertColour(gp32.Registers[c.RVideoColour]))
				}

			case c.IntVideoClear:
				rl.ClearBackground(colour.ConvertColour(gp32.Registers[c.RVideoColour]))
			case c.IntVideoPixel:
				rl.DrawPixel(
					int32(gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoY0]),
					colour.ConvertColour(gp32.Registers[c.RVideoColour]))
			case c.IntVideoLine:
				rl.DrawLine(
					int32(gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoY0]),
					int32(gp32.Registers[c.RVideoX1]),
					int32(gp32.Registers[c.RVideoY1]),
					colour.ConvertColour(gp32.Registers[c.RVideoColour]),
				)
			case c.IntVideoArea:
				rl.DrawRectangle(
					int32(gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoY0]),
					int32(gp32.Registers[c.RVideoX1]-gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoY1]-gp32.Registers[c.RVideoY0]),
					colour.ConvertColour(gp32.Registers[c.RVideoColour]),
				)
			case c.IntSoundFlush:
				sound.PlaySound(gp32.Registers[c.RSoundWave], gp32.Registers[c.RSoundTone], gp32.Registers[c.RSoundVolume])
			case c.IntSoundStop:
				speaker.Clear()

			}
		default:
		}

		rl.EndTextureMode()

		//Draw render textures to screen

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.DrawTexture(VMStatusRenderTexture.Texture, 0, 0, rl.White)

		rl.DrawTextureRec(
			VMStatusRenderTexture.Texture,
			rl.Rectangle{X: 0,
				Y:      0,
				Width:  640,
				Height: -float32(rendering.DebugUISize),
			},
			rl.Vector2{X: 0, Y: 0},
			rl.White,
		)

		rl.DrawTextureRec(
			IOStatusRenderTexture.Texture,
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  640,
				Height: -float32(rendering.IOUISize),
			},
			rl.Vector2{X: 0, Y: float32(rendering.DebugUISize)},
			rl.White,
		)

		rl.DrawLine(0, int32(rendering.IOUISize+rendering.DebugUISize+3), 640, int32(rendering.IOUISize+rendering.DebugUISize+3), rl.LightGray)

		rl.DrawTextureRec(VideoRenderTexture.Texture, rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  640,
			Height: -480,
		},
			rl.Vector2{X: 0, Y: float32(rendering.TotalYOffset)},
			rl.White,
		)

		//Handle subscribed interrupts

		var key int32

		//Keyboard
		for {

			key = rl.GetKeyPressed()

			if key != 0 {

				if rl.IsKeyDown(key) {
					gp32.Registers[c.RKeyboardCurrent] = uint32(key)
					gp32.Registers[c.RKeyboardPressed] = 1

					if gp32.Subscribed(c.IntKeyboardDown) {
						gp32_subbed_chan <- c.IntKeyboardDown
					}
				} else if rl.IsKeyUp(key) {
					gp32.Registers[c.RKeyboardCurrent] = uint32(key)
					gp32.Registers[c.RKeyboardPressed] = 0

					if gp32.Subscribed(c.IntKeyboardUp) {
						gp32_subbed_chan <- c.IntKeyboardUp
					}
				} else {

					gp32.Registers[c.RKeyboardCurrent] = uint32(key)

					if gp32.Subscribed(c.IntKeyboardDown) && gp32.Subscribed(c.IntKeyboardUp) {

						gp32_subbed_chan <- c.IntKeyboardDown
						gp32_subbed_chan <- c.IntKeyboardUp
					}

				}
			} else {
				break
			}

		}

		//Mouse

		if uint32(rl.GetMouseX()) != previous_mouse.MouseX && uint32(CorrectedMouseY()) != previous_mouse.MouseY {

			gp32.Registers[c.RMouseX] = uint32(rl.GetMouseX())
			gp32.Registers[c.RMouseY] = uint32(CorrectedMouseY())

			previous_mouse.MouseX = uint32(rl.GetMouseX())
			previous_mouse.MouseY = uint32(CorrectedMouseY())

			if gp32.Subscribed(c.IntMouseMove) {
				gp32_subbed_chan <- c.IntMouseMove
			}

		}

		//Loop through buttons and check each one
		for i := 0; i < rl.MouseMiddleButton+1; i++ {

			if rl.IsMouseButtonDown(int32(i)) && i != int(previous_mouse.Button) {

				gp32.Registers[c.RMouseButton] = uint32(i)
				previous_mouse.Button = uint32(i)

				if gp32.Subscribed(c.IntMouseDown) {
					gp32_subbed_chan <- c.IntMouseDown
				}

			} else if rl.IsMouseButtonReleased(int32(i)) && i != int(previous_mouse.Button) {
				gp32.Registers[c.RMouseButton] = uint32(i)

				if gp32.Subscribed(c.IntMouseUp) {
					gp32_subbed_chan <- c.IntMouseUp
				}

			}

		}

		//For updating IO

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

			for index := range IOToggleSwitches {

				if IOToggleSwitches[index].Update(rl.Vector2{
					X: float32(rl.GetMouseX()),
					Y: float32(rl.GetMouseY()),
				}) {

					if IOToggleSwitches[index].Toggled {
						gp32.Registers[int(c.RIO08)+index] = 1
					} else {
						gp32.Registers[int(c.RIO08)+index] = 0
					}

					if gp32.Subscribed(c.Interrupt(int(c.IntIO08) + index)) {

						gp32_subbed_chan <- c.Interrupt(int(c.IntIO08) + index)
					}

				}

			}

		}

		rl.EndDrawing()

	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func CorrectedMouseY() int32 {

	return rl.GetMouseY() - int32(rendering.TotalYOffset)

}
