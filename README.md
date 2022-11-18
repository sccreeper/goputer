# goputer
<sup>`Go + Computer = goputer`</sup>

---

A computer emulator/virtual machine intended to demonstrate how computers work at a low level.

---

### Features

**Note:** Most of these have not been fully implemented yet.

###### Complete

- Custom assembly language and compiler.
- Custom runtime.
- Standalone executables.

###### Working on

- Frontend that shows output from VM backend.

###### In the future

- WASM based frontend.
- High level language.
- IDE for easy development.

---

### Documentation

See the [project wiki](https://github.com/sccreeper/goputer/wiki).


---

### Credits

- Raylib Go Bindings - [gen2brain/raylib-go](https://github.com/gen2brain/raylib-go)
  - Rendering library used by default gp32 frontend.
  - See also: Raylib - [raysan5/raylib](https://github.com/raysan5/raylib)
- Beep - [faiface/beep](https://github.com/faiface/beep)
  - Used for producing sound in the default gp32 frontend.
- CLI - [urfave/cli](https://github.com/urfave/cli)
  - Used by the compiler for getting input from the terminal.
- Color - [fatih/color](https://github.com/fatih/color)
  - Used for colouring & formatting the terminal output i.e. making it look nice.
- Termlink - [savioxavier/termlink](https://github.com/savioxavier/termlink)
  - Inserting links into the terminal.
  
---

### License

MIT - Refer to the `LICENSE` file