export CGO_ENABLED="1"
export CC="x86_64-w64-mingw32-gcc"
export CXX="x86_64-w64-mingw32-g++"
export GOOS="windows"
export GOARCH="amd64"

mage build