#include <raylib.h>
#include <unistd.h>

Color convertColour(__uint32_t col) {

    Color result = {
        .r = (col >> 24) & 0xFF,
        .g = (col >> 16) & 0xFF,
        .b = (col >> 8) &0xFF,
        .a = (col) & 0xFF,
    };

    return result;

}