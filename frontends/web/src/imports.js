import { CodeTabElement } from "./editor/code_tab"
import globals from "./globals"
import { goputer } from "./goputer"
import { clamp } from "./util";

const rowLength = 16;

const filesContainer = document.getElementById("code-names-container");

const newFile = document.getElementById("new-file");
newFile.addEventListener("click", NewFileUI)

// Editor elements

/** @type {HTMLDivElement} */
const codeEditorDiv = document.getElementById("code-editor")

/** @type {HTMLTextAreaElement} */
const codeArea = document.getElementById("code-textarea");

/** @type {HTMLDivElement} */
const binDisplay = document.getElementById("bin-display");
/** @type {HTMLDivElement} */
const binDisplayData = document.getElementById("bin-display-data");

/** @type {HTMLInputElement} */
const binDisplayInput = document.getElementById("bin-display-offset-input");

/** @type {HTMLDivElement} */
const imageDisplay = document.getElementById("img-display");
/** @type {HTMLImageElement} */
const imageDisplayImage = document.getElementById("img-display-img");
/** @type {HTMLParagraphElement} */
const imageDisplayInfo = document.getElementById("img-display-info");
/** @type {HTMLInputElement} */
const imageDisplayTrueSize = document.getElementById("img-display-true-size");

/** @type {Map<string, {blob: Blob, objectUrl: string, dimensions: number[]}>} */
export const imageMap = new Map();

imageDisplayTrueSize.addEventListener("change", 
    /** @param {Event} e  */
    (e) => {
        if (imageDisplayTrueSize.checked) {
            imageDisplayImage.style.height = "auto";
        } else {
            imageDisplayImage.style.height = "100%";
        }
    }
)

binDisplayInput.addEventListener("input", 
    /** @param {Event} e  */
    (e) => {

        const realOffset = clamp(binDisplayInput.value, 0, goputer.files.size(globals.focusedFile)-1)
        const elementOffset = Math.floor(
            realOffset / rowLength
        )

        if (elementOffset != 0) {
            binDisplayData.children[elementOffset-1].scrollIntoView();
        } else {
            binDisplayData.scroll(0, 0)
        }

        // Highlight element

        const selection = window.getSelection()
        selection.removeAllRanges()

        const range = document.createRange()
        range.setStart(
            binDisplayData.children[elementOffset].children[1].childNodes[0], 
            (realOffset % rowLength) * 2
        )
        range.setEnd(
            binDisplayData.children[elementOffset].children[1].childNodes[0], 
            ((realOffset % rowLength) * 2) + 2
        )

        selection.addRange(range)

    }
)

codeArea.addEventListener("input", (e) => {

    let encoder = new TextEncoder();
    let encoded = encoder.encode(codeArea.value);

    goputer.files.update(globals.focusedFile, encoded, encoded.length, "text");
})

codeEditorDiv.addEventListener("drop", 
    /** @type {DragEvent} */
    async (e) => {
        e.preventDefault();

        if (e.dataTransfer.files.length > 0) {
        
            const f = e.dataTransfer.files[0]
            let filename = ""

            if (goputer.files.exists(f.name)) {
                filename = `dup_${f.name}`
            } else {
                filename = f.name
            }

            NewFile(filename)

            // Determine extension and thus filetype

            /** @type {import("./editor/code_tab").FileType} */
            let fileType = ""

            switch (filename.split(".").pop()) {
                case "txt":
                case "gpasm":
                    fileType = "text"
                    break;
                case "png":
                case "jpg":
                    fileType = "image"

                    const imgObjectUrl = window.URL.createObjectURL(f)

                    const loadedImage = new Image()
                    loadedImage.src = imgObjectUrl

                    loadedImage.onload = (e) => {
                        imageMap.set(
                            filename,
                            {
                                blob: f,
                                objectUrl: imgObjectUrl,
                                dimensions: [loadedImage.width, loadedImage.height]
                            },
                        ) 
                        SwitchFocus(filename) 
                    }


                    break;
                default:
                    fileType = "bin"
                    break;
            }

            goputer.files.update(filename, await f.bytes(), f.size, fileType, true)
            SwitchFocus(filename)

        }
    }
)

