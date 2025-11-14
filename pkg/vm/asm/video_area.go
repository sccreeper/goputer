//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {

	TEXT("VideoAreaAsm", NOSPLIT, "func(array *byte, red uint8, green uint8, blue uint8, x uint32, y uint32, x1 uint32, y1 uint32)")

	shuffle_mask_low := GLOBL("shuffle_mask_low", RODATA|NOPTR)
	DATA(0, String([]byte{
		0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0,
	}))

	shuffle_mask_high := GLOBL("shuffle_mask_high", RODATA|NOPTR)
	DATA(0, String([]byte{
		1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1,
	}))

	Comment("Load parameters")

	ptr := Load(Param("array"), GP64())
	red := Load(Param("red"), GP8())
	green := Load(Param("green"), GP8())
	blue := Load(Param("blue"), GP8())

	x := Load(Param("x"), GP32())
	y := Load(Param("y"), GP32())
	x1 := Load(Param("x1"), GP32())
	y1 := Load(Param("y1"), GP32())

	colour := GP32()

	xmm_low := XMM()
	xmm_high := XMM()

	data := YMM()

	counter_x := GP32()
	counter_y := GP32()

	width := GP32()
	rows := GP32()

	tmp_32 := GP32()
	tmp_offset := GP64()
	tmp_operand := GP64()

	Comment("Offset pointer by x and y")
	MOVL(y, tmp_offset.As32())
	MOVQ(U64(960), tmp_operand)
	IMULQ(tmp_operand, tmp_offset)
	ADDQ(tmp_offset, ptr)

	MOVL(x, tmp_offset.As32())
	MOVQ(U64(3), tmp_operand)
	IMULQ(tmp_operand, tmp_offset)
	ADDQ(tmp_offset, ptr)

	Comment("Bounds")

	MOVL(y1, rows)
	SUBL(y, rows)

	MOVL(x1, width)
	SUBL(x, width)

	Comment("Construct colour")

	MOVB(blue, colour.As8())
	SHLL(Imm(8), colour)
	MOVB(green, colour.As8())
	SHLL(Imm(8), colour)
	MOVB(red, colour.As8())

	MOVD(colour, xmm_low)
	VPSHUFB(shuffle_mask_low, xmm_low, xmm_low)
	MOVD(colour, xmm_high)
	VPSHUFB(shuffle_mask_high, xmm_high, xmm_high)

	Comment("Insert data into each of the lanes in YMM")
	VINSERTF128(Imm(0), xmm_low, data, data)
	VINSERTF128(Imm(1), xmm_high, data, data)

	XORL(counter_x, counter_x)
	XORL(counter_y, counter_y)

	Comment("Loop to fill for no alpha")
	Label("na_loop")

	CMPL(width, Imm(10))
	JB(LabelRef("na_blit_remaining"))

	MOVQ(ptr, RDI)

	Label("na_loop_x")

	Comment("Check if less than 10 pixels to blit")

	MOVL(counter_x, tmp_32)
	ADDL(Imm(10), tmp_32)
	CMPL(tmp_32, width)
	JA(LabelRef("na_blit_remaining"))

	Comment("Otherwise blit 10 pixels at a time")

	VMOVDQU(data, Mem{Base: RDI})

	ADDQ(Imm(30), RDI)
	ADDL(Imm(10), counter_x)

	CMPL(counter_x, width)
	JBE(LabelRef("na_loop_x"))
	JA(LabelRef("na_loop_end"))

	Label("na_blit_remaining")

	MOVB(red, Mem{Base: RDI})
	MOVB(green, Mem{Base: RDI, Disp: 1})
	MOVB(blue, Mem{Base: RDI, Disp: 2})

	ADDQ(Imm(3), RDI)
	INCL(counter_x)

	CMPL(counter_x, width)
	JB(LabelRef("na_blit_remaining"))

	Label("na_loop_end")

	Comment("Cleanup")

	XORL(counter_x, counter_x)
	ADDQ(U32(960), ptr)
	INCL(counter_y)
	CMPL(counter_y, rows)
	JB(LabelRef("na_loop"))

	RET()

	Generate()

}
