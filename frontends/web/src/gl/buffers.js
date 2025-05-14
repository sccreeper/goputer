/**
 * 
 * @param {WebGL2RenderingContext} gl
 * @returns {{position: WebGLBuffer, colour: WebGLBuffer}}
 */
function initBuffers(gl) {

    const positionBuffer = initPositionBuffer(gl)
    const colourBuffer = initColourBuffer(gl)

    return {
        position: positionBuffer,
        colour: colourBuffer,
    }

}

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @returns {WebGLBuffer}
 */
function initPositionBuffer(gl) {
    const positionBuffer = gl.createBuffer()
    gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer)
    const positions = [1.0, 1.0, -1.0, 1.0, 1.0, -1.0, -1.0, -1.0]

    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(positions), gl.STATIC_DRAW)

    return positionBuffer;
}

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 */
function initColourBuffer(gl) {
    const colours = [
        1.0, // White
        1.0,
        1.0,
        1.0,
        1.0, // Red
        0.0,
        0.0,
        1.0,
        0.0, // Green
        1.0,
        0.0,
        1.0,
        0.0, // Blue
        0.0,
        1.0,
        0.0,
    ]

    const colourBuffer = gl.createBuffer()
    gl.bindBuffer(gl.ARRAY_BUFFER, colourBuffer)
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(colours), gl.STATIC_DRAW)

    return colourBuffer
}

export {initBuffers}