import { Compile, handleKeyDown, handleKeyUp, handleMouseMove, IOToggle, PeekRegister, Run } from "./app";
import { clearCanvas } from "./canvas_util";
import globals from "./globals";
import { DownloadProgram, GetSharedCode, ShareCode } from "./sharing";
import { ExamplesInit } from "./examples";
import { NewFileUI, SwitchFocus } from "./imports";

//Cycles per second
export const CPS = 240;

//Init Go WASM
const go = new Go();
await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

document.getElementById("code-textarea").addEventListener("change", (e) => {
    updateFile(globals.focused_file, document.getElementById("code-textarea").value);
})

export const new_file = document.getElementById("new-file");
new_file.addEventListener("click", NewFileUI)

document.getElementById("main-gpasm").addEventListener("click", SwitchFocus)

export const files_container = document.getElementById("code-names-container");

//Set the version

fetch("/ver")
.then((response) => response.text())
.then((data) => {

    let hash = data.split(/\r?\n/)[0];
    let time = data.split(/\r?\n/)[1];

    document.getElementById("version").textContent = `${hash.substring(0, 10)}`;
    document.getElementById("version").setAttribute("href", `https://github.com/sccreeper/goputer/commit/${hash}`);
    document.getElementById("build-date").textContent = time;

})

export const examplesDiv = document.getElementById("examples-div");

ExamplesInit();

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
document.getElementById("compile-code-button").addEventListener("click", Compile)
document.getElementById("run-code-button").addEventListener("click", Run)
document.getElementById("download-code-button").addEventListener("click", DownloadProgram)
document.getElementById("share-code-button").addEventListener("click", ShareCode)

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
    canvas.setAttribute("running", "false");
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

globals.error_div = document.getElementById("error-div")

document.getElementById("error-clear-button").addEventListener("click", (e) => {

    let p = document.createElement("p")
    p.textContent = "No notifications."
    p.classList.add("text-center", "w-full");

    globals.error_count = 0;

    globals.error_div.replaceChildren(p);

})

// Populate register keys

var peek_reg_datalist = document.getElementById("peek-reg-datalist");

const peekRegHTML = document.getElementById("peek-reg-value");

Object.keys(registerInts).forEach(element => {
    
    let el = document.createElement("option")
    el.value = element;
    peek_reg_datalist.appendChild(el);

});

const peekRegInput = document.getElementById("peek-reg");
peekRegInput.value = "";
peekRegInput.addEventListener("input", PeekRegister);

GetSharedCode();

export {canvas, renderContext, programCounterHTML, currentInstructionHTML, peekRegHTML, peekRegInput}
