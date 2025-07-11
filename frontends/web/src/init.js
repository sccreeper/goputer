import { Compile, handleKeyDown, handleKeyUp, handleMouseMove, IOToggle, PeekRegister, Run } from "./app";
import globals from "./globals";
import { DownloadProgram, GetSharedCode, ShareCode, UploadBinary } from "./sharing";
import { ExamplesInit } from "./examples";
import { NewFileUI, SwitchFocus } from "./imports";
import { glInit } from "./gl/index";

//Cycles per second
export const CPS = 240;

//Init Go WASM
const go = new Go();
await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

document.getElementById("code-textarea").addEventListener("change", (e) => {
    updateFile(globals.focusedFile, document.getElementById("code-textarea").value);
})

export const new_file = document.getElementById("new-file");
new_file.addEventListener("click", NewFileUI)

document.getElementById("main-gpasm").addEventListener("click", SwitchFocus)

export const filesContainer = document.getElementById("code-names-container");

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
/**
 * @type {HTMLCanvasElement}
 */
const canvas = document.getElementById("render-canvas")
const gl = canvas.getContext("webgl2") ?? alert("Your browser does not support WebGL. goputer will not work.");
glInit(gl);

//Init audio

globals.audioContext = new (window.AudioContext || window.webkitAudioContext)();
globals.oscillator = globals.audioContext.createOscillator();
globals.audioVolume = globals.audioContext.createGain();

globals.audioVolume.gain.value = 0.0;
globals.oscillator.connect(globals.audioVolume);
globals.audioVolume.connect(globals.audioContext.destination);

//Init event listeners.
document.getElementById("compile-code-button").addEventListener("click", Compile)
document.getElementById("run-code-button").addEventListener("click", Run)
document.getElementById("download-code-button").addEventListener("click", DownloadProgram)
document.getElementById("share-code-button").addEventListener("click", ShareCode)
document.getElementById("upload-binary-button").addEventListener("click", UploadBinary)

document.getElementById("stop-code-button").addEventListener("click", function (e) {  
    clearInterval(globals.runInterval);

    // Clear sound

    globals.oscillator.frequency.value = 0;
    globals.audioVolume.gain.value = 0;

    if (globals.soundStarted) {
        globals.oscillator.stop();   
    }
    globals.soundStarted = false;

    // Clear IO lights

    for(let reg in globals.ioBulbs) {
        globals.ioBulbs[reg].setAttribute("on", "false")
    }

    // Clear other data

    globals.vmIsAlive = false;
    globals.videoText = "";

    // Clear canvas

    gl.clearColor(0.0, 0.0, 0.0, 1.0)
    gl.clear(gl.COLOR_BUFFER_BIT)

    canvas.setAttribute("running", "false");

    // Show editor again
    document.getElementById("code-editor").style.visibility = "visible"
    document.getElementById("binary-message").style.visibility = "hidden"
    document.getElementById("compile-code-button").disabled = false
    document.getElementById("share-code-button").disabled = false

})

//Init IO elements

for (let i = 0; i < document.getElementById("bulb-container").children.length; i++) {
    globals.ioBulbs[document.getElementById("bulb-container").children[i].getAttribute("reg")] = document.getElementById("bulb-container").children[i]
}

for (let i = 0; i < document.getElementById("switch-container").children.length; i++) {
    document.getElementById("switch-container").children[i].addEventListener("click", IOToggle)
}

canvas.addEventListener("mouseenter", () => {globals.mouseOverDisplay = true})
canvas.addEventListener("mouseleave", () => {globals.mouseOverDisplay = false})
canvas.addEventListener("mousemove", handleMouseMove)

document.addEventListener("keydown", handleKeyDown)
document.addEventListener("keyup", handleKeyUp)

globals.errorDiv = document.getElementById("error-div")

document.getElementById("error-clear-button").addEventListener("click", (e) => {

    let p = document.createElement("p")
    p.textContent = "No notifications."
    p.classList.add("text-center", "w-full");

    globals.errorCount = 0;

    globals.errorDiv.replaceChildren(p);

})

// Populate register keys

var peekRegDatalist = document.getElementById("peek-reg-datalist");

const peekRegHTML = document.getElementById("peek-reg-value");

peekRegHTML.addEventListener("click", e => {
    if (peekRegHTML.style.height != "auto") {
        peekRegHTML.style.height = "auto"
    } else {
        peekRegHTML.style.height = "2rem"
    }
})

Object.keys(registerInts).forEach(element => {
    
    let el = document.createElement("option")
    el.value = element;
    peekRegDatalist.appendChild(el);

});

const peekRegInput = document.getElementById("peek-reg");
peekRegInput.value = "";
peekRegInput.addEventListener("input", PeekRegister);
document.getElementById("peek-format-select").addEventListener("change", PeekRegister)

GetSharedCode();

export {canvas, gl as glContext, programCounterHTML, currentInstructionHTML, peekRegHTML, peekRegInput}
