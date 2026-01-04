import { goputer } from "../goputer";
import { setVersion } from "../shared/version";
import { DisplayDisassembledCode } from "./display";
import shared from "./shared";

var uploadedBytes = []
var fileUploaded = false
var file_disassembled = false
var fileName = ""

goputer.workerInit();

setVersion()

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

ButtonDisassemble.addEventListener("click", async (e) => {

    if (String.fromCharCode(...uploadedBytes.slice(0, 4)) != "GPTR") {
        alert("Invalid file!")
        return
    }

    let code = await goputer.disassembleCode(uploadedBytes)
    
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