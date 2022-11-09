// Default frontend for goputer
package main

import (
	"fmt"
	"log"
	"math/rand"
	c "sccreeper/govm/pkg/constants"
	"sccreeper/govm/pkg/vm"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Name string = "GP32"
var Description string = "Default graphical front end"
var Authour string = "Oscar Peace (sccreeper)"
var Repository string = "https://github.com/sccreeper/goputer"

func Run(program []byte, args []string) {

	rand.Seed(time.Now().UnixNano())

	log.Println("GP32 frontend starting...")
	fmt.Println()

	var gp32 vm.VM
	var gp32_chan chan c.Interrupt = make(chan c.Interrupt)
	var text_string string = ""

	vm.InitVM(&gp32, program, gp32_chan)

	go gp32.Run()

	rl.InitWindow(640, 480, fmt.Sprintf("gp32 - %s", "Program Name"))
	rl.SetTargetFPS(128)

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
				str_temp := string(gp32.TextBuffer[:])
				str_temp = strings.ReplaceAll(str_temp, "\x00", "")
				text_string += strings.ReplaceAll(str_temp, `\n`, "\n")
				rl.DrawText(text_string, 0, 0, 16, rl.White)
			case c.IntVideoClear:
				rl.ClearBackground(rl.Black)
			case c.IntVideoPixel:
				rl.DrawPixel(
					int32(gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoX1]),
					rl.White)
			case c.IntVideoLine:
				rl.DrawLine(
					int32(gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoY0]),
					int32(gp32.Registers[c.RVideoX1]),
					int32(gp32.Registers[c.RVideoY1]),
					rl.White,
				)
			case c.IntVideoArea:
				rl.DrawRectangle(
					int32(gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoY0]),
					int32(gp32.Registers[c.RVideoX1]-gp32.Registers[c.RVideoX0]),
					int32(gp32.Registers[c.RVideoY1]-gp32.Registers[c.RVideoY0]),
					rl.White,
				)
			}
		default:
		}

		rl.EndDrawing()
	}

	for !rl.WindowShouldClose() {
	}

	rl.CloseWindow()
}
