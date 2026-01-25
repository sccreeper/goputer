// Default frontend for goputer
package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"sccreeper/goputer/frontends/gp32/keyboard"
	"sccreeper/goputer/frontends/gp32/rendering"
	"sccreeper/goputer/frontends/gp32/sound"
	c "sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/expansions"
	"sccreeper/goputer/pkg/profiler"
	"sccreeper/goputer/pkg/vm"
	"strconv"
	"time"

	"github.com/faiface/beep/speaker"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/cpu"
)

var Name string = "GP32"
var Description string = "Default graphical front end"
var Authour string = "Oscar Peace (sccreeper)"
var Repository string = "https://github.com/sccreeper/goputer"

var useProfiler bool
var profilerOut string

var programPath string

//To avoid double firing interrupts

type PreviousMousePos struct {
	MouseX uint32
	MouseY uint32
	Button uint32
}

func Run(ctx *cli.Context) error {

	log.Println("Reading file...")

	programBytes, err := os.ReadFile(programPath)
	if err != nil {
		log.Fatal(err)
	}

	//Init

	log.Println("GP32 frontend starting...")

	if cpu.X86.HasAVX2 {
		log.Printf("ArchVideoClear: %s", strconv.FormatBool(vm.HaveArchVideoClear))
		log.Printf("ArchVideoAreaNoAlpha: %s", strconv.FormatBool(vm.HaveArchVideoArea))
		log.Printf("ArchVideoAreaAlpha: %s", strconv.FormatBool(vm.HaveArchVideoAreaAlpha))
	}

	rl.InitWindow(640, 480+int32(rendering.TotalYOffset), fmt.Sprintf("gp32 - %s", programPath))
	defer rl.CloseWindow()

	var gp32 *vm.VM

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

	for i := 0; i < len(VideoIntermediate); i++ {
		VideoIntermediate[i].A = 255
	}

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

	gp32, _ = vm.NewVM(programBytes, expansions.ModuleExists, expansions.Interaction)
	expansions.LoadExpansions(gp32)

	var pr *profiler.Profiler

	if useProfiler {

		var err error
		pr, err = profiler.NewProfiler(gp32)
		if err != nil {
			panic(err)
		}

	}

	expansions.SetAttribute("goputer.sys", "name", []byte("gp32"))

	sound.SoundInit()

	var startTime int64 = time.Now().UnixMilli()
	var framesCompleted int = 0
	var shouldCycle bool = false

	toggleManualButton := rendering.NewButton("Toggle step", 560, 0, 80, 24, func() {
		gp32.Mutex.Lock()
		shouldCycle = !shouldCycle
		gp32.Mutex.Unlock()
	})

	cycleButton := rendering.NewButton("Cycle", 560, 24, 80, 24, func() {
		if !shouldCycle {
			gp32.Cycle()
		}
	})

	// Start running

	go func() {

		for {
			if shouldCycle {
				if useProfiler {
					pr.Cycle()
				} else {
					gp32.Cycle()
				}
			}
		}

	}()

	//Start rendering

	for !rl.WindowShouldClose() {

		if framesCompleted == 0 {
			gp32.Mutex.Lock()
			shouldCycle = true
			gp32.Mutex.Unlock()
		}

		//Render IO

		rl.BeginTextureMode(IOStatusRenderTexture)
		rendering.RenderIO(ioStatus[:], ioToggleSwitches[:])
		rl.EndTextureMode()

		rl.BeginTextureMode(VMStatusRenderTexture)
		rendering.RenderVMDebug(gp32)
		toggleManualButton.Draw()
		cycleButton.Draw()
		rl.EndTextureMode()

		//Handle interrupts

		for len(gp32.InterruptQueue) > 0 {

			gp32.Mutex.Lock()
			var x c.Interrupt
			x, gp32.InterruptQueue = gp32.InterruptQueue[0], gp32.InterruptQueue[1:]
			gp32.Mutex.Unlock()

			switch x {
			// Sound interrupts
			case c.IntSoundFlush:
				gp32.Mutex.Lock()
				sound.PlaySound(gp32.Registers[c.RSoundWave], gp32.Registers[c.RSoundTone], gp32.Registers[c.RSoundVolume])
				gp32.Mutex.Unlock()
			case c.IntSoundStop:
				speaker.Clear()
			case c.IntIOFlush:
				gp32.Mutex.Lock()
				for i := 0; i < 8; i++ {

					if gp32.Registers[i+int(c.RIO00)] != 0 {
						ioStatus[i] = true
					} else {
						ioStatus[i] = false
					}

				}
				gp32.Mutex.Unlock()
			case c.IntVideoFlush:
				// Update video texture

				gp32.Mutex.Lock()
				for i := 0; i < int(vm.VideoBufferSize/3); i++ {

					VideoIntermediate[i].R = gp32.MemArray[i*3]
					VideoIntermediate[i].G = gp32.MemArray[(i*3)+1]
					VideoIntermediate[i].B = gp32.MemArray[(i*3)+2]

				}
				gp32.Mutex.Unlock()

				rl.UpdateTexture(
					VideoRenderTexture.Texture,
					VideoIntermediate[:],
				)

			default:
				continue
			}
		}

		rl.BeginTextureMode(VideoRenderTexture)

		// Draw video brightness

		var b float64

		gp32.Mutex.Lock()
		if gp32.Registers[c.RVideoBrightness] == 0 {
			b = 0xFF
		} else {
			b = (1 - math.Pow(math.Pow(float64(gp32.Registers[c.RVideoBrightness]), -1)*255.0, -1)) * 255
		}
		gp32.Mutex.Unlock()

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

				if rl.IsKeyDown(key) {

					gp32.Mutex.Lock()

					gp32.Registers[c.RKeyboardCurrent] = uint32(keyboard.MapKey(key))
					gp32.Registers[c.RKeyboardPressed] = 1

					if gp32.Subscribed(c.IntKeyboardDown) {
						gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.IntKeyboardDown)
					}

					gp32.Mutex.Unlock()

				} else if rl.IsKeyUp(key) {

					gp32.Mutex.Lock()

					gp32.Registers[c.RKeyboardCurrent] = uint32(keyboard.MapKey(key))
					gp32.Registers[c.RKeyboardPressed] = 0

					if gp32.Subscribed(c.IntKeyboardUp) {
						gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.IntKeyboardUp)
					}

					gp32.Mutex.Unlock()

				} else {

					gp32.Mutex.Lock()

					gp32.Registers[c.RKeyboardCurrent] = uint32(keyboard.MapKey(key))
					gp32.Registers[c.RKeyboardPressed] = 0

					if gp32.Subscribed(c.IntKeyboardDown) {
						gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.IntKeyboardDown)
					}

					if gp32.Subscribed(c.IntKeyboardUp) {
						gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.IntKeyboardUp)
					}

					gp32.Mutex.Unlock()

				}
			} else {
				break
			}

		}

		//Mouse

		if uint32(rl.GetMouseX()) != previousMouse.MouseX && uint32(CorrectedMouseY()) != previousMouse.MouseY {

			gp32.Mutex.Lock()
			gp32.Registers[c.RMouseX] = uint32(rl.GetMouseX()) / 2
			gp32.Registers[c.RMouseY] = uint32(CorrectedMouseY())
			gp32.Mutex.Unlock()

			previousMouse.MouseX = uint32(rl.GetMouseX()) / 2
			previousMouse.MouseY = uint32(CorrectedMouseY())

			if gp32.Subscribed(c.IntMouseMove) {
				gp32.Mutex.Lock()
				gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.IntMouseMove)
				gp32.Mutex.Unlock()
			}

		}

		//Loop through buttons and check each one
		for i := 0; i < int(rl.MouseMiddleButton)+1; i++ {

			if rl.IsMouseButtonDown(rl.MouseButton(i)) && i != int(previousMouse.Button) {

				gp32.Mutex.Lock()
				gp32.Registers[c.RMouseButton] = uint32(i)
				gp32.Mutex.Unlock()
				previousMouse.Button = uint32(i)

				if gp32.Subscribed(c.IntMouseDown) {
					log.Println("Interrupt: Mouse down")

					gp32.Mutex.Lock()
					gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.IntMouseDown)
					gp32.Mutex.Unlock()

				}

			} else if rl.IsMouseButtonReleased(rl.MouseButton(i)) && i != int(previousMouse.Button) {
				gp32.Mutex.Lock()
				gp32.Registers[c.RMouseButton] = uint32(i)
				gp32.Mutex.Unlock()

				if gp32.Subscribed(c.IntMouseUp) {
					log.Println("Interrupt: Mouse up")

					gp32.Mutex.Lock()
					gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.IntMouseUp)
					gp32.Mutex.Unlock()
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
						gp32.Mutex.Lock()
						gp32.SubscribedInterruptQueue = append(gp32.SubscribedInterruptQueue, c.Interrupt(int(c.IntIO08)+index))
						gp32.Mutex.Unlock()
					}

				}

			}

			toggleManualButton.Update(rl.Vector2{
				X: float32(rl.GetMouseX()),
				Y: float32(rl.GetMouseY()),
			})

			cycleButton.Update(rl.Vector2{
				X: float32(rl.GetMouseX()),
				Y: float32(rl.GetMouseY()),
			})

		}

		rl.EndDrawing()

		framesCompleted++

		//Check if finished and then exit program loop

		if gp32.Finished {
			gp32.Mutex.Lock()
			shouldCycle = false
			gp32.Mutex.Unlock()
			break
		}

	}

	fmt.Printf("Frames completed: %d\n", framesCompleted)
	fmt.Printf("Time elapsed: %dms\n", time.Now().UnixMilli()-startTime)
	fmt.Printf("Mean time per frame: %fms\n", float64(time.Now().UnixMilli()-startTime)/float64(framesCompleted))

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.EndDrawing()
	}

	if useProfiler {
		f, err := os.OpenFile(
			filepath.Join(os.Getenv("GOPUTER_CWD"), ctx.String("profilerout")),
			os.O_WRONLY|os.O_CREATE,
			0600,
		)

		if err != nil {
			panic(err)
		}

		pr.Finish()
		pr.Dump(f)
	}

	return nil
}

func CorrectedMouseY() int32 {

	return (rl.GetMouseY() - int32(rendering.TotalYOffset)) / 2

}

func main() {

	if val, exists := os.LookupEnv("GOPUTER_ROOT"); !exists {
		log.Fatal("GOPUTER_ROOT not set")
	} else {
		log.Printf("GOPUTER_ROOT: %s\n", val)
	}

	app := &cli.App{
		Name:   "gp32 Frontend",
		Usage:  "golang frontend for goputer",
		Action: Run,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "executable",
				Aliases:     []string{"e"},
				Usage:       "Run `FILE`",
				Required:    true,
				Destination: &programPath,
			},
			&cli.BoolFlag{
				Name:        "useprofiler",
				Destination: &useProfiler,
			},
			&cli.StringFlag{
				Name:        "profilerout",
				Usage:       "Output to `FILE`",
				DefaultText: "out.gppr",
				Destination: &profilerOut,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
