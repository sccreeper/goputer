import { CodeTabElement } from "./editor/code_tab"
import globals from "./globals"
import { goputer } from "./goputer"
import { clamp } from "./util";

const rowLength = 16;

const filesContainer = document.getElementById("code-names-container");

const newFile = document.getElementById("new-file");
newFile.addEventListener("click", NewFileUI)

/** @type {HTMLTextAreaElement} */
const codeArea = document.getElementById("code-textarea");

/** @type {HTMLDivElement} */
const binDisplay = document.getElementById("bin-display");
/** @type {HTMLDivElement} */
const binDisplayData = document.getElementById("bin-display-data");

/** @type {HTMLInputElement} */
const binDisplayInput = document.getElementById("bin-display-offset-input");

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

codeArea.addEventListener("drop", 
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
                    break;
                default:
                    fileType = "bin"
                    break;
            }

            goputer.files.update(filename, await f.bytes(), f.size, fileType)
            SwitchFocus(filename)

        }
    }
)

codeArea.addEventListener("dragover", (e) => {
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
            codeArea.style.display = "block";

            let textDecoder = new TextDecoder()
            let decodedText = textDecoder.decode(fileBytes)
            codeArea.value = decodedText

            break;
        case "image":
            break;
        case "bin":
            codeArea.style.display = "none";
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