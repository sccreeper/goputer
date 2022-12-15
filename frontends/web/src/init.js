import { Compile, handleKeyDown, handleKeyUp, handleMouseMove, IOToggle, Run } from "./app";
import { clearCanvas } from "./canvas_util";
import globals from "./globals";

//Cycles per second
export const CPS = 240;

//Init Go WASM
const go = new Go();
await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

//Debug UI
const programCounterHTML = document.getElementById("program-counter");
const currentInstructionHTML = document.getElementById("current-instruction");

//Init Canvas
const canvas = document.getElementById("render-canvas")
const renderContext = canvas.getContext('2d');

clearCanvas(renderContext, "black");

//Init audio

globals.audio_context = new (window.AudioContext || window.webkitAudioContext)();
globals.oscillator = globals.audio_context.createOscillator();
globals.audio_volume = globals.audio_context.createGain();

globals.audio_volume.gain.value = 0.0;
globals.oscillator.connect(globals.audio_volume);
globals.audio_volume.connect(globals.audio_context.destination);

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
    globals.audio_volume.gain.value = 0;

    if (globals.sound_started) {
        globals.oscillator.stop();   
    }

    globals.sound_started = false;
    globals.vmIsAlive = false;
    globals.video_text = "";

    clearCanvas(renderContext, "black");
})

//Init IO elements

for (let i = 0; i < document.getElementById("bulb-container").children.length; i++) {
    globals.io_bulbs[document.getElementById("bulb-container").children[i].getAttribute("reg")] = document.getElementById("bulb-container").children[i]
}

for (let i = 0; i < document.getElementById("switch-container").children.length; i++) {
    document.getElementById("switch-container").children[i].addEventListener("click", IOToggle)
}

canvas.addEventListener("mouseenter", () => {globals.mouse_over_display = true})
canvas.addEventListener("mouseleave", () => {globals.mouse_over_display = false})
canvas.addEventListener("mousemove", handleMouseMove)

document.addEventListener("keydown", handleKeyDown)
document.addEventListener("keyup", handleKeyUp)

export {canvas, renderContext, programCounterHTML, currentInstructionHTML}
