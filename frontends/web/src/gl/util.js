/**
 * @import { ProgramInfo } from "./index.js"
 */

/**
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

export {setColourAttribute, setPositionAttribute}