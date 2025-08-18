$env:CGO_ENABLED = "1"
$env:CC = "x86_64-w64-mingw32-gcc"
$env:CXX = "x86_64-w64-mingw32-g++"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

mage build
