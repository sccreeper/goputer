// Default frontend for goputer
package main

import (
	"fmt"
	"log"
	"os"
	c "sccreeper/govm/pkg/constants"
	"sccreeper/govm/pkg/vm"
)

var Name string = "GP32"
var Description string = "Default graphical front end"
var Authour string = "Oscar Peace (sccreeper)"
var Repository string = "https://github.com/sccreeper/goputer"

func Run(program []byte, args []string) {

	log.Println("GP32 frontend starting...")
	fmt.Println()

	var gp32 vm.VM
	var gp32_chan chan c.Interrupt = make(chan c.Interrupt, 128)

	vm.InitVM(&gp32, program, gp32_chan)

	go gp32.Run()

	for {

		if gp32.Finished {
			os.Exit(0)
		}

		select {
		case x := <-gp32_chan:

			switch x {
			case c.IntVideoText:
				fmt.Println(string(gp32.TextBuffer[:]))
				continue
			default:
				continue

			}

		}

	}

}
