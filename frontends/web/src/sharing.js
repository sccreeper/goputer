import { ErrorTypes, ShowError } from "./error";
import globals from "./globals";
import { NewFile } from "./imports";

// Extracts shared code from URL
export function GetSharedCode() {

    if (window.location.search.length == 0) {
        return;
    }

    let base64_string = window.location.search;
    base64_string = base64_string.substring(3, base64_string.length);

    let code_json = JSON.parse(atob(base64_string))

    for (const [key, value] of Object.entries(code_json)) {

        if (key != "main.gpasm") {
            NewFile(key)
        }
        updateFile(key, atob(value))
        document.getElementById("code-textarea").value = atob(value);

    }

    ShowError(ErrorTypes.Success, "Imported code from shared URL");

}

// Converts text area to base64 and copies shareable URL to clipboard.
export function ShareCode(e) {

    let files = getFiles()

    var files_object = {}

    files.forEach(element => {
        files_object[element] = btoa(getFile(element))
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

    let program_bytes_array = getProgramBytes()

    let program_bytes = new Uint8Array(program_bytes_array.length)

    for (let i = 0; i < program_bytes_array.length; i++) {
        program_bytes[i] = program_bytes_array[i]
    }

    let date = new Date()

    if (!globals.codeHasBeenCompiled) {
        return
    }

    let blob = new Blob([program_bytes], { type: "application/octet-stream" })
    let link = document.createElement("a")
    link.href = window.URL.createObjectURL(blob)

    let filename = `program_${date.getHours().toString().padStart(2, "0")}${date.getMinutes().toString().padStart(2, "0")}${date.getSeconds().toString().padStart(2, "0")}`

    link.download = filename;
    link.click();

}