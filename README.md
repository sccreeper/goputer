# goputer <!-- omit in toc -->
<sup>`Go + Computer = goputer`</sup>

---

A computer emulator/virtual machine intended to demonstrate how computers work at a low level.

---

**Contents**
- [Features](#features)
  - [Complete](#complete)
  - [Working on](#working-on)
  - [In the future](#in-the-future)
- [Documentation \& getting started.](#documentation--getting-started)
- [Project layout](#project-layout)
- [Credits](#credits)
  - [GP32 Frontend](#gp32-frontend)
  - [CLI tool](#cli-tool)
  - [GUI launcher](#gui-launcher)
    - [Other](#other)
- [License](#license)

---

### Features

**Note:** Most of these have not been fully implemented yet.

#### Complete

- Custom assembly language and compiler.
- Custom runtime.
- Standalone executables.

#### Working on

- Frontend that shows output from VM backend.

#### In the future

- WASM based frontend.
- High level language.
- IDE for easy development.

---

### Documentation & getting started.

See the [project wiki](https://github.com/sccreeper/goputer/wiki).

---

### Project layout

- `frontends` Contains source for the frontends.
- `examples` A list of example code to get started with.
- `cmd/goputer` The CLI tool for compiling, running & disassembling code.
- `pkg` Shared code. Includes the compiler & the VM runtime.

---

### Credits

#### GP32 Frontend

- Raylib Go Bindings - [gen2brain/raylib-go](https://github.com/gen2brain/raylib-go)
  - Very useful rendering library.
  - See also: Raylib - [raysan5/raylib](https://github.com/raysan5/raylib)
- Beep - [faiface/beep](https://github.com/faiface/beep)
  - Used for producing sound on the fly.
- TOML - [BurntSushi/toml](https://github.com/BurntSushi/toml)
  - Configuration format used by frontends.
  - See also: TOML - [toml-lang/toml](https://github.com/toml-lang/toml)

#### CLI tool

- CLI - [urfave/cli](https://github.com/urfave/cli)
  - Used by the compiler for getting input from the terminal.
- Color - [fatih/color](https://github.com/fatih/color)
  - Used for colouring & formatting the terminal output i.e. making it look nice.
- Termlink - [savioxavier/termlink](https://github.com/savioxavier/termlink)
  - Inserting links into the terminal.

#### GUI launcher
- Fyne - [fyne-io/fyne](https://github.com/fyne-io/fyne)
  - UI library used
- dialog - [sqweek/dialog](https://github.com/sqweek/dialog)
  - Used for file open dialogs.

##### Other
- Mage - [magefile/mage](https://github.com/magefile/mage)
  - Build system used for development.
  
---

### License

MIT - Refer to the `LICENSE` file