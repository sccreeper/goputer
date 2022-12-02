import sys
from goputerpy import goputerpy
from goputerpy import constants as c

f_name = sys.argv[1]

f_bytes = open(f_name, "rb")
f_bytes = f_bytes.read()

goputerpy.Init(list(f_bytes))
goputerpy.Run()

while True:

    match goputerpy.GetInterrupt():
        case c.Interrupt.IntVideoText:
            #Decode string
            s = ""
            text = goputerpy.GetBuffer(c.Register.RVideoText)
            
            for b in text:
                s += b.decode()
            
            print(s)
