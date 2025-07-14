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

    let tabDiv = document.createElement("div")
    tabDiv.classList.add("code-name")
    tabDiv.setAttribute("data-selected", "true")
    tabDiv.setAttribute("data-file-name", fileName)

    tabDiv.addEventListener("click", SwitchFocus)

    let fileNameP = document.createElement("p")
    fileNameP.textContent = fileName;

    tabDiv.appendChild(fileNameP)

    let deleteFileI = document.createElement("i")
    deleteFileI.classList.add("bi", "bi-x", "delete-file-button")
    deleteFileI.title = "Delete file"

    deleteFileI.addEventListener("click", DeleteFile)

    tabDiv.appendChild(deleteFileI)

    filesContainer.insertBefore(tabDiv, newFile);

    globals.focusedFile = fileName;

    document.getElementById("code-textarea").value = ""

    SwitchFocusedStyle();

}

export function DeleteFile(e) {
    
    let fileName = e.currentTarget.parentElement.getAttribute("data-file-name");
    globals.focusedFile = e.currentTarget.previousSibling.getAttribute("data-file-name");

    e.currentTarget.parentElement.remove()
    goputer.files.remove(fileName)

    document.getElementById("code-textarea").value = goputer.files.get(globals.focusedFile)
    SwitchFocusedStyle();

}

// Switch focus from one file to another
/**
 * 
 * @param {MouseEvent} e 
 */
export function SwitchFocus(e) {

    let fileName = e.currentTarget.getAttribute("data-file-name");
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
        
        if (fileElements[i].getAttribute("data-file-name") != globals.focusedFile) {
            fileElements[i].setAttribute("data-selected", "false")
        } else {
            fileElements[i].setAttribute("data-selected", "true")
        }
        
    }

}