//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	// 230400 is divisible by 32 so no need to handle overflow
	TEXT("VideoClearAsm", NOSPLIT, "func(array *byte, red uint8, green uint8, blue uint8)")

	Comment("Load params")

	ptr := Load(Param("array"), GP64())
	red := Load(Param("red"), GP32())
	green := Load(Param("green"), GP32())
	blue := Load(Param("blue"), GP32())
	max := GP64()
	MOVQ(ptr, max)
	ADDQ(Imm(230400), max)

	Comment("Generate colour")

	ymm := YMM()
	xmm := XMM()
	tmp := GP32()

	MOVL(U32(0), tmp)
	MOVL(red, tmp)
	SHLL(Imm(8), tmp)
	ORL(green, tmp)
	SHLL(Imm(8), tmp)
	ORL(blue, tmp)
	SHLL(Imm(8), tmp)
	ORL(red, tmp)
	MOVD(tmp, xmm)
	VBROADCASTSS(xmm, ymm)

	Comment("Fill colour")

	Label("loop")
	VMOVDQU(ymm, Mem{Base: ptr})
	ADDQ(Imm(32), ptr)
	CMPQ(ptr, max)
	JBE(LabelRef("loop"))

	RET()

	Generate()

}
