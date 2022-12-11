import { Compile, Run } from "./app";
import { clearCanvas } from "./canvas_util";

export const FPS = 60;

//Init Go WASM
const go = new Go();
await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

//Init Canvas
const canvas = document.getElementById("render-canvas")
const renderContext = canvas.getContext('2d');

clearCanvas(renderContext, "black");

//Init event listeners.
document.getElementById("compile-code-button").addEventListener("click", function (e) {  
    Compile();
})

document.getElementById("run-code-button").addEventListener("click", function (e) {  
    Run();
})


export {canvas, renderContext}
