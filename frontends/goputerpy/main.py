import sys
import os

import pygame as pg
from pygame.mixer import pre_init

from goputerpy import goputerpy as gppy
from goputerpy import constants as c
from goputerpy import util
from goputerpy.sound import SoundManager

def correct_mouse_y(y: int) -> int:
    return y - r.TOTAL_Y_OFFSET


if not len(sys.argv) >= 2:
    print("Pass in a file to run!")
    exit()

#Sound init

pre_init(44100, -16, 1, 1024)

#pygame init
pg.init()
import rendering as r
from rendering.io import Switch, Light

size = width, height = 640, r.TOTAL_Y_OFFSET + 480

screen = pg.display.set_mode(size)
pg.display.set_caption(f"goputerpy - {os.path.basename(sys.argv[1])}")
font = pg.font.SysFont(None, 32)

video_surface = pg.surface.Surface((640, 480))
brightness_surface = pg.surface.Surface((640, 480))
io_surface = pg.surface.Surface((640, r.IO_UI_SIZE))
debug_surface = pg.surface.Surface((640, r.DEBUG_UI_SIZE))

#Read the code file
f_name = sys.argv[1]

print(f"Loading program {f_name}...")

f_bytes = open(f_name, "rb")
f_bytes = f_bytes.read()

print(f"Program size: {len(list(f_bytes))}")

sound_manager = SoundManager()

#goputer init
gppy.Init(list(f_bytes))

#io init

io_state = [False for i in range(16)]
io_lights: list[Light] = []
io_switches: list[Switch] = []

for i in range(8):
    io_lights.append(
        Light(i * r.IO_SWITCH_SIZE, i)
    )

for i in range(8):
    io_switches.append(
        Switch((i + 8) * r.IO_SWITCH_SIZE, i)
    )

video_text = ""

prev_mouse_pos = (0, 0)

clock = pg.time.Clock()


