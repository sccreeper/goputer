import { goputer } from "../goputer";
import { DisplayDisassembledCode } from "./display";
import shared from "./shared";

var uploadedBytes = []
var fileUploaded = false
var file_disassembled = false
var fileName = ""

const go = new Go();
await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

fetch("/ver")
.then((response) => response.text())
.then((data) => {

    let hash = data.split(/\r?\n/)[0];
    let time = data.split(/\r?\n/)[1];

    document.getElementById("version").textContent = `${hash.substring(0, 10)}`;
    document.getElementById("version").setAttribute("href", `https://github.com/sccreeper/goputer/commit/${hash}`);
    document.getElementById("build-date").textContent = time;

})

const ButtonUpload = document.getElementById("button-upload")
const ButtonDisassemble = document.getElementById("button-disassemble")
const ButtonDownload = document.getElementById("button-download")
const FileForm = document.getElementById("file-form")

export const instructionsContainer = document.getElementById("container-instructions")
export const interruptTableContainer = document.getElementById("container-interrupt-table")
export const definitionsContainer = document.getElementById("container-definitions")

FileForm.addEventListener("change", (e) => {

    fileName = e.target.files[0].name;

    var reader = new FileReader()

    reader.onload = function () { 

        var arrayBuffer = this.result,
        array = new Uint8Array(arrayBuffer)

        uploadedBytes = array

    }

    reader.readAsArrayBuffer(e.target.files[0])
    fileUploaded = true;

    ButtonDisassemble.removeAttribute("disabled")

})

ButtonUpload.addEventListener("click", (e) => {
    FileForm.click()
})

ButtonDisassemble.addEventListener("click", (e) => {

    if (String.fromCharCode(...uploadedBytes.slice(0, 4)) != "GPTR") {
        alert("Invalid file!")
        return
    }

    let code = goputer.disassembleCode(uploadedBytes)
    
    DisplayDisassembledCode(code)

    file_disassembled = true;

    ButtonDownload.removeAttribute("disabled")

})

ButtonDownload.addEventListener("click", (e) => {

    var file_element = document.createElement("a")
    file_element.setAttribute("href", `data:text/plain;charset=utf-8,${encodeURIComponent(shared.file_json)}`)
    file_element.setAttribute("download", `disassembled_${fileName}.json`)
    file_element.style.display = "none";

    document.body.appendChild(file_element)
    file_element.click()
    document.body.removeChild(file_element)

})