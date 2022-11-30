# Wrapper for py32 SO
import ctypes

_lib = ctypes.cdll.LoadLibrary("./bindings.so")

_init = _lib.Init
_init.argtypes = [ctypes.Array, ctypes.c_int]

_run = _lib.Run
_run.restype = ctypes.c_void_p

_get_interrupt = _lib.GetInterrupt
_get_interrupt.restype = ctypes.c_ulong

_send_interrupt = _lib.SendInterrupt
_send_interrupt.argtypes = [ctypes.c_ulong]

_get_buffer = _lib.GetBuffer
_get_buffer.argtypes = [ctypes.c_ulong]
_send_interrupt.restype = ctypes.Array

_set_register = _lib.SetRegister
_set_register.argtypes = [ctypes.c_ulong]

