# Handles rendering IO
from . import *
import pygame as pg

#Pretty much just a data class for handling IO lights
class Light:
    on: bool
    
    pos: int
    _id: int

    def __init__(self, pos: int, id: int) -> None:
        self.on = False
        self.pos = pos
        self._id = id

    def draw(self, surface: pg.Surface):
        pg.draw.rect(
            surface, 
            MIDDLE_GREY,
            pg.Rect(
                self.pos,
                0,
                IO_UI_SIZE,
                IO_UI_SIZE,
            )
        )

        pg.draw.rect(
            surface, 
            LIGHT_ON if self.on else BLACK,
            pg.Rect(
                self.pos + 5,
                0 + 5,
                IO_UI_SIZE - 10,
                IO_UI_SIZE - 10,
            )
        )


class Switch:
    enabled: bool
    
    pos: int
    _id: int

    def __init__(self, pos: int, id: int) -> bool:
        self.enabled = False
        self.pos = pos
        self._id = id

    def update(self, click_pos: tuple) -> None:
        if (
            (click_pos[0] > self.pos and click_pos[0] < self.pos + IO_UI_SIZE) 
            and 
            (click_pos[1] > DEBUG_UI_SIZE + SEPERATOR_SIZE and click_pos[1] < DEBUG_UI_SIZE + SEPERATOR_SIZE + IO_SWITCH_SIZE)
            ):

            self.enabled = not self.enabled
            return True
        else:
            return False

    def draw(self, surface: pg.Surface):
        #Big switch background
        
        pg.draw.rect(
            surface, 
            MIDDLE_GREY,
            pg.Rect(
                self.pos,
                0,
                IO_UI_SIZE,
                IO_UI_SIZE,
            )
        )

        #Switch background

        pg.draw.rect(
            surface, 
            DARK_GREY,
            pg.Rect(
                self.pos + 5,
                5,
                IO_UI_SIZE - 10,
                IO_UI_SIZE - (IO_SWITCH_SIZE / 2),
            )
        )

        #Switch click bit

        pg.draw.rect(
            surface,
            WHITE,
            pg.Rect(
                self.pos + 5 + (IO_SWITCH_SIZE / 2) if self.enabled else self.pos + 5,
                5,
                IO_SWITCH_SIZE - 10 if self.enabled else IO_SWITCH_SIZE - 5 - (IO_SWITCH_SIZE / 2),
                IO_UI_SIZE - (IO_SWITCH_SIZE / 2),
            )
        )

        txt_img = IO_SWITCH_FONT.render("On" if self.enabled else "Off", True, WHITE)

        surface.blit(txt_img, (self.pos + 10, (IO_SWITCH_SIZE / 2) + 5))

        

        