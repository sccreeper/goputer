attribute vec4 aVertexPosition;
attribute vec4 aVertexColour;

uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;

varying lowp vec4 vColour;

void main() {
    gl_Position = uProjectionMatrix * uModelViewMatrix * aVertexPosition;
    vColour = aVertexColour;
}
