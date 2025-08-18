![goputer logo](./.github/logo-readme.png)

<sup>`Go + Computer = goputer`</sup>

---

A computer emulator/virtual machine that intends to demonstrate how basic computers work at a low level.

---

**Contents**
- [Features](#features)
  - [Complete](#complete)
  - [Working on](#working-on)
- [Documentation \& getting started.](#documentation--getting-started)
- [Project layout](#project-layout)
- [Build instructions](#build-instructions)
  - [Docker](#docker)
  - [Linux](#linux)
      - [Other](#other)
  - [Windows](#windows)
  - [Steps](#steps)
- [Development](#development)
  - [Testing](#testing)
- [Credits](#credits)
  - [Core](#core)
  - [GP32 Frontend](#gp32-frontend)
  - [goputerpy Frontend](#goputerpy-frontend)
  - [Web playground/frontend.](#web-playgroundfrontend)
  - [CLI tool](#cli-tool)
  - [GUI launcher](#gui-launcher)
  - [Other](#other-1)
- [License](#license)

---

### Features

#### Complete

- Custom assembly language and compiler.
- Custom runtime.
- Standalone executables.
- Frontends to show VM output.
- A [WASM based runtime](https://goputer.oscarcp.net) that runs in a web browser.
- Expansion modules written in Lua.

#### Working on

- IDE for easy development.
- High level language.

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
- `expansions` Source code for all of the default expansions.

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

Partial cross compilation targeting Windows systems can be done on Linux, see [Windows](#windows).

**Perquisites**

- Languages
  - Python ^3.10
  - Go ^1.23
  - NodeJS ^18.X

- Build tools
  - [Poetry](https://python-poetry.org/)
  - [Mage](https://magefile.org/)

For Node I would recommend installing [NVM](https://github.com/nvm-sh/nvm) (Node Version Manager).

##### Fedora <!-- omit in toc -->

Tested on Fedora 42.

###### x11 <!-- omit in toc -->

```
mesa-libGL-devel libXi-devel libXcursor-devel libXrandr-devel libXinerama-devel libXxf86vm-devel 
```

###### Wayland <!-- omit in toc -->

```
mesa-libGL-devel wayland-devel libxkbcommon-devel
```

###### Other

```
gtk3-devel golang golang-tests python3 
```

###### Audio <!-- omit in toc -->
```
alsa-lib-devel
```

**Building (For Linux)**

1. Install the prerequisites that are mentioned above.
2. Check that everything works.
   ```
   $ node --version
    v22.13.1
   $ mage --version
    Mage Build Tool <not set>
    Build Date: <not set>
    Commit: <not set>
    built with: go1.24.5
   $ poetry --version
    Poetry version 1.1.13
   $ go version 
    go version go1.24.5 linux/amd64
   $ python3 --version
    Python 3.13.2
   ```
3. Clone the repository from GitHub

    ```
    
    git clone https://github.com/sccreeper/goputer
    cd goputer
    
    ```

4. Activate the virtual environment for Python (unneeded if not using or building the Python frontend):
  
  ```
  eval $(poetry env activate)
  ```

5. Build the project. (This step shouldn't take that long depending on your hardware)
   
   ```
   mage build
   ```

    Alternatively you can run `mage dev frontend.gp32,expansion.goputer.sys` to only build the `gp32` frontend.


6. Go to the build directory and run the `hello_world.gp` example.
   ```
   cd build
   ./goputer run -f gp32 -e ./examples/hello_world.gp
   ```

#### Windows

**TLDR;**
- Windows executables can be built using cross compilation on Linux (run `build_win.sh`). 
- [mingw](https://www.mingw-w64.org/) must be installed, in addition to the prerequisites listed under [Linux](#linux).
- The gp32 frontend will cross compile however the goputerpy one will not, this is because [pyinstaller](https://pyinstaller.org/en/stable/) does not support cross compilation.
- [MSYS2](https://www.msys2.org/) is the recommended compilation environment if you are compiling on Windows and targeting Windows.

#### Steps

1. Install MSYS2.
2. Once in MSYS2 (mingw64) install the required packages:
    ```
    pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-go mingw-w64-x86_64-nodejs mingw-w64-x86_64-python mingw-w64-x86_64-python-poetry mingw-w64-python-certifi mingw-w64-python-pip mingw-w64-x86_64-pyinstaller git
    ```
3. Install mage:
   ```
   export GOBIN=$HOME/go/bin
   export PATH=$GOBIN:$PATH
   go install github.com/magefile/mage@latest
   ```
   *Note:* You may need to set those environment variables again.
4. Clone the repository:
    ```
    https://github.com/sccreeper/goputer.git
    ```
5. cd into the repository and install dependencies:
   ```
   poetry install
   npm install
   ```
6. Build:
    ```
    cd goputer
    sh ./build_win.sh
    ```
    Optionally run `eval $(poetry env activate)` if building the Python frontend as well.

---

### Development

In addition to the dependencies in [Building](#build-instructions), you should also have [`golangci-lint`](https://golangci-lint.run) (V2) installed.

#### Testing

There a small suite of tests written for testing the core of goputer.

To run them use:

```
go test ./tests -v
```

---

### Credits

#### Core

- Lua runtime - [Shopify/go-lua](https://github.com/Shopify/go-lua)
  - The Lua runtime used by expansions.

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
- Mediabunny - [Vanilagy/mediabunny](https://github.com/Vanilagy/mediabunny)
  - Used for recording and muxing MP4 files
- Zip.js - [gildas-lormeau/zip.js](https://github.com/gildas-lormeau/zip.js/)
  - Used to create and load ZIP archives
- Dexie.js - [dexie/Dexie.js](https://github.com/dexie/Dexie.js)
  - Used for interacting with IndexedDB.

#### CLI tool

- CLI - [urfave/cli](https://github.com/urfave/cli)
  - Used by the compiler for getting input from the terminal.
- Color - [fatih/color](https://github.com/fatih/color)
  - Used for colouring & formatting the terminal output i.e. making it look nice.
- Termlink - [savioxavier/termlink](https://github.com/savioxavier/termlink)
  - Inserting links into the terminal.
- Tview - [rivo/tview](https://github.com/rivo/tview)
  - Used for the profiler UI.

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
