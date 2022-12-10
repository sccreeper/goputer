import pygame as pg
from . import constants as c

def convert_colour(c: int) -> pg.Color:
    c = c.to_bytes(4, byteorder="little")

    return pg.Color(c[0], c[1], c[2], a=c[3])

def generate_instruction_string(itn: int, large_arg: int, args: list[2]) -> str:

    if itn + large_arg + args[0] + args[1] == 0:
        return "Finished"


    itn_string = ""

    itn_string += c.InstructionStrings[itn]

    if itn in c.SingleArgInstructions:
        
        if itn == c.Instructions["jmp"] or itn == c.Instructions["cndjmp"] or itn == c.Instructions["call"] or itn == c.Instructions["cndcall"] or itn == c.Instructions["lda"] or itn == c.Instructions["sta"]:
            itn_string += f" {large_arg}"
        elif itn == c.Instructions["int"]:
            itn_string += f" {c.InterruptStrings[large_arg]}"
        else:
            itn_string += f" {c.RegisterStrings[large_arg]}"

    else:
        itn_string += f" {c.RegisterStrings[args[0]]} {c.RegisterStrings[args[1]]}"
    
    return itn_string
