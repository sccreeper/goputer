cmake_minimum_required(VERSION 3.5)
project(GPC)

# Set output directories
set(CMAKE_BINARY_DIR "build/")

set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)

add_subdirectory(lib/raylib)

file(GLOB SOURCES "src/*.c")

add_custom_target(
    bindings
    COMMAND go build -buildmode=c-shared -o ${CMAKE_BINARY_DIR}/bin/bindings.so ${CMAKE_SOURCE_DIR}/go/bindings.go
)

add_executable(GPC ${SOURCES})

target_include_directories(
    GPC
    PRIVATE 
        "lib/raylib/src"
        "src"
)

add_dependencies(GPC bindings)
target_link_libraries(GPC PRIVATE raylib)