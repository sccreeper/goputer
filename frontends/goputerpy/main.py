import sys
from goputerpy import goputerpy as gppy
from goputerpy import constants as c
from goputerpy import util
import pygame as pg
from pygame.mixer import get_init, pre_init
from goputerpy.sound import SoundManager

#Sound init

pre_init(44100, -16, 1, 1024)
sound_manager = SoundManager()

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

    #Handle called interrupts

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
                    util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour)),
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
                gppy.GetRegister(c.Register.RVideoY1),
            )
            )
        
        case c.Interrupt.IntVideoLine:
            pg.draw.line(
                screen,
                util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour)),
                pg.Vector2(
                    x=gppy.GetRegister(c.Register.RVideoX0),
                    y=gppy.GetRegister(c.Register.RVideoY0),
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
                util.ConvertColour(gppy.GetRegister(c.Register.RVideoColour)),
                )
        case c.Interrupt.IntSoundFlush:
            sound_manager.play(
                gppy.GetRegister(c.Register.RSoundTone),
                gppy.GetRegister(c.Register.RSoundVolume) / 255,
                c.SoundWave(gppy.GetRegister(c.Register.RSoundWave)),
            )
        case c.Interrupt.IntSoundStop:
            sound_manager.stop()
        
    pg.display.update()
