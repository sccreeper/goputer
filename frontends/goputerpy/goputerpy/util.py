import pygame as pg

def ConvertColour(c: int) -> pg.Color:
    c = c.to_bytes(4, byteorder="little")

    return pg.Color(c[0], c[1], c[2], a=c[3])