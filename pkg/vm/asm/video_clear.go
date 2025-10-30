//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	// 230400 is divisible by 32 so no need to handle overflow
	TEXT("VideoClearAsm", NOSPLIT, "func(array *byte, red uint8, green uint8, blue uint8)")

	first_mask := GLOBL("first_mask", RODATA|NOPTR)
	DATA(0, String([]byte{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0}))

	second_mask := GLOBL("second_mask", RODATA|NOPTR)
	DATA(0, String([]byte{1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1}))

	third_mask := GLOBL("third_mask", RODATA|NOPTR)
	DATA(0, String([]byte{2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2}))

	Comment("Load params")

	ptr := Load(Param("array"), GP64())
	red := Load(Param("red"), GP8())
	green := Load(Param("green"), GP8())
	blue := Load(Param("blue"), GP8())
	max := GP64()
	MOVQ(ptr, max)
	ADDQ(Imm(230400), max)

	Comment("Generate colour")

	xmm0 := XMM()
	xmm1 := XMM()
	xmm2 := XMM()
	tmp := GP32()

	MOVL(U32(0), tmp)
	MOVB(blue, tmp.As8())
	SHLL(Imm(8), tmp)
	ORB(green, tmp.As8())
	SHLL(Imm(8), tmp)
	ORB(red, tmp.As8())

	Comment("Shuffle XMM")

	MOVD(tmp, xmm0)
	VPSHUFB(first_mask, xmm0, xmm0)
	VPSHUFB(second_mask, xmm0, xmm1)
	VPSHUFB(third_mask, xmm0, xmm2)

	Comment("Fill colour")

	Label("loop")

	MOVNTDQ(xmm0, Mem{Base: ptr})
	MOVNTDQ(xmm1, Mem{Base: ptr, Disp: 16})
	MOVNTDQ(xmm2, Mem{Base: ptr, Disp: 32})

	MOVNTDQ(xmm0, Mem{Base: ptr, Disp: 48})
	MOVNTDQ(xmm1, Mem{Base: ptr, Disp: 64})
	MOVNTDQ(xmm2, Mem{Base: ptr, Disp: 80})

	MOVNTDQ(xmm0, Mem{Base: ptr, Disp: 96})
	MOVNTDQ(xmm1, Mem{Base: ptr, Disp: 112})
	MOVNTDQ(xmm2, Mem{Base: ptr, Disp: 128})

	MOVNTDQ(xmm0, Mem{Base: ptr, Disp: 144})
	MOVNTDQ(xmm1, Mem{Base: ptr, Disp: 160})
	MOVNTDQ(xmm2, Mem{Base: ptr, Disp: 176})

	ADDQ(Imm(192), ptr)
	CMPQ(ptr, max)
	JBE(LabelRef("loop"))

	RET()

	Generate()

}
