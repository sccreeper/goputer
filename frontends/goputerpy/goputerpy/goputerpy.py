# Wrapper for py32 SO
import ctypes
from . import errors, constants
import os

print("Loading SO...")

files = os.listdir(".")

if "goputer" in files:
    _lib = ctypes.cdll.LoadLibrary("./frontends/goputerpy/bindings.so")
else:
    _lib = ctypes.cdll.LoadLibrary("./bindings.so")

_init = _lib.Init
_init.argtypes = [ctypes.Array, ctypes.c_int]
_init.restype = ctypes.c_void_p

_get_interrupt = _lib.GetInterrupt
_get_interrupt.restype = ctypes.c_uint32

_send_interrupt = _lib.SendInterrupt
_send_interrupt.argtypes = [ctypes.c_uint32]
_send_interrupt.restype = ctypes.c_void_p

_is_subscribed = _lib.IsSubscribed
_is_subscribed.argtypes = [ctypes.c_uint32]
_is_subscribed.restype = ctypes.c_uint32

_get_buffer = _lib.GetBuffer
_get_buffer.argtypes = [ctypes.c_uint32]
_get_buffer.restype = ctypes.POINTER(ctypes.c_char * 128)

_set_register = _lib.SetRegister
_set_register.argtypes = [ctypes.c_uint32]
_set_register.restype = ctypes.c_void_p

_get_register = _lib.GetRegister
_get_register.argtypes = [ctypes.c_uint32]
_get_register.restype = ctypes.c_void_p

_is_finished = _lib.IsFinished
_is_finished.restype = ctypes.c_uint32

_cycle = _lib.Cycle
_cycle.restype = ctypes.c_void_p

_free = _lib.free
_free.argtypes = [ctypes.c_void_p]

_get_arg = _lib.GetArgs
_get_arg.restype = ctypes.c_uint32

_get_current_instruction = _lib.GetCurrentInstruction
_get_current_instruction.restype = ctypes.c_uint32

_get_instruction_string = _lib.GetInstructionString
_get_instruction_string.restype = ctypes.c_void_p

_copy_video_buffer = _lib.CopyFrameBuffer
_copy_video_buffer.restype = ctypes.c_void_p

_set_exp_attribute = _lib.SetExpansionModuleAttribute
_set_exp_attribute.argtypes = [ctypes.c_void_p, ctypes.c_void_p, ctypes.c_void_p, ctypes.c_size_t]
_set_exp_attribute.restype = ctypes.c_void_p

_get_exp_attribute = _lib.GetExpansionModuleAttribute
_get_exp_attribute.argtypes = [ctypes.c_void_p, ctypes.c_void_p]
_get_exp_attribute.restype = ctypes.c_void_p

print("SO loaded!")

_vm_inited = False
_vm_alive = False

def Init(program_bytes: list) -> None:
    global _vm_inited
    a = (ctypes.c_char * len(program_bytes))(*program_bytes) 
    _vm_inited = True

    _init(a, ctypes.c_int(len(program_bytes)))


#Pops the last inerrupt off of the interupt array.
def GetInterrupt() -> constants.Interrupt:
    x = int(_get_interrupt())

    if x > constants.Interrupt.IntKeyboardDown:
        return None
    else:
        return constants.Interrupt(x)

def SendInterrupt(interrupt: constants.Interrupt):
    if type(interrupt) != constants.Interrupt:
        raise ValueError(f"Wrong type (type {type(interrupt)})")
    elif not(interrupt > 0 and interrupt <= constants.Interrupt.IntKeyboardDown):
        raise ValueError(f"Not valid interrupt {interrupt}")
    
    _send_interrupt(ctypes.c_uint32(interrupt))

def IsSubscribed(i: constants.Interrupt) -> bool:

    subscribed = int(_is_subscribed(ctypes.c_uint32(i)))

    return True if subscribed == 1 else False

def GetBuffer(b: constants.Register) -> list:
    if b == constants.Register.RVideoText or b == constants.Register.RData:
        a = _get_buffer(ctypes.c_uint32(b))

        l = [x for x in a.contents]
        _free(a)

        return l
    else:
        raise ValueError("Not a buffer!")

def SetRegister(r: constants.Register, v: int) -> None:
    if type(v) != int:
        raise TypeError(f"v should of type int (type {type(v)})")
    elif type(r) != constants.Register:
        raise TypeError(f"r should be of type Interrupt (type {type(r)})")

    _set_register(ctypes.c_uint32(r), ctypes.c_uint32(v))

def GetRegister(r: constants.Register) -> int:
    x = _get_register(ctypes.c_uint32(r))

    if x == None:
        x = 0

    return int(x)

def IsFinished() -> bool:
    finished = int(_is_finished())

    return True if finished == 1 else False

def Cycle() -> None:
    _cycle()

def GetLargeArg() -> int:

    return int(_get_arg())

def GetSmallArgs() -> tuple[2]:

    x = int(_get_arg())

    x = x.to_bytes(4, byteorder="little")

    return (
        int.from_bytes(x[0:2], byteorder="little"),
        int.from_bytes(x[0:2], byteorder="little")
    )

def GetInstruction() -> int:
    return int(_get_current_instruction())

def GetInstructionString() -> str:

    itn_ptr = _get_instruction_string()
    itn_bytes = ctypes.string_at(itn_ptr)
    _free(itn_ptr)

    return itn_bytes.decode("utf-8")

FRAMEBUFFER_SIZE = 320 * 240 * 3

video_buffer = None

def InitVideoBuffer() -> None:
    global video_buffer

    _lib.GetVideoBufferPtr.restype = ctypes.POINTER(ctypes.c_uint8)

    ptr = _lib.GetVideoBufferPtr()

    video_buffer = (ctypes.c_uint8 * FRAMEBUFFER_SIZE).from_address(ctypes.addressof(ptr.contents))

def UpdateVideoBuffer() -> None:
    _copy_video_buffer()

def SetExpAttribute(exp: str, attrib: str, val: bytes) -> None:
    _set_exp_attribute(exp.encode("utf-8"), attrib.encode("utf-8"), val, len(val))

def GetExpAttribute(exp: str, attrib: str) -> bytes:
    
    res_ptr = _get_exp_attribute(exp.encode("utf-8"), attrib.encode("utf-8"))
    res_bytes = ctypes.string_at(res_ptr)
    _free(res_ptr)

    return res_bytes