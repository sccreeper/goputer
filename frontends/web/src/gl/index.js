import { mat4 } from "gl-matrix";
import { initShaderProgram } from "./shaders"
import { setPositionAttribute, setColourAttribute } from "./util";
import { initBuffers } from "./buffers";
import { createTexture } from "./textures";
import vertexSource from "./shaders/vs.vert"
import fragmentSource from "./shaders/fs.frag"
import globals from "../globals";

/**
 * @type {WebGLTexture}
 */
var drawTexture = null;

/**
 * @type {ProgramInfo}
 */
var programInfo = null;

/**
 * @type {{position: WebGLBuffer, textureCoord: WebGLBuffer}}
 */
var buffers = null;

/**
 * @typedef {Object} ProgramInfo
 * @property {WebGLProgram} program
 * @property {{vertexPosition: GLint, textureCoord: GLint}} attribLocations
 * @property {{uSampler: GLint}} uniformLocations
 */

// "Main" functions

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 */
function glInit(gl) {

    // Clear canvas
    gl.clearColor(0.0, 0.0, 0.0, 1.0)
    gl.clear(gl.COLOR_BUFFER_BIT)
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

    window.textureData.fill(0)
    drawTexture = createTexture(gl, window.textureData, 320, 240, gl.RGB)
    gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true)
    
    // Load shaders
    const shaderProgram = initShaderProgram(gl, vertexSource, fragmentSource);

    programInfo = {
        program: shaderProgram,
        attribLocations: {
            vertexPosition: gl.getAttribLocation(shaderProgram, "aVertexPosition"),
            textureCoord: gl.getAttribLocation(shaderProgram, "aTextureCoord"),
        },
        uniformLocations: {
            uSampler: gl.getUniformLocation(shaderProgram, "uSampler"),
        }
    }

    buffers = initBuffers(gl)

    drawScene(gl, programInfo, buffers, drawTexture)

}

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @param {ProgramInfo} programInfo 
 * @param {{position: WebGLBuffer}} buffers 
 * @param {WebGLTexture} texture
 */
function drawScene(gl, programInfo, buffers, texture) {
    gl.clearColor(0.0, 0.0, 0.0, 1.0)
    gl.clearDepth(1.0)
    gl.enable(gl.DEPTH_TEST)
    gl.depthFunc(gl.LEQUAL)

    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    setPositionAttribute(gl, buffers, programInfo)
    setTextureAttribute(gl, buffers, programInfo)

    gl.bindTexture(gl.TEXTURE_2D, texture)
    gl.texImage2D(
        gl.TEXTURE_2D,
        0,
        gl.RGB,
        320,
        240,
        0,
        gl.RGB,
        gl.UNSIGNED_BYTE,
        window.textureData
    )

    gl.activeTexture(gl.TEXTURE0)
    gl.bindTexture(gl.TEXTURE_2D, texture)

    gl.useProgram(programInfo.program)
    gl.uniform1i(programInfo.uniformLocations.uSampler, 0)

    {
        const offset = 0;
        const vertexCount = 4;
        gl.drawArrays(gl.TRIANGLE_STRIP, offset, vertexCount)
    }
}

function drawSceneSimple(gl) {
    drawScene(gl, programInfo, buffers, drawTexture)
}

/**
 * 
 * @param {WebGL2RenderingContext} gl 
 * @param {{position: WebGLBuffer, textureCoord: WebGLBuffer}} buffers 
 * @param {ProgramInfo} programInfo 
 */
function setTextureAttribute(gl, buffers, programInfo) {
    const num = 2; 
    const type = gl.FLOAT; 
    const normalize = false;
    const stride = 0; 
    const offset = 0;
    
    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.textureCoord);
    
    gl.vertexAttribPointer(
        programInfo.attribLocations.textureCoord,
        num,
        type,
        normalize,
        stride,
        offset,
    );
    
    gl.enableVertexAttribArray(programInfo.attribLocations.textureCoord);
}

export {glInit, drawSceneSimple}