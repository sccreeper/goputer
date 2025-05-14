import { mat4 } from "gl-matrix";
import { initShaderProgram } from "./shaders"
import { setPositionAttribute, setColourAttribute } from "./util";
import { initBuffers } from "./buffers";
import vertexSource from "./shaders/vs.vert"
import fragmentSource from "./shaders/fs.frag"

/**
 * @typedef {Object} ProgramInfo
 * @property {WebGLProgram} program
 * @property {{vertexPosition: GLint, vertexColour: GLint}} attribLocations
 * @property {{projectionMatrix: WebGLUniformLocation, modelViewMatrix: WebGLUniformLocation}} uniformLocations
 */

// "Main" functions

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 */
export function glInit(gl) {

    // Clear canvas
    gl.clearColor(0.0, 0.0, 0.0, 1.0)
    gl.clear(gl.COLOR_BUFFER_BIT)
    
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