codeEditorDiv.addEventListener("dragover", (e) => {
    e.preventDefault()
})

// Creation of a file from the UI.
export function NewFileUI(e) {
    
    let file_name = prompt("New file name:", `new_${goputer.files.numFiles+1}.gpasm`)

    if (file_name == null || file_name == "" || file_name == ".gpasm") {
        file_name = `new_${goputer.files.numFiles+1}.gpasm`
    }

    NewFile(file_name)
}

/**
 * Initializes a new file
 * @param {String} fileName 
 */
export function NewFile(fileName) {

    goputer.files.update(fileName, new Uint8Array(), 0, "text")

    /** @type {CodeTabElement} */
    let newTab = document.createElement("code-tab")
    newTab.filename = fileName

    if (fileName === "main.gpasm") {
        newTab.deletable = false
        newTab.renamable = false
    }
    newTab.type = "text"

    newTab.addEventListener(
        "fileselect",
        /** @param {CustomEvent} e  */
        (e) => {
            SwitchFocus(e.detail)
        }
    )

    newTab.addEventListener(
        "filedelete",
        /** @param {CustomEvent} e  */
        (e) => {
            DeleteFile(e.detail)
        }
    )

    filesContainer.insertBefore(newTab, newFile);

    globals.focusedFile = fileName;

    codeArea.value = ""

    SwitchFocusedStyle();

}

export function DeleteFile(fileName) {

    SwitchFocus("main.gpasm")

}

// Switch focus from one file to another
/**
 * 
 * @param {string} fileName 
 */
export function SwitchFocus(fileName) {

    const fileBytes = new Uint8Array(goputer.files.size(fileName))
    goputer.files.get(fileName, fileBytes)

    switch (goputer.files.type(fileName)) {
        case "text":
            binDisplay.style.display = "none";
            imageDisplay.style.display = "none";
            codeArea.style.display = "block";

            let textDecoder = new TextDecoder()
            let decodedText = textDecoder.decode(fileBytes)
            codeArea.value = decodedText

            break;
        case "image":
            codeArea.style.display = "none";
            binDisplay.style.display = "none";
            imageDisplay.style.display = "grid";

            imageDisplayImage.src = imageMap.get(fileName).objectUrl
            imageDisplayInfo.innerText = `Dimensions: ${imageMap.get(fileName).dimensions[0]}x${imageMap.get(fileName).dimensions[1]}
Original size: ${imageMap.get(fileName).blob.size} bytes Encoded size: ${goputer.files.size(fileName)} bytes`

            break;
        case "bin":
            codeArea.style.display = "none";
            imageDisplay.style.display = "none";
            binDisplay.style.display = "block";

            binDisplayData.replaceChildren();

            for (let i = 0; i < fileBytes.length; i += rowLength) {
                
                let hexString = ""

                fileBytes.slice(i, i + rowLength).forEach(element => {
                    hexString += element.toString(16).padStart(2, "0")
                });

                // Create elements

                const container = document.createElement("div")
                
                const offsetAddress = document.createElement("span")
                offsetAddress.innerText = i.toString(16).padStart(8, "0")
                offsetAddress.title = i.toString()
                
                const row = document.createElement("span")
                row.innerText = hexString

                container.appendChild(offsetAddress)
                container.appendChild(row)

                binDisplayData.appendChild(container)

            }

            break;
        default:
            break;
    }

    globals.focusedFile = fileName;
    SwitchFocusedStyle();

}

export function SwitchFocusedStyle() {
    
    let fileElements = filesContainer.children

    for (let i = 0; i < fileElements.length; i++) {
        
        if (fileElements[i].getAttribute("filename") != globals.focusedFile) {
            fileElements[i].setAttribute("selected", "false")
        } else {
            fileElements[i].setAttribute("selected", "true")
        }
        
    }

}