while True:

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

                # Handle newlines

                text_lines = video_text.splitlines()

                line_count = 0

                for l in text_lines:
                    txt_img = font.render(
                        l.replace("\x00", ""), 
                        True, 
                        util.convert_colour(gppy.GetRegister(c.Register.RVideoColour)),
                        )
                    
                    video_surface.blit(txt_img, (gppy.GetRegister(c.Register.RVideoX0), gppy.GetRegister(c.Register.RVideoY0)  + (line_count * 32)))

                    line_count += 1

        case c.Interrupt.IntVideoClear:
            video_surface.fill(util.convert_colour(gppy.GetRegister(c.Register.RVideoColour)))

        case c.Interrupt.IntVideoArea:
            pg.draw.rect(
            video_surface, 
            util.convert_colour(gppy.GetRegister(c.Register.RVideoColour)),
            pg.Rect(
                gppy.GetRegister(c.Register.RVideoX0),
                gppy.GetRegister(c.Register.RVideoY0),
                gppy.GetRegister(c.Register.RVideoX1) - gppy.GetRegister(c.Register.RVideoX0),
                gppy.GetRegister(c.Register.RVideoY1) - gppy.GetRegister(c.Register.RVideoY0),
            )
            )
        
        case c.Interrupt.IntVideoLine:
            pg.draw.line(
                video_surface,
                util.convert_colour(gppy.GetRegister(c.Register.RVideoColour)),
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
            video_surface.set_at(
                (gppy.GetRegister(c.Register.RVideoX0),
                gppy.GetRegister(c.Register.RVideoY0)),
                util.convert_colour(gppy.GetRegister(c.Register.RVideoColour)),
                )
        case c.Interrupt.IntSoundFlush:
            sound_manager.play(
                gppy.GetRegister(c.Register.RSoundTone),
                gppy.GetRegister(c.Register.RSoundVolume) / 255,
                c.SoundWave(gppy.GetRegister(c.Register.RSoundWave)),
            )
        case c.Interrupt.IntSoundStop:
            sound_manager.stop()
        case c.Interrupt.IntIOFlush:
            #Get IO states from registers.

            for i in range(len(io_state)):
                io_state[i] = True if gppy.GetRegister(c.Register(c.Register.RIO00 + i)) > 0 else False

    for event in pg.event.get():
        match event.type:
            case pg.QUIT: 
                sys.exit()
        
            case pg.MOUSEMOTION:
                if pg.mouse.get_pos()[0] != prev_mouse_pos[0] or pg.mouse.get_pos()[1] != prev_mouse_pos[1]:
                    gppy.SetRegister(c.Register.RMouseX, pg.mouse.get_pos()[0])
                    gppy.SetRegister(c.Register.RMouseY, correct_mouse_y(pg.mouse.get_pos()[1]))

                    prev_mouse_pos = pg.mouse.get_pos()

                    if gppy.IsSubscribed(c.Interrupt.IntMouseMove):
                        gppy.SendInterrupt(c.Interrupt.IntMouseMove)
            
            case pg.MOUSEBUTTONDOWN:
                #Check if in bounds of rendering area.
                if pg.mouse.get_pos()[1] > r.TOTAL_Y_OFFSET:
                    gppy.SetRegister(c.Register.RMouseButton, event.button)

                    if gppy.IsSubscribed(c.Interrupt.IntMouseDown):
                        gppy.SendInterrupt(c.Interrupt.IntMouseDown)
                
            case pg.MOUSEBUTTONUP:
                
                #IO Switch interrupts
                if event.button == 1:
                    for s in io_switches:
                        x = s.update(pg.mouse.get_pos())

                        gppy.SetRegister(c.Register((s._id + c.Register.RIO08)), 1 if s.enabled else 0)

                        if x:
                            if gppy.IsSubscribed(c.Interrupt(s._id + c.Interrupt.IntIO08)):
                                gppy.SendInterrupt(c.Interrupt(s._id + c.Interrupt.IntIO08))
                
                #Check if in bounds of rendering area.
                elif pg.mouse.get_pos()[1] > r.TOTAL_Y_OFFSET:
                    gppy.SetRegister(c.Register.RMouseButton, event.button)

                    if gppy.IsSubscribed(c.Interrupt.IntMouseUp):
                        gppy.SendInterrupt(c.Interrupt.IntMouseUp)
            
            case pg.KEYDOWN:
                gppy.SetRegister(c.Register.RKeyboardCurrent, event.key)

                if gppy.IsSubscribed(c.Interrupt.IntKeyboardDown):
                    gppy.SendInterrupt(c.Interrupt.IntKeyboardDown)

            case pg.KEYUP:
                gppy.SetRegister(c.Register.RKeyboardPressed, event.key)

                if gppy.IsSubscribed(c.Interrupt.IntKeyboardUp):
                    gppy.SendInterrupt(c.Interrupt.IntKeyboardDown)

            case _:
                continue
    
    #Video brightness

    alpha = 0

    if gppy.GetRegister(c.Register.RVideoBrightness) == 0:
        alpha = 255
    else:
        alpha = round((1 - pow(pow(gppy.GetRegister(c.Register.RVideoBrightness), -1) * 255.0, -1)) * 255)

    brightness_surface.set_alpha(alpha)
    brightness_surface.fill((0, 0, 0))

    video_surface.blit(brightness_surface, (0, 0))

    #Draw IO

    io_surface.fill(r.GREY)

    for i in range(len(io_lights)):
        
        io_lights[i].on = io_state[i]
        io_lights[i].draw(io_surface)

    for s in io_switches:
        s.draw(io_surface)

    #Draw debug UI

    prc_img = font.render(f"Program counter: {util.convert_to_hex(gppy.GetRegister(c.Register.RProgramCounter))}", True, r.WHITE)
    itn_img = font.render(f"Current instruction: {util.generate_instruction_string(gppy.GetInstruction(), gppy.GetLargeArg(), gppy.GetSmallArgs())}", True, r.WHITE)

    debug_surface.fill(r.GREY)

    debug_surface.blit(prc_img, (5, 0))
    debug_surface.blit(itn_img, (5, 16))

    screen.fill(r.BLACK)
    screen.blit(debug_surface, (0, 0))
    screen.blit(io_surface, (0, r.DEBUG_UI_SIZE + r.SEPERATOR_SIZE))
    screen.blit(video_surface, (0, r.TOTAL_Y_OFFSET))

    if gppy.IsFinished():
        pg.display.flip()
        break
    
    pg.display.flip()

    gppy.Cycle()

    clock.tick(120)

#Hang once finished executing code
while True:
    for event in pg.event.get():
        if event.type == pg.QUIT: 
            sys.exit()