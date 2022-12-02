import sys
from goputerpy import goputerpy as gppy
from goputerpy import constants as c
from goputerpy import util
import pygame as pg

#pygame init
pg.init()
size = width, height = 640, 480

screen = pg.display.set_mode(size)
font = pg.font.SysFont(None, 32)

#Read the code file
f_name = sys.argv[1]

f_bytes = open(f_name, "rb")
f_bytes = f_bytes.read()

#goputer init
gppy.Init(list(f_bytes))
gppy.Run()

video_text = ""

while True:

    for event in pg.event.get():
        if event.type == pg.QUIT: sys.exit()


    match gppy.GetInterrupt():
        case c.Interrupt.IntVideoText:
            #Decode string
            text = gppy.GetBuffer(c.Register.RVideoText)

            if [0] == 0:
                video_text = ""
            else:
                for b in text:
                    video_text += b.decode()
                
                txt_img = font.render(
                    video_text.replace("\x00", ""), 
                    True, 
                    util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour))
                    )
                
                screen.blit(txt_img, (0, 0))

        case c.Interrupt.IntVideoClear:
            screen.fill(util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour)))

        case c.Interrupt.IntVideoArea:
            pg.draw.rect(
            screen, 
            util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour)),
            pg.Rect(
                gppy.GetRegister(c.Register.RVideoX0),
                gppy.GetRegister(c.Register.RVideoY0),
                gppy.GetRegister(c.Register.RVideoX1),
                gppy.GetRegister(c.Register.RVideoY1)
            )
            )
        
        case c.Interrupt.IntVideoLine:
            pg.draw.line(
                screen,
                util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour)),
                pg.Vector2(
                    x=gppy.GetRegister(c.Register.RVideoX0),
                    y=gppy.GetRegister(c.Register.RVideoY0)
                    ),
                pg.Vector2(
                    x=gppy.GetRegister(c.Register.RVideoX1),
                    y=gppy.GetRegister(c.Register.RVideoY1),
                )
            )
        case c.Interrupt.IntVideoPixel:
            screen.set_at(
                (gppy.GetRegister(c.Register.RVideoX0),
                gppy.GetRegister(c.Register.RVideoY0)),
                util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour))
                )
        
        

    pg.display.update()
