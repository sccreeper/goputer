rm -r ./build

mkdir build
go build -buildmode=plugin -o "goputer.sys.so"
mv "goputer.sys.so" "./build/goputer.sys.so"