import pygame as pg
from . import constants as c

def convert_colour(c: int) -> pg.Color:
    c = c.to_bytes(4, byteorder="little")

    return pg.Color(c[0], c[1], c[2], a=c[3])

def convert_to_hex(i: int) -> str:
    return f"0x{i:0>8X}"
