import { Compile, handleKeyDown, handleKeyUp, handleMouseMove, PeekRegister, Run, SaveVideo } from "./app";
import globals from "./globals";
import { DownloadProgram, DownloadAll, UploadBinary } from "./sharing";
import { ExamplesInit } from "./examples";
import { glInit } from "./gl/index";
import { ToggleRecording } from "./recording";
import { codeArea, InitImage, NewFile, NewFileUI, SwitchFocus, tabsContainer } from "./imports";
import { db, fileTableName } from "./db";
import { goputer } from "./goputer";
import "./ui/io.js";
import { clamp } from "./util.js";

//Init Go WASM before anything else
const go = new Go();
await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

// Editor init

let files = await db.table(fileTableName).toArray()

if (files.length != 0) {
    files.forEach(element => {

        NewFile(element.name, false)

        if (element.type == "image") {

            
            const imgBlob = new Blob([element.data])
            InitImage(imgBlob, element.name)

        }
        
        goputer.files.update(
            element.name,
            element.data,
            element.data.length,
            element.type,
            true,
            false,
        )

    });

    // Make sure main.gpasm is at beginning

    const mainTab = document.querySelector(`code-tab[filename="main.gpasm"]`)
    document.getElementById("code-names-container").prepend(mainTab)

    SwitchFocus("main.gpasm");   
} else {
    NewFile("main.gpasm");
}

//Cycles per second
export const CPS = 240;

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

// Examples init

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
const gl = canvas.getContext("webgl2", {preserveDrawingBuffer: true}) ?? alert("Your browser does not support WebGL. goputer will not work.");
glInit(gl);

//Init audio

globals.audioContext = new (window.AudioContext || window.webkitAudioContext)();
globals.oscillator = globals.audioContext.createOscillator();
globals.audioVolume = globals.audioContext.createGain();
globals.audioMediaStreamDestination = globals.audioContext.createMediaStreamDestination();

globals.audioVolume.gain.value = 0.0;
globals.oscillator.connect(globals.audioVolume);
globals.audioVolume.connect(globals.audioContext.destination);
globals.audioVolume.connect(globals.audioMediaStreamDestination);

//Init event listeners for general button UI.

document.getElementById("compile-code-button").addEventListener("click", Compile)
document.getElementById("run-code-button").addEventListener("click", Run)
document.getElementById("download-code-button").addEventListener("click", DownloadProgram)
document.getElementById("download-all-button").addEventListener("click", DownloadAll)
document.getElementById("upload-binary-button").addEventListener("click", UploadBinary)
document.getElementById("save-video-button").addEventListener("click", SaveVideo)
document.getElementById("record-video-button").addEventListener("click", ToggleRecording)

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
        globals.ioBulbs[reg].setAttribute("enabled", "false")
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
    document.getElementById("download-all-button").disabled = false

})

document.getElementById("toggle-contrast-button").addEventListener("click", () => {
    document.documentElement.classList.toggle("high-contrast");
})

//Init IO elements

for (let i = 0; i < document.getElementById("bulb-container").children.length; i++) {
    globals.ioBulbs[document.getElementById("bulb-container").children[i].getAttribute("reg")] = document.getElementById("bulb-container").children[i]
}

canvas.addEventListener("mouseenter", () => {globals.mouseOverDisplay = true})
canvas.addEventListener("mouseleave", () => {globals.mouseOverDisplay = false})
canvas.addEventListener("mousemove", handleMouseMove)

document.addEventListener("keydown", handleKeyDown)
document.addEventListener("keyup", handleKeyUp)

// Error messages

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

// Global keyboard shortcuts

document.addEventListener("keydown", (e) => {
    if (e.altKey && e.key.charCodeAt(0) >= 49 && e.key.charCodeAt(0) <= 56) { // Alt+1-8 - use switch
        
        e.preventDefault();
        document.getElementById("switch-container").children[e.key.charCodeAt(0)-49].click();
    
    } else if (e.ctrlKey && e.altKey && (e.key == "ArrowLeft" || e.key == "ArrowRight")) { // Ctrl+Alt+<-/-> - change tab
        
        e.preventDefault();

        const increment = e.key == "ArrowLeft" ? -1 : 1

        const tabs = tabsContainer.querySelectorAll("code-tab");
        let selectedIndex = 0;

        for (let i = 0; i < tabs.length; i++) {
            
            if (tabs[i].selected) {
                selectedIndex = i;
                break;
            }
            
        }

        selectedIndex += increment;
        selectedIndex = clamp(selectedIndex, 0, tabsContainer.length - 1);

        SwitchFocus(tabs[selectedIndex].filename);
        
    } else if (e.ctrlKey && e.altKey && e.code == "KeyN") { // Ctrl+Alt+N - new file
        e.preventDefault();
        NewFileUI(e)
    } else if (e.ctrlKey && e.altKey && e.code == "KeyW") { // Ctrl+Alt+W - delete file
        e.preventDefault();
        tabsContainer.querySelector("code-tab[selected=true]").deleteFile(e);
    } else if (e.ctrlKey && e.altKey && e.code == "KeyR") { // Ctrl+Alt+R - rename file
        e.preventDefault();
        tabsContainer.querySelector("code-tab[selected=true]").renameFile(e);
    } else if (e.altKey && e.code == "F1") { // Compile
        e.preventDefault()
        Compile()
    } else if (e.altKey && e.code == "F2") { // Run
        e.preventDefault()
        Run()
    } else if (e.altKey && e.code == "F3") { // Reset
        e.preventDefault()
        document.getElementById("stop-code-button").click()
    } else if (e.altKey && e.code == "F5") { // Clear errors
        e.preventDefault()
        document.getElementById("error-clear-button").click()
    } else if (e.ctrlKey && e.altKey && e.code == "KeyF") { // Ctrl+Alt+F - focus text file in editor
        if (goputer.files.type(globals.focusedFile) == "text") {
            e.preventDefault()
            codeArea.focus()
        }
    } else if (e.ctrlKey && e.code == "Slash") { // Toggle help menu
        e.preventDefault()
        document.getElementById("key-shortcuts-summary").toggleAttribute("open")
        document.getElementById("key-shortcuts-summary").scrollIntoView()
    } else if (e.ctrlKey && e.altKey && e.code == "KeyS") { // Take screenshot of output
        e.preventDefault()
        SaveVideo()
    } else if (e.ctrlKey && e.altKey && e.code == "KeyV") { // Start/end video recording
        e.preventDefault()
        ToggleRecording()
    }
})

export {canvas, gl as glContext, programCounterHTML, currentInstructionHTML, peekRegHTML, peekRegInput}
