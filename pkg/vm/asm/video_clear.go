//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	TEXT("VideoClearAsm", NOSPLIT, "func(array *byte, red uint8, green uint8, blue uint8)")

	ptr := Load(Param("array"), GP64())
	red := Load(Param("red"), GP8())
	green := Load(Param("green"), GP8())
	blue := Load(Param("blue"), GP8())
	max := GP64()
	MOVQ(ptr, max)
	ADDQ(Imm(230400), max)

	Label("loop")
	MOVB(red, Mem{Base: ptr})
	INCQ(ptr)
	MOVB(green, Mem{Base: ptr})
	INCQ(ptr)
	MOVB(blue, Mem{Base: ptr})
	INCQ(ptr)
	CMPQ(ptr, max)
	JBE(LabelRef("loop"))

	RET()

	Generate()

}
