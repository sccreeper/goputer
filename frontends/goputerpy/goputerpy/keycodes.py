import pygame as pg

_KEY_MAP: dict[int, int] = {
    pg.K_a: 1,
    pg.K_b: 2,
    pg.K_c: 3,
    pg.K_d: 4,
    pg.K_e: 5,
    pg.K_f: 6,
    pg.K_g: 7,
    pg.K_h: 8,
    pg.K_i: 9,
    pg.K_j: 10,
    pg.K_k: 11,
    pg.K_l: 12,
    pg.K_m: 13,
    pg.K_n: 14,
    pg.K_o: 15,
    pg.K_p: 16,
    pg.K_q: 17,
    pg.K_r: 18,
    pg.K_s: 19,
    pg.K_t: 20,
    pg.K_u: 21,
    pg.K_v: 22,
    pg.K_w: 23,
    pg.K_x: 24,
    pg.K_y: 25,
    pg.K_z: 26,

    pg.K_SPACE: 27,

    pg.K_1: 28,
    pg.K_2: 29,
    pg.K_3: 30,
    pg.K_4: 31,
    pg.K_5: 32,
    pg.K_6: 33,
    pg.K_7: 34,
    pg.K_8: 35,
    pg.K_9: 36,
    pg.K_0: 37,

    pg.K_MINUS: 38,
    pg.K_EQUALS: 39,

    pg.K_LEFTBRACKET: 40,
    pg.K_RIGHTBRACKET: 41,

    pg.K_SEMICOLON: 42,
    pg.K_QUOTE: 43,

    pg.K_COMMA: 44,
    pg.K_PERIOD: 45,

    pg.K_SLASH: 46,
    pg.K_BACKSLASH: 47,

    pg.K_LCTRL: 48,
    pg.K_RCTRL: 49,

    pg.K_LALT: 50,
    pg.K_RALT: 51,

    pg.K_LSHIFT: 52,
    pg.K_RSHIFT: 53,

    pg.K_CAPSLOCK: 54,
    pg.K_TAB: 55,
    pg.K_ESCAPE: 56,
    pg.K_BACKQUOTE: 57,

    pg.K_BACKSPACE: 58,
    pg.K_RETURN: 59,

    pg.K_LSUPER: 60,
    pg.K_RSUPER: 61,

    pg.K_MENU: 62,

    pg.K_PRINTSCREEN: 63,
    pg.K_SCROLLLOCK: 64,
    pg.K_PAUSE: 65,

    pg.K_INSERT: 66,
    pg.K_HOME: 67,
    pg.K_DELETE: 68,
    pg.K_END: 69,

    pg.K_PAGEUP: 70,
    pg.K_PAGEDOWN: 71,

    pg.K_UP: 72,
    pg.K_RIGHT: 73,
    pg.K_DOWN: 74,
    pg.K_LEFT: 75,

    pg.K_NUMLOCK: 76,
    pg.K_KP_DIVIDE: 77,
    pg.K_KP_MULTIPLY: 78,
    pg.K_KP_MINUS: 79,
    pg.K_KP_PLUS: 80,
    pg.K_KP_ENTER: 81,

    pg.K_KP1: 82,
    pg.K_KP2: 83,
    pg.K_KP3: 84,
    pg.K_KP4: 85,
    pg.K_KP5: 86,
    pg.K_KP6: 87,
    pg.K_KP7: 88,
    pg.K_KP8: 89,
    pg.K_KP9: 90,
    pg.K_KP0: 91,
    pg.K_KP_PERIOD: 92,

    pg.K_F1: 93,
    pg.K_F2: 94,
    pg.K_F3: 95,
    pg.K_F4: 96,
    pg.K_F5: 97,
    pg.K_F6: 98,
    pg.K_F7: 99,
    pg.K_F8: 100,
    pg.K_F9: 101,
    pg.K_F10: 102,
    pg.K_F11: 103,
    pg.K_F12: 104,
}

def map_keycode(pg_code: int) -> int:

    if not pg_code in _KEY_MAP:
        return 0
    else:
        return _KEY_MAP[pg_code]
