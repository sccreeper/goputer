import vertexSource from "./shaders/vs.vert";
import fragmentSource from "./shaders/fs.frag";
import { mat4 } from "gl-matrix";

// "Main" functions

/**
 * @typedef {Object} ProgramInfo
 * @property {WebGLProgram} program
 * @property {{vertexPosition: GLint, vertexColour: GLint}} attribLocations
 * @property {{projectionMatrix: WebGLUniformLocation, modelViewMatrix: WebGLUniformLocation}} uniformLocations
 */

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 */
export function glInit(gl) {

    // Clear canvas
    gl.clearColor(0.0, 0.0, 0.0, 1.0)
    gl.clear(gl.COLOR_BUFFER_BIT)

    console.log(vertexSource)
    
    // Load shaders
    const shaderProgram = initShaderProgram(gl, vertexSource, fragmentSource);

    const programInfo = {
        program: shaderProgram,
        attribLocations: {
            vertexPosition: gl.getAttribLocation(shaderProgram, "aVertexPosition"),
            vertexColour: gl.getAttribLocation(shaderProgram, "aVertexColour")
        },
        uniformLocations: {
          projectionMatrix: gl.getUniformLocation(shaderProgram, "uProjectionMatrix"),
          modelViewMatrix: gl.getUniformLocation(shaderProgram, "uModelViewMatrix"),
        },
    }

    const buffers = initBuffers(gl)

    drawScene(gl, programInfo, buffers)

}

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @param {ProgramInfo} programInfo 
 * @param {{position: WebGLBuffer}} buffers 
 */
function drawScene(gl, programInfo, buffers) {
    gl.clearColor(0.0, 0.0, 0.0, 1.0)
    gl.clearDepth(1.0)
    gl.enable(gl.DEPTH_TEST)
    gl.depthFunc(gl.LEQUAL)

    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    const fieldOfView = (45 * Math.PI) / 180
    const aspect = gl.canvas.clientWidth / gl.canvas.clientHeight
    const zNear = 0.1
    const zFar = 100.0
    const projectionMatrix = mat4.create()

    mat4.perspective(projectionMatrix, fieldOfView, aspect, zNear, zFar)

    const modelViewMatrix = mat4.create()

    mat4.translate(
        modelViewMatrix,
        modelViewMatrix,
        [-0.0, 0.0, -6.0]
    )

    setPositionAttribute(gl, buffers, programInfo)
    setColourAttribute(gl, buffers, programInfo)

    gl.useProgram(programInfo.program)

    gl.uniformMatrix4fv(
        programInfo.uniformLocations.projectionMatrix,
        false,
        projectionMatrix
    )

    gl.uniformMatrix4fv(
        programInfo.uniformLocations.modelViewMatrix,
        false,
        modelViewMatrix
    )

    {
        const offset = 0;
        const vertexCount = 4;
        gl.drawArrays(gl.TRIANGLE_STRIP, offset, vertexCount)
    }
}

// "Accessory" functions

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

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @param {{position: WebGLBuffer, colour: WebGLBuffer}} buffers 
 * @param {ProgramInfo} programInfo 
 */
function setPositionAttribute(gl, buffers, programInfo) {
    const numComponents = 2
    const type = gl.FLOAT
    const normalise = false
    const stride = 0
    const offset = 0

    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.position)
    gl.vertexAttribPointer(
        programInfo.attribLocations.vertexPosition,
        numComponents,
        type,
        normalise,
        stride,
        offset,
    )
    gl.enableVertexAttribArray(programInfo.attribLocations.vertexPosition)
}

/**
 * 
 * @param {WebGLRenderingContext} gl 
 * @param {{position: WebGLBuffer, colour: WebGLBuffer}} buffers 
 * @param {ProgramInfo} programInfo 
 */
function setColourAttribute(gl, buffers, programInfo) {
    
    const numComponents = 4;
    const type = gl.FLOAT;
    const normalise = false;
    const stride = 0;
    const offset = 0;

    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.colour)

    gl.vertexAttribPointer(
        programInfo.attribLocations.vertexColour,
        numComponents,
        type,
        normalise,
        stride,
        offset
    )

    gl.enableVertexAttribArray(programInfo.attribLocations.vertexColour)

}