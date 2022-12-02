import sys
from goputerpy import goputerpy

f_name = sys.argv[1]

f_bytes = open(f_name, "rb")
f_bytes = f_bytes.read()

goputerpy.Init(f_bytes)
