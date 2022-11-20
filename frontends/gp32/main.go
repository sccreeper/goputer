// Default frontend for goputer
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sccreeper/goputer/frontends/gp32/colour"
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

	rand.Seed(time.Now().UnixNano())

	log.Println("GP32 frontend starting...")
	fmt.Println()

	var gp32 vm.VM
	var gp32_chan chan c.Interrupt = make(chan c.Interrupt)
	var gp32_subbed_chan chan c.Interrupt = make(chan c.Interrupt)
	var text_string string = ""

	var previous_mouse PreviousMousePos = PreviousMousePos{
		Button: 69,
	}

	vm.InitVM(&gp32, program, gp32_chan, gp32_subbed_chan)

	go gp32.Run()

	rl.InitWindow(640, 480, fmt.Sprintf("gp32 - %s", args[0]))
	rl.SetTargetFPS(128)

	sound.SoundInit()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		if gp32.Finished {
			break
		}

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

		//Handle inputs

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

		if uint32(rl.GetMouseX()) != previous_mouse.MouseX && uint32(rl.GetMouseY()) != previous_mouse.MouseY {

			gp32.Registers[c.RMouseX] = uint32(rl.GetMouseX())
			gp32.Registers[c.RMouseY] = uint32(rl.GetMouseY())

			previous_mouse.MouseX = uint32(rl.GetMouseX())
			previous_mouse.MouseY = uint32(rl.GetMouseY())

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

		rl.EndDrawing()
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
