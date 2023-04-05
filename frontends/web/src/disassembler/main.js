import { DisplayDisassembledCode } from "./display";
import shared from "./shared";

var uploaded_bytes = []
var file_uploaded = false
var file_disassembled = false
var file_name = ""

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

export const ContainerInstructions = document.getElementById("container-instructions")
export const ContainerJumpBlocks = document.getElementById("container-jump-blocks")
export const InterruptTableContainer = document.getElementById("container-interrupt-table")
export const DefinitionsContainer = document.getElementById("container-definitions")

FileForm.addEventListener("change", (e) => {

    file_name = e.target.files[0].name;

    var reader = new FileReader()

    reader.onload = function () { 

        var arrayBuffer = this.result,
        array = new Uint8Array(arrayBuffer)

        uploaded_bytes = array

    }

    reader.readAsArrayBuffer(e.target.files[0])
    file_uploaded = true;

    ButtonDisassemble.removeAttribute("disabled")

})

ButtonUpload.addEventListener("click", (e) => {
    FileForm.click()
})

ButtonDisassemble.addEventListener("click", (e) => {

    if (String.fromCharCode(...uploaded_bytes.slice(0, 4)) != "GPTR") {
        alert("Invalid file!")
        return
    }

    let code_string = disassembleCode(uploaded_bytes)
    
    DisplayDisassembledCode(code_string)

    file_disassembled = true;

    ButtonDownload.removeAttribute("disabled")

})

ButtonDownload.addEventListener("click", (e) => {

    var file_element = document.createElement("a")
    file_element.setAttribute("href", `data:text/plain;charset=utf-8,${encodeURIComponent(shared.file_json)}`)
    file_element.setAttribute("download", `disassembled_${file_name}.json`)
    file_element.style.display = "none";

    document.body.appendChild(file_element)
    file_element.click()
    document.body.removeChild(file_element)

})