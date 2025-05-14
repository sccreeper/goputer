/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @param {String} vsSource 
 * @param {String} fsSource 
 */
function initShaderProgram(gl, vsSource, fsSource) {
    
    // Load shaders
    const vertexShader = loadShader(gl, gl.VERTEX_SHADER, vsSource)
    const fragmentShader = loadShader(gl, gl.FRAGMENT_SHADER, fsSource)

    // Create shader program

    const shaderProgram = gl.createProgram()
    gl.attachShader(shaderProgram, vertexShader)
    gl.attachShader(shaderProgram, fragmentShader)
    gl.linkProgram(shaderProgram)

    if (!gl.getProgramParameter(shaderProgram, gl.LINK_STATUS)) {
        alert(`Unable to init shader`)

        return null
    } else {
        return shaderProgram
    }
}

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @param {GLenum} type 
 * @param {String} source 
 */
function loadShader(gl, type, source) {
    
    const shader = gl.createShader(type)

    gl.shaderSource(shader, source);
    gl.compileShader(shader)

    if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
        alert(`Unable to compile shader program ${gl.getShaderInfoLog(shader)}`)
        gl.deleteShader(shader)
        return null
    } else {
        return shader
    }
}

export {initShaderProgram}