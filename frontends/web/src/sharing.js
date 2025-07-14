import { ErrorTypes, ShowError } from "./error";
import globals from "./globals";
import { goputer } from "./goputer";
import { NewFile } from "./imports";

// Extracts shared code from URL
export function GetSharedCode() {

    if (window.location.search.length == 0) {
        return;
    }

    let base64Sting = window.location.search;
    base64Sting = base64Sting.substring(3, base64Sting.length);

    let codeJson = JSON.parse(atob(base64Sting))

    for (const [key, value] of Object.entries(codeJson)) {

        if (key != "main.gpasm") {
            NewFile(key)
        }

        let encoder = new TextEncoder()
        let encoded = encoder.encode(atob(value))

        goputer.files.update(key, encoded, encoded.length)
        document.getElementById("code-textarea").value = atob(value);

    }

    ShowError(ErrorTypes.Success, "Imported code from shared URL");

}

// Converts text area to base64 and copies shareable URL to clipboard.
export function ShareCode(e) {

    let files = goputer.files.fileNames

    var files_object = {}

    files.forEach(element => {
        files_object[element] = btoa(goputer.files.get(element))
    });

    let code_json = JSON.stringify(files_object)

    // Convert to base64

    let base64 = btoa(code_json);

    let port = ("" == window.location.port) ? "" : `:${window.location.port}`

    let shareable_url = `${window.location.protocol}//${window.location.hostname}${port}/?c=${base64}`

    navigator.clipboard.writeText(shareable_url);

    ShowError(ErrorTypes.Success, "Code copied to clipboard");
}

// Download program bytes.
export function DownloadProgram(e) {

    // Convert method return value to bytes first.

    let programBytesArray = goputer.getProgramBytes()

    let programBytes = new Uint8Array(programBytesArray.length)

    for (let i = 0; i < programBytesArray.length; i++) {
        programBytes[i] = programBytesArray[i]
    }

    let date = new Date()

    if (!globals.codeHasBeenCompiled) {
        return
    }

    let blob = new Blob([programBytes], { type: "application/octet-stream" })
    let link = document.createElement("a")
    link.href = window.URL.createObjectURL(blob)

    let filename = `program_${date.getHours().toString().padStart(2, "0")}${date.getMinutes().toString().padStart(2, "0")}${date.getSeconds().toString().padStart(2, "0")}`

    link.download = filename;
    link.click();

}

export function UploadBinary(e) {

    let uploadForm = document.createElement("input")
    uploadForm.type = "file"
    uploadForm.accept = ".gp"
    uploadForm.multiple = false

    uploadForm.addEventListener("change", async (e) => {
    
        let file = uploadForm.files[0]

        let fileBytes = new Uint8Array(await file.arrayBuffer())
        console.log(`Read file with ${fileBytes.length} byte(s)`)

        goputer.setProgramBytes(fileBytes, fileBytes.length)

        document.getElementById("run-code-button").disabled = false
        document.getElementById("download-code-button").disabled = false
        document.getElementById("code-textarea").disabled = true

        globals.codeHasBeenCompiled = true
        globals.compileFailed = false

        // Hide code editor

        document.getElementById("code-editor").style.visibility = "hidden"
        document.getElementById("binary-message").style.display = "block"
        document.getElementById("compile-code-button").disabled = true
        document.getElementById("share-code-button").disabled = true

        ShowError(ErrorTypes.Success, "Binary uploaded successfully")

    })

    uploadForm.click()
    
}