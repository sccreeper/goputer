![goputer logo](./.github/logo-readme.png)

<sup>`Go + Computer = goputer`</sup>

---

A computer emulator/virtual machine that intends to demonstrate how basic computers work at a low level.

---

**Contents**
- [Features](#features)
  - [Complete](#complete)
  - [Working on](#working-on)
  - [In the future](#in-the-future)
- [Documentation \& getting started.](#documentation--getting-started)
- [Project layout](#project-layout)
- [Build instructions](#build-instructions)
  - [Docker](#docker)
  - [Linux](#linux)
- [Testing](#testing)
- [Credits](#credits)
  - [GP32 Frontend](#gp32-frontend)
  - [goputerpy Frontend](#goputerpy-frontend)
  - [Web playground/frontend.](#web-playgroundfrontend)
  - [CLI tool](#cli-tool)
  - [GUI launcher](#gui-launcher)
  - [Other](#other)
- [License](#license)

---

### Features

#### Complete

- Custom assembly language and compiler.
- Custom runtime.
- Standalone executables.
- Frontends to show VM output.
- A [WASM based runtime](https://goputer.oscarcp.net) that runs in a web browser.
- Expansion cards/modules.

#### Working on

- IDE for easy development.

#### In the future

- Rewrite of compiler.
- High level language.
- Non-native plugins via Lua.

---

### Documentation & getting started.

See the [project wiki](https://github.com/sccreeper/goputer/wiki) or try the playground at [goputer.oscarcp.net](https://goputer.oscarcp.net).

---

### Project layout

- `frontends` Contains source for the frontends.
 - `frontends/web` The WASM frontend.
 - `frontends/gp32` The Go frontend.
 - `frontends/goputerpy` The Python frontend.
- `examples` A list of example code to get started with.
- `cmd/goputer` The CLI tool for compiling, running & disassembling code.
- `cmd/launcher` The GUI for running code.
- `pkg` Shared code. Includes the compiler, VM runtime and constants for instructions and registers.

---

### Build instructions

Build instructions for Linux, other platforms are not supported at the moment as native plugins do not work at all. See [plugin](https://pkg.go.dev/plugin)

#### Docker

If you have [Docker](https://www.docker.com/) installed, you can build goputer in Docker without installing additional dependencies by running:

```
./docker_entrypoint.sh
```

This will build the container and then run `mage dev` inside the container, outputting to the build directory.

**Note:** This uses a Fedora container, so the output will only work on Linux based systems.

#### Linux

**Perquisites**

- Languages
  - Python ^3.10
  - Go ^1.19
  - NodeJS ^18.X

- Build tools
  - [Poetry](https://python-poetry.org/)
  - [Mage](https://magefile.org/)

- Library requirements
  
  Requirements for different display servers & audio.
  
  ##### Ubuntu <!-- omit in toc -->
  
  ###### x11 <!-- omit in toc -->
  
  ```
  libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev
  ```

  ###### Wayland <!-- omit in toc -->

  ```
  libgl1-mesa-dev libwayland-dev libxkbcommon-dev 
  ```

  ###### Audio <!-- omit in toc -->

  ```
  libasound2-dev
  ```

  ##### Fedora <!-- omit in toc -->

  ###### x11 <!-- omit in toc -->

  ```
  mesa-libGL-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel
  ```

  ###### Wayland <!-- omit in toc -->

  ```
  mesa-libGL-devel wayland-devel libxkbcommon-devel
  ```

  ###### Audio <!-- omit in toc -->
  ```
  alsa-lib-devel
  ```

**Building**

1. Install the prerequisites that are mentioned above.
2. Check that everything works.
   ```
   $ node --version
   v18.12.1
   $ mage --version
   Mage Build Tool v1.14.0
   Build Date: 2022-11-09T16:45:13Z
   Commit: 300bbc8
   built with: go1.19.3
   $ poetry --version
   Poetry version 1.1.13
   ```
3. Clone the repository from GitHub

    ```
    git clone https://github.com/sccreeper/goputer
    cd goputer
    ```

4. Build the project. (This step shouldn't take that long depending on your hardware)
   
   ```
   mage dev
   ```
    Alternatively you can run `mage all` to *not* build the examples.

5. Go to the build directory and run the `hello_world` example.
   ```
   cd build
   ./goputer run -f gp32 -e ./examples/hello_world
   ```

---

### Testing

There a small suite of tests written for testing the core of goputer.

To run them use:

```
go test ./tests -v
```

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

#### goputerpy Frontend

- pygame - [pygame/pygame](https://github.com/pygame/pygame)
  - Used for rendering & sound output.
- numpy - [numpy/numpy](https://github.com/numpy/numpy)
  - Used partially for sound generation
- poetry - [python-poetry/poetry](https://github.com/python-poetry/poetry)
  - Very useful for Python dependency management.

#### Web playground/frontend.
- Tailwind CSS - [tailwindlabs/tailwindcss](https://github.com/tailwindlabs/tailwindcss)
  - CSS library used to build the UI.
- Parcel - [parcel-bundler/parcel](https://github.com/parcel-bundler/parcel)
  - Module bundler.
- Parcel Static Files Copy - [elwin013/parcel-reporter-static-files-copy](https://github.com/elwin013/parcel-reporter-static-files-copy)
  - Used for copying static files into the final dist directory.
- Bootstrap Icons - [twbs/icons](https://github.com/twbs/icons)
  - Icons used heavily by the UI.
- Microtip - [ghosh/microtip](https://github.com/ghosh/microtip)
  - Tooltip library used in UI.

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

#### Other
- Mage - [magefile/mage](https://github.com/magefile/mage)
  - Build system used for development.
  
---

### License

MIT - Refer to the `LICENSE` file
