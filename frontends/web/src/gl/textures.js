/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @param {Uint8Array} source
 * @param {number} width 
 * @param {number} height 
 * @param {number} format 
 */
function createTexture(gl, source, width, height, format) {

    const texture = gl.createTexture()
    gl.bindTexture(gl.TEXTURE_2D, texture)
    
    const level = 0;
    const border = 0;

    gl.texImage2D(
        gl.TEXTURE_2D,
        level,
        format,
        width,
        height,
        border,
        format,
        gl.UNSIGNED_BYTE,
        source
    )

    if (isPowerOf2(width) && isPowerOf2(height)) {
        gl.generateMipmap(gl.TEXTURE_2D)
    } else {
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    }

    return texture

}

function isPowerOf2(value) {
    return (value & (value - 1)) === 0;
}

export {createTexture}