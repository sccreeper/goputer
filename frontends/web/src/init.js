import { Compile, Run } from "./app";
import { clearCanvas } from "./canvas_util";
import globals from "./globals";

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

//Init audio

globals.audioContext = new (window.AudioContext || window.webkitAudioContext)();
globals.oscillator = globals.audioContext.createOscillator();
globals.audioVolume = globals.audioContext.createGain();

globals.audioVolume.gain.value = 0.0;
globals.oscillator.connect(globals.audioVolume);
globals.audioVolume.connect(globals.audioContext.destination);

//Init event listeners.
document.getElementById("compile-code-button").addEventListener("click", function (e) {  
    Compile();
})

document.getElementById("run-code-button").addEventListener("click", function (e) {  
    Run();
})

document.getElementById("stop-code-button").addEventListener("click", function (e) {  
    clearInterval(globals.runInterval);

    globals.oscillator.frequency.value = 0;
    globals.audioVolume.gain.value = 0;

    globals.oscillator.stop();

    globals.sound_started = false;
})


export {canvas, renderContext}
