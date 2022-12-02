# Wrapper for py32 SO
import ctypes
from . import errors, constants
import os

print("Loading SO...")

files = os.listdir(".")

if "goputer" in files:
    _lib = ctypes.cdll.LoadLibrary("./frontends/goputerpy/goputerpy/bindings.so")
else:
    _lib = ctypes.cdll.LoadLibrary("./bindings.so")

_init = _lib.Init
_init.argtypes = [ctypes.Array, ctypes.c_int]
_init.restype = ctypes.c_void_p

_run = _lib.Run
_run.restype = ctypes.c_void_p

_get_interrupt = _lib.GetInterrupt
_get_interrupt.restype = ctypes.c_ulong

_send_interrupt = _lib.SendInterrupt
_send_interrupt.argtypes = [ctypes.c_ulong]
_send_interrupt.restype = ctypes.c_void_p

_get_buffer = _lib.GetBuffer
_get_buffer.argtypes = [ctypes.c_ulong]
_get_buffer.restype = ctypes.POINTER(ctypes.c_char * 128)

_set_register = _lib.SetRegister
_set_register.argtypes = [ctypes.c_ulong]
_set_register.restype = ctypes.c_void_p

_get_register = _lib.GetRegister
_get_register.argtypes = [ctypes.c_ulong]
_get_register.restype = ctypes.c_void_p

_free = _lib.free
_free.argtypes = [ctypes.c_void_p]

print("SO loaded!")

_vm_inited = False
_vm_alive = False

def Init(program_bytes: list) -> None:
    global _vm_inited
    a = (ctypes.c_char * len(program_bytes))(*program_bytes) 
    _vm_inited = True

    _init(a, ctypes.c_int(len(program_bytes)))
    

def Run() -> None:
    global _vm_alive

    if _vm_inited:
        _run()
        _vm_alive = True
    else:
        raise errors.VMNotInitialized("Must call Init() before running VM")

#Pops the last inerrupt off of the interupt array.
def GetInterrupt() -> constants.Interrupt:
    x = int(_get_interrupt())

    if x > constants.Interrupt.IntKeyboardDown:
        return None
    else:
        return constants.Interrupt(x)

def SendInterrupt(interrupt: constants.Interrupt):
    if type(interrupt) != constants.Interrupt:
        raise ValueError("Wrong type")
    elif not(interrupt > 0 and interrupt < constants.Interrupt.IntKeyboardDown):
        raise ValueError("Not valid interrupt")
    
    _send_interrupt(ctypes.c_ulong(interrupt))

def GetBuffer(b: constants.Register) -> list:
    if b == constants.Register.RVideoText or b == constants.Register.RData:
        a = _get_buffer(ctypes.c_ulong(b))

        l = [x for x in a.contents]
        _free(a)

        return l
    else:
        raise ValueError("Not a buffer!")

def SetRegister(r: constants.Register, v: int) -> None:
    _set_register(ctypes.c_ulong(r), ctypes.c_ulong(v))

def GetRegister(r: constants.Register) -> int:
    x = _get_register(ctypes.c_ulong(r))

    if x == None:
        x = 0

    return int(x)
