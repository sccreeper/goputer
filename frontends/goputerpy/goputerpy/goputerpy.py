# Wrapper for py32 SO
import ctypes
from . import errors, constants
import os
import platform

from threading import Thread, Lock

mutex = Lock()
shouldRun = False

print("Loading SO...")

files = os.listdir(".")

lib_extension = ".dll" if platform.system() == "Windows" else ".so"

if "goputer" in files:
    _lib = ctypes.cdll.LoadLibrary(f"./frontends/goputerpy/bindings{lib_extension}")
else:
    _lib = ctypes.cdll.LoadLibrary(f"./bindings{lib_extension}")

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

_free = _lib.CFree
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
    mutex.acquire()
    x = int(_get_interrupt())
    mutex.release()

    if x > constants.Interrupt.IntKeyboardDown:
        return None
    else:
        return constants.Interrupt(x)

def SendInterrupt(interrupt: constants.Interrupt):
    if type(interrupt) != constants.Interrupt:
        raise ValueError(f"Wrong type (type {type(interrupt)})")
    elif not(interrupt > 0 and interrupt <= constants.Interrupt.IntKeyboardDown):
        raise ValueError(f"Not valid interrupt {interrupt}")
    
    mutex.acquire()
    _send_interrupt(ctypes.c_uint32(interrupt))
    mutex.release()

def IsSubscribed(i: constants.Interrupt) -> bool:

    mutex.acquire()
    subscribed = int(_is_subscribed(ctypes.c_uint32(i)))
    mutex.release()

    return True if subscribed == 1 else False

def GetBuffer(b: constants.Register) -> list:
    if b == constants.Register.RVideoText or b == constants.Register.RData:
        mutex.acquire()
        a = _get_buffer(ctypes.c_uint32(b))
        mutex.release()

        l = [x for x in a.contents]
        mutex.acquire()
        _free(a)
        mutex.release()

        return l
    else:
        raise ValueError("Not a buffer!")

def SetRegister(r: constants.Register, v: int) -> None:
    if type(v) != int:
        raise TypeError(f"v should of type int (type {type(v)})")
    elif type(r) != constants.Register:
        raise TypeError(f"r should be of type Interrupt (type {type(r)})")

    mutex.acquire()
    _set_register(ctypes.c_uint32(r), ctypes.c_uint32(v))
    mutex.release()

def GetRegister(r: constants.Register) -> int:
    mutex.acquire()
    x = _get_register(ctypes.c_uint32(r))
    mutex.release()

    if x == None:
        x = 0

    return int(x)

def IsFinished() -> bool:
    mutex.acquire()
    finished = int(_is_finished())
    mutex.release()

    return True if finished == 1 else False

def Cycle() -> None:
    _cycle()

def GetLargeArg() -> int:
    mutex.acquire()
    x = _get_arg()
    mutex.release()

    return int(x)

def GetSmallArgs() -> tuple[2]:

    mutex.acquire()
    x = int(_get_arg())
    mutex.release()

    x = x.to_bytes(4, byteorder="little")

    return (
        int.from_bytes(x[0:2], byteorder="little"),
        int.from_bytes(x[0:2], byteorder="little")
    )

def GetInstruction() -> int:
    mutex.acquire()
    x = _get_current_instruction()
    mutex.release()

    return x

def GetInstructionString() -> str:
    mutex.acquire()

    try:
        itn_ptr = _get_instruction_string()
        itn_bytes = ctypes.string_at(itn_ptr)
        _free(itn_ptr)

        return itn_bytes.decode("utf-8")
    finally:
        mutex.release()

FRAMEBUFFER_SIZE = 320 * 240 * 3

video_buffer = None

def InitVideoBuffer() -> None:
    global video_buffer

    _lib.GetVideoBufferPtr.restype = ctypes.POINTER(ctypes.c_uint8)

    ptr = _lib.GetVideoBufferPtr()

    video_buffer = (ctypes.c_uint8 * FRAMEBUFFER_SIZE).from_address(ctypes.addressof(ptr.contents))

def UpdateVideoBuffer() -> None:
    mutex.acquire()
    _copy_video_buffer()
    mutex.release()

def SetExpAttribute(exp: str, attrib: str, val: bytes) -> None:
    mutex.acquire()
    _set_exp_attribute(exp.encode("utf-8"), attrib.encode("utf-8"), val, len(val))
    mutex.release()

def GetExpAttribute(exp: str, attrib: str) -> bytes:
    
    mutex.acquire()
    res_ptr = _get_exp_attribute(exp.encode("utf-8"), attrib.encode("utf-8"))
    res_bytes = ctypes.string_at(res_ptr)
    _free(res_ptr)
    mutex.release()

    return res_bytes

def Run():
    def _run():
        while shouldRun:
            Cycle()

    p = Thread(target=_run)
    p.start()