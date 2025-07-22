import { db, fileTableName } from "./db";
import { CodeTabElement } from "./editor/code_tab"
import globals from "./globals"
import { goputer } from "./goputer"
import { OpenSharedArchive } from "./sharing";
import { clamp } from "./util";

const rowLength = 16;

export const tabsContainer = document.getElementById("code-names-container");

const newFile = document.getElementById("new-file");
newFile.addEventListener("click", NewFileUI)

// Editor elements

/** @type {HTMLDivElement} */
const codeEditorDiv = document.getElementById("code-editor")

/** @type {HTMLTextAreaElement} */
export const codeArea = document.getElementById("code-textarea");

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

/** @type {HTMLButtonElement} */
const deleteAllButton = document.getElementById("delete-all-button");

/** @type {Map<string, {blob: Blob, objectUrl: string, dimensions: number[]}>} */
export const imageMap = new Map();

deleteAllButton.addEventListener("click", (e) => {

    if (confirm("Are you sure you want permanently to delete all files?")) {
        
        RemoveAll()

        NewFile("main.gpasm")
        SwitchFocus("main.gpasm")

    }

});

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

            if (f.size > goputer.usableMemorySize) {
                if (!confirm(
                    `This file is ${f.size} bytes in size which is larger than goputer's usable memory size of ${goputer.usableMemorySize} bytes. Press OK if you know what you are doing.`
                )) {
                    return;    
                }
            }

            if (f.name.split(".").pop() == "zip") {
                OpenSharedArchive(f)
                return
            }

            let fileName = ""

            if (goputer.files.exists(f.name)) {
                fileName = `dup_${f.name}`
            } else {
                fileName = f.name
            }

            NewFile(fileName)

            // Determine extension and thus filetype

            /** @type {import("./editor/code_tab").FileType} */
            let fileType = ""

            switch (fileName.split(".").pop()) {
                case "txt":
                case "gpasm":
                    fileType = "text"
                    break;
                case "png":
                case "jpg":
                    fileType = "image"
                    InitImage(new Blob([f]), fileName, () => {SwitchFocus(fileName)})
                    break;
                default:
                    fileType = "bin"
                    break;
            }

            goputer.files.update(fileName, await f.bytes(), f.size, fileType, true)

            if (fileType != "image") {
                SwitchFocus(fileName)                
            }
        }
    }
)

codeEditorDiv.addEventListener("dragover", (e) => {
    e.preventDefault()
})

/**
 * Deletes all files and performs required cleanup.
 */
export function RemoveAll() {
    goputer.files.fileNames.forEach(fileName => {

        if (goputer.files.type(fileName) == "image") {
            window.URL.revokeObjectURL(imageMap.get(fileName).objectUrl)
        }

        goputer.files.remove(fileName)

        document.querySelector(`code-tab[filename="${fileName}"]`).remove()
    });

    imageMap.clear()
}

/**
 * 
 * @param {Blob} imageBlob 
 * @param {string} key 
 * @param {() => void|undefined} loadCallback
 */
export function InitImage(imageBlob, key, loadCallback = undefined) {
    const imgUrl = window.URL.createObjectURL(imageBlob)

    const image = new Image()
    
    image.onload = (e) => {
        imageMap.set(
            key,
            {
                blob: imageBlob,
                objectUrl: imgUrl,
                dimensions: [image.width, image.height]
            }
        )

        if (typeof loadCallback !== "undefined") {
            loadCallback()
        }
    }

    image.src = imgUrl
}

// Creation of a file from the UI.
export function NewFileUI(e) {
    
    let fileName = prompt("New file name:", `new_${goputer.files.numFiles+1}.gpasm`)

    if (fileName == null || fileName == "" || fileName == ".gpasm") {
        fileName = `new_${goputer.files.numFiles+1}.gpasm`
    }

    NewFile(fileName)
    SwitchFocus(fileName)
}

/**
 * Initializes a new file
 * @param {String} fileName 
 * @param {boolean} [doCreation=true] - defaults to true
 */
export function NewFile(fileName, doCreation = true) {

    if (doCreation) {
        goputer.files.update(fileName, new Uint8Array(), 0, "text", true)   
    }

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

    tabsContainer.insertBefore(newTab, newFile);

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
    
    let fileElements = tabsContainer.children

    for (let i = 0; i < fileElements.length; i++) {
        
        if (fileElements[i].getAttribute("filename") != globals.focusedFile) {
            fileElements[i].setAttribute("selected", "false")
        } else {
            fileElements[i].setAttribute("selected", "true")
        }
        
    }

}