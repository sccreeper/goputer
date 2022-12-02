import sys
from goputerpy import goputerpy
from goputerpy import constants as c

f_name = sys.argv[1]

f_bytes = open(f_name, "rb")
f_bytes = f_bytes.read()

goputerpy.Init(list(f_bytes))
goputerpy.Run()

while True:
    print(str(goputerpy.GetBuffer(c.Register.RVideoText)))
    continue
