/**
 * 
 * @param {WebGL2RenderingContext} gl
 * @returns {{position: WebGLBuffer, textureCoord: WebGLBuffer}}
 */
function initBuffers(gl) {

    const positionBuffer = initPositionBuffer(gl)
    const textureCoordBuffer = initTextureBuffer(gl)

    return {
        position: positionBuffer,
        textureCoord: textureCoordBuffer,
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
    const positions = [
        1.0, 1.0, 
        -1.0, 1.0, 
        1.0, -1.0, 
        -1.0, -1.0
    ]

    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(positions), gl.STATIC_DRAW)

    return positionBuffer;
}

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 */
function initTextureBuffer(gl) {
    const textureCoordBuffer = gl.createBuffer()
    gl.bindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)

    const textureCoordinates = [
        0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0
    ]

    gl.bufferData(
        gl.ARRAY_BUFFER,
        new Float32Array(textureCoordinates),
        gl.STATIC_DRAW,
    )

    return textureCoordBuffer
}

export {initBuffers}