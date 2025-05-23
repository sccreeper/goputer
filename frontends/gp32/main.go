// Default frontend for goputer
package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"sccreeper/goputer/frontends/gp32/rendering"
	"sccreeper/goputer/frontends/gp32/sound"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/vm"
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

	log.Println("GP32 frontend starting...")
	fmt.Println()

	rl.InitWindow(640, 480+int32(rendering.TotalYOffset), fmt.Sprintf("gp32 - %s", args[0]))

	var gp32 vm.VM

	var ioStatus [16]bool = [16]bool{}
	var ioToggleSwitches [8]rendering.IOSwitch = [8]rendering.IOSwitch{}

	for index := range ioToggleSwitches {

		ioToggleSwitches[index] = rendering.IOSwitch{
			Toggled: false,
			ID:      uint32(index) + 8,
			X:       float32((index * rendering.IOUISize) + (8 * rendering.IOUISize)),
			Y:       0,
		}

	}

	var VideoRenderTexture rl.RenderTexture2D = rl.LoadRenderTexture(320, 240)
	var VideoIntermediate [vm.VideoBufferWidth * vm.VideoBufferHeight]color.RGBA = [vm.VideoBufferWidth * vm.VideoBufferHeight]color.RGBA{}
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

	var previousMouse PreviousMousePos = PreviousMousePos{
		Button: 64,
	}

	vm.InitVM(&gp32, program, true)

	expansions.SetAttribute("goputer.sys", "name", "gp32")

	sound.SoundInit()

	var startTime int64 = time.Now().UnixMilli()
	var cyclesCompleted int = 0

	//Start rendering

	for !rl.WindowShouldClose() {

		gp32.Cycle()

		//Render IO

		rl.BeginTextureMode(IOStatusRenderTexture)
		rendering.RenderIO(ioStatus[:], ioToggleSwitches[:])
		rl.EndTextureMode()

		rl.BeginTextureMode(VMStatusRenderTexture)
		rendering.RenderVMDebug(&gp32)
		rl.EndTextureMode()

		//Check if finished and then exit program loop

		if gp32.Finished {
			break
		}

		//Handle interrupts

		for len(gp32.InterruptQueue) > 0 {

			var x c.Interrupt
			x, gp32.InterruptQueue = gp32.InterruptQueue[0], gp32.InterruptQueue[1:]

			switch x {
			// Sound interrupts
			case c.IntSoundFlush:
				sound.PlaySound(gp32.Registers[c.RSoundWave], gp32.Registers[c.RSoundTone], gp32.Registers[c.RSoundVolume])
			case c.IntSoundStop:
				speaker.Clear()
			case c.IntIOFlush:
				for i := 0; i < 8; i++ {

					if gp32.Registers[i+int(c.RIO00)] != 0 {
						ioStatus[i] = true
					} else {
						ioStatus[i] = false
					}

				}
			default:
				continue
			}
		}

		// Update video texture

		for i := 0; i < int(vm.VideoBufferSize); i += 3 {

			VideoIntermediate[i/3] = color.RGBA{
				gp32.MemArray[i],
				gp32.MemArray[i+1],
				gp32.MemArray[i+2],
				255,
			}

		}

		rl.UpdateTexture(
			VideoRenderTexture.Texture,
			VideoIntermediate[:],
		)

		rl.BeginTextureMode(VideoRenderTexture)

		// Draw video brightness

		var b float64

		if gp32.Registers[c.RVideoBrightness] == 0 {
			b = 0xFF
		} else {
			b = (1 - math.Pow(math.Pow(float64(gp32.Registers[c.RVideoBrightness]), -1)*255.0, -1)) * 255
		}

		rl.DrawRectangle(
			0,
			0,
			int32(vm.VideoBufferWidth),
			int32(vm.VideoBufferHeight),
			rl.Color{
				R: 0,
				G: 0,
				B: 0,
				A: uint8(b),
			},
		)

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

		rl.DrawTexturePro(
			VideoRenderTexture.Texture,
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  float32(VideoRenderTexture.Texture.Width),
				Height: float32(VideoRenderTexture.Texture.Height),
			},
			rl.Rectangle{
				X:      0,
				Y:      float32(rendering.TotalYOffset),
				Width:  640,
				Height: -480,
			},
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)

		//Handle subscribed interrupts

		var key int32

		//Keyboard
		for {

			key = rl.GetKeyPressed()

			if key != 0 {

				log.Println("Interrupt: Key pressed")

				if rl.IsKeyDown(key) {

					gp32.Registers[c.RKeyboardCurrent] = uint32(key)
					gp32.Registers[c.RKeyboardPressed] = 1

					if gp32.Subscribed(c.IntKeyboardDown) {
						log.Println("Interrupt: Key down")

						gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.IntKeyboardDown)
					}
				} else if rl.IsKeyUp(key) {
					log.Println("Interrupt: Bozo")

					gp32.Registers[c.RKeyboardCurrent] = uint32(key)
					gp32.Registers[c.RKeyboardPressed] = 0

					if gp32.Subscribed(c.IntKeyboardUp) {
						log.Println("Interrupt: Key up")

						gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.IntKeyboardUp)
					}
				} else {

					log.Println("Interrupt: Triggering ku, kd")

					gp32.Registers[c.RKeyboardCurrent] = uint32(key)

					if gp32.Subscribed(c.IntKeyboardDown) {
						gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.IntKeyboardDown)
					}

					if gp32.Subscribed(c.IntKeyboardUp) {
						gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.IntKeyboardUp)
					}

				}
			} else {
				break
			}

		}

		//Mouse

		if uint32(rl.GetMouseX()) != previousMouse.MouseX && uint32(CorrectedMouseY()) != previousMouse.MouseY {

			gp32.Registers[c.RMouseX] = uint32(rl.GetMouseX()) / 2
			gp32.Registers[c.RMouseY] = uint32(CorrectedMouseY())

			previousMouse.MouseX = uint32(rl.GetMouseX()) / 2
			previousMouse.MouseY = uint32(CorrectedMouseY())

			if gp32.Subscribed(c.IntMouseMove) {
				gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.IntMouseMove)
			}

		}

		//Loop through buttons and check each one
		for i := 0; i < int(rl.MouseMiddleButton)+1; i++ {

			if rl.IsMouseButtonDown(rl.MouseButton(i)) && i != int(previousMouse.Button) {

				gp32.Registers[c.RMouseButton] = uint32(i)
				previousMouse.Button = uint32(i)

				if gp32.Subscribed(c.IntMouseDown) {
					log.Println("Interrupt: Mouse down")

					gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.IntMouseDown)
				}

			} else if rl.IsMouseButtonReleased(rl.MouseButton(i)) && i != int(previousMouse.Button) {
				gp32.Registers[c.RMouseButton] = uint32(i)

				if gp32.Subscribed(c.IntMouseUp) {
					log.Println("Interrupt: Mouse up")

					gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.IntMouseUp)
				}

			}

		}

		//For updating IO

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

			for index := range ioToggleSwitches {

				if ioToggleSwitches[index].Update(rl.Vector2{
					X: float32(rl.GetMouseX()),
					Y: float32(rl.GetMouseY()),
				}) {

					if ioToggleSwitches[index].Toggled {
						gp32.Registers[int(c.RIO08)+index] = 1
					} else {
						gp32.Registers[int(c.RIO08)+index] = 0
					}

					if gp32.Subscribed(c.Interrupt(int(c.IntIO08) + index)) {
						gp32.SubbedInterruptQueue = append(gp32.SubbedInterruptQueue, c.Interrupt(int(c.IntIO08)+index))
					}

				}

			}

		}

		rl.EndDrawing()

		cyclesCompleted++

	}

	fmt.Printf("Cycles completed: %d\n", cyclesCompleted)
	fmt.Printf("Time elapsed: %dms\n", time.Now().UnixMilli()-startTime)
	fmt.Printf("Mean time per cycle: %fms\n", float64(time.Now().UnixMilli()-startTime)/float64(cyclesCompleted))

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func CorrectedMouseY() int32 {

	return (rl.GetMouseY() - int32(rendering.TotalYOffset)) / 2

}
