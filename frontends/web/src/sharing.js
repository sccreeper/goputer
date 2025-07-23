import { BlobReader, BlobWriter, Uint8ArrayReader, Uint8ArrayWriter, ZipReader, ZipWriter } from "@zip.js/zip.js";
import { ErrorTypes, ShowError } from "./error";
import globals from "./globals";
import { goputer } from "./goputer";
import { imageMap, InitImage, NewFile, RemoveAll, SwitchFocus } from "./imports";

/**
 * Open a shared archive. Expects all files to be in the top level directory, will not tree walk.
 * @param {Blob} fileBlob 
 * @returns 
 */
export async function OpenSharedArchive(fileBlob) {

    if (confirm("Importing this will delete all existing files. Are you sure you want to do this?")) {
        RemoveAll()
    } else {
        return
    }

    const zipFileReader = new BlobReader(fileBlob);

    const zipReader = new ZipReader(zipFileReader);

    for (const ent of await zipReader.getEntries()) {
        if (!ent.directory && ent.filename.split("/").length == 1) {
            
            const fileWriter = new Uint8ArrayWriter();
            await ent.getData(fileWriter)
            const fileData = await fileWriter.getData()

            /** @type {import("./editor/code_tab").FileType} */
            let fileType;

            switch (ent.filename.split(".").pop()) {
                case "txt":
                case "gpasm":
                    fileType = "text"
                    break;
                case "png":
                case "jpg":
                    fileType = "image"

                    InitImage(new Blob([fileData]), ent.filename)

                    break;
                default:
                    fileType = "bin"
                    break;
            }

            goputer.files.update(
                ent.filename, 
                fileData,
                fileData.length,
                fileType,
                true
            )
            NewFile(ent.filename, false)
        }
    }

    SwitchFocus("main.gpasm");
    ShowError(ErrorTypes.Success, "Imported code from archive");

}

// Converts text area to base64 and copies shareable URL to clipboard.
export async function DownloadAll(e) {

    const zipFileWriter = new BlobWriter();
    const zipWriter = new ZipWriter(zipFileWriter);

    for (const fileName of goputer.files.fileNames) {
        
        /** @type {Uint8Array} */
        let fileBytes;

        if (goputer.files.type(fileName) == "image") {

            fileBytes = new Uint8Array(await imageMap.get(fileName).blob.bytes())

        } else {

            fileBytes = new Uint8Array(goputer.files.size(fileName))
            goputer.files.get(fileName, fileBytes)

        }

        const fileReader = new Uint8ArrayReader(fileBytes)
        await zipWriter.add(fileName, fileReader)

    }

    await zipWriter.close()

    const objUrl = window.URL.createObjectURL(await zipFileWriter.getData())

    const date = new Date()

    /** @type {HTMLAnchorElement} */
    const linkElement = document.createElement("a")
    linkElement.href = objUrl
    linkElement.download = `files_${date.getHours().toString().padStart(2, "0")}${date.getMinutes().toString().padStart(2, "0")}${date.getSeconds().toString().padStart(2, "0")}.zip`
    linkElement.click()

    window.URL.revokeObjectURL(objUrl)

    ShowError(ErrorTypes.Success, "Zip file exported successfully.")

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

    let filename = `program_${date.getHours().toString().padStart(2, "0")}${date.getMinutes().toString().padStart(2, "0")}${date.getSeconds().toString().padStart(2, "0")}.gp`

    link.download = filename;
    link.click();

    window.URL.revokeObjectURL(blob);

}

export function UploadBinary(e) {

    let uploadForm = document.createElement("input")
    uploadForm.type = "file"
    uploadForm.accept = ".gp"
    uploadForm.multiple = false

    uploadForm.addEventListener("change", async (e) => {
    
        let file = uploadForm.files[0]

        const fileBytes = await file.bytes()

        if (new TextDecoder().decode(fileBytes.slice(0, 4)) != "GPTR") {
            alert("Invalid file.")
            return
        }

        fileBytes.slice(0, 4)

        console.log(`Read file with ${fileBytes.length} byte(s)`)

        goputer.setProgramBytes(fileBytes, fileBytes.length)

        document.getElementById("run-code-button").disabled = false
        document.getElementById("download-code-button").disabled = false

        globals.codeHasBeenCompiled = true
        globals.compileFailed = false

        // Hide code editor

        document.getElementById("code-editor").style.visibility = "hidden"
        document.getElementById("binary-message").style.display = "block"
        document.getElementById("compile-code-button").disabled = true
        document.getElementById("download-all-button").disabled = true

        ShowError(ErrorTypes.Success, "Binary uploaded successfully")

    })

    uploadForm.click()
    
}