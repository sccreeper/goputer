import { CodeTabElement } from "./editor/code_tab"
import globals from "./globals"
import { goputer } from "./goputer"
import { filesContainer, newFile } from "./init"

// Creation of a file from the UI.
export function NewFileUI(e) {
    
    let file_name = prompt("New file name:", `new_${goputer.files.numFiles+1}.gpasm`)

    if (file_name == null || file_name == "" || file_name == ".gpasm") {
        file_name = `new_${goputer.files.numFiles+1}.gpasm`
    }

    NewFile(file_name)
}

/**
 * Initializes a new file, asking the user for it's name.
 * @param {String} fileName 
 */
export function NewFile(fileName) {

    goputer.files.update(fileName, new Uint8Array(), 0)

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

    document.getElementById("code-textarea").value = ""

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

    let dest = new Uint8Array(goputer.files.size(fileName))
    goputer.files.get(fileName, dest)

    let textDecoder = new TextDecoder()
    let decodedText = textDecoder.decode(dest)

    document.getElementById("code-textarea").value = decodedText

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