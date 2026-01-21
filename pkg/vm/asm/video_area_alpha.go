//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {

	TEXT("VideoAreaAlphaAsm", NOSPLIT, "func(array *byte, red uint8, green uint8, blue uint8, alpha uint8, x uint32, y uint32, x1 uint32, y1 uint32)")

	colour_shuffle_mask := GLOBL("shuffle_mask", RODATA|NOPTR)
	DATA(0, String([]byte{
		0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0,
	}))

	Comment("Load parameters")

	ptr := Load(Param("array"), GP64())
	red := Load(Param("red"), GP8())
	red_16 := GP16()
	green := Load(Param("green"), GP8())
	green_16 := GP16()
	blue := Load(Param("blue"), GP8())
	blue_16 := GP16()
	alpha := Load(Param("alpha"), GP8())
	inv_alpha := GP8()
	alpha_16 := GP16()
	inv_alpha_16 := GP16()
	alpha_values := XMM()
	widened_alpha_values := YMM()
	inverted_alpha_values := XMM()
	widened_inverted_alpha_values := YMM()

	x := Load(Param("x"), GP32())
	y := Load(Param("y"), GP32())
	x1 := Load(Param("x1"), GP32())
	y1 := Load(Param("y1"), GP32())

	colour := GP32()

	shuffled_colour_bytes := XMM()
	widened_shuffled_colour_bytes := YMM()

	memory_data := XMM()
	widened_memory_data := YMM()

	counter_x := GP32()
	counter_y := GP32()

	width := GP32()
	rows := GP32()

	last_byte := GP32()
	XORL(last_byte, last_byte)
	tmp_16 := GP16()
	tmp_32 := GP32()
	tmp_offset := GP64()
	tmp_operand := GP64()
	tmp_ymm := YMM()
	VPXOR(tmp_ymm, tmp_ymm, tmp_ymm)

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

	Comment("Fill alpha register")
	MOVBLZX(alpha, tmp_32)
	MOVD(tmp_32, alpha_values)
	VPBROADCASTB(alpha_values, alpha_values)
	VPMOVZXBW(alpha_values, widened_alpha_values)

	MOVB(alpha, inv_alpha)
	NOTB(inv_alpha)
	MOVBWZX(inv_alpha, inv_alpha_16)
	MOVBLZX(inv_alpha, tmp_32)
	MOVD(tmp_32, inverted_alpha_values)
	VPBROADCASTB(inverted_alpha_values, inverted_alpha_values)
	VPMOVZXBW(inverted_alpha_values, widened_inverted_alpha_values)

	Comment("Construct colour")

	MOVB(blue, colour.As8())
	SHLL(Imm(8), colour)
	MOVB(green, colour.As8())
	SHLL(Imm(8), colour)
	MOVB(red, colour.As8())

	MOVD(colour, shuffled_colour_bytes)
	VPSHUFB(colour_shuffle_mask, shuffled_colour_bytes, shuffled_colour_bytes)
	VPMOVZXBW(shuffled_colour_bytes, widened_shuffled_colour_bytes)

	VPMULLW(widened_shuffled_colour_bytes, widened_alpha_values, widened_shuffled_colour_bytes)

	XORL(counter_x, counter_x)
	XORL(counter_y, counter_y)

	Comment("Construct colour for remaining")
	MOVBWZX(red, red_16)
	MOVBWZX(green, green_16)
	MOVBWZX(blue, blue_16)
	MOVBWZX(alpha, alpha_16)

	IMULW(alpha_16, red_16)
	IMULW(alpha_16, green_16)
	IMULW(alpha_16, blue_16)

	Comment("Loop to fill for alpha")
	Label("a_loop")

	MOVQ(ptr, RDI)

	CMPL(width, Imm(5))
	JB(LabelRef("a_blit_remaining"))

	Label("a_loop_x")

	Comment("Check if less than 5 pixels to blit")

	MOVL(counter_x, tmp_32)
	ADDL(Imm(5), tmp_32)
	CMPL(tmp_32, width)
	JA(LabelRef("a_blit_remaining"))

	Comment("Otherwise blit 5 pixels at a time")

	Comment("Load data to be modified from memory")
	VMOVDQU(Mem{Base: RDI}, memory_data)
	XORL(last_byte, last_byte)
	PEXTRB(Imm(15), memory_data, last_byte)

	Comment("Modify memory data")
	VPMOVZXBW(memory_data, widened_memory_data)
	VPMULLW(widened_memory_data, widened_inverted_alpha_values, widened_memory_data)

	VPADDW(widened_memory_data, widened_shuffled_colour_bytes, widened_memory_data)

	VPSRLW(Imm(8), widened_memory_data, widened_memory_data)

	VPACKUSWB(widened_memory_data, widened_memory_data, widened_memory_data)
	VPERMQ(Imm(0xd8), widened_memory_data, widened_memory_data)
	VEXTRACTI128(Imm(0), widened_memory_data, memory_data)

	Comment("Move data back to memory")

	PINSRB(Imm(15), last_byte, memory_data)
	VMOVDQU(memory_data, Mem{Base: RDI})

	ADDQ(Imm(15), RDI)
	ADDL(Imm(5), counter_x)

	CMPL(counter_x, width)
	JBE(LabelRef("a_loop_x"))
	JA(LabelRef("a_loop_end"))

	Label("a_blit_remaining")

	MOVBWZX(Mem{Base: RDI}, tmp_16)
	IMULW(inv_alpha_16, tmp_16)
	ADDW(red_16, tmp_16)
	SHRW(Imm(8), tmp_16)
	MOVB(tmp_16.As8L(), Mem{Base: RDI})

	MOVBWZX(Mem{Base: RDI, Disp: 1}, tmp_16)
	IMULW(inv_alpha_16, tmp_16)
	ADDW(green_16, tmp_16)
	SHRW(Imm(8), tmp_16)
	MOVB(tmp_16.As8L(), Mem{Base: RDI, Disp: 1})

	MOVBWZX(Mem{Base: RDI, Disp: 2}, tmp_16)
	IMULW(inv_alpha_16, tmp_16)
	ADDW(blue_16, tmp_16)
	SHRW(Imm(8), tmp_16)
	MOVB(tmp_16.As8L(), Mem{Base: RDI, Disp: 2})

	ADDQ(Imm(3), RDI)
	INCL(counter_x)

	CMPL(counter_x, width)
	JB(LabelRef("a_blit_remaining"))

	Comment("Cleanup")
	Label("a_loop_end")

	XORL(counter_x, counter_x)
	ADDQ(U32(960), ptr)
	INCL(counter_y)
	CMPL(counter_y, rows)
	JB(LabelRef("a_loop"))

	RET()

	Generate()

}
