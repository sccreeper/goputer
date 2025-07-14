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

    let tab_div = document.createElement("div")
    tab_div.classList.add("code-name")
    tab_div.setAttribute("data-selected", "true")
    tab_div.setAttribute("data-file-name", fileName)

    tab_div.addEventListener("click", SwitchFocus)

    let file_name_p = document.createElement("p")
    file_name_p.textContent = fileName;

    tab_div.appendChild(file_name_p)

    let delete_file_i = document.createElement("i")
    delete_file_i.classList.add("bi", "bi-x")
    delete_file_i.title = "Delete file"

    delete_file_i.addEventListener("click", DeleteFile)

    tab_div.appendChild(delete_file_i)

    filesContainer.insertBefore(tab_div, newFile);

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
export function SwitchFocus(e) {
    
    let fileName = e.currentTarget.getAttribute("data-file-name");
    document.getElementById("code-textarea").value = goputer.files.get(fileName)

    globals.focusedFile = fileName;
    SwitchFocusedStyle();

}

export function SwitchFocusedStyle() {
    
    let file_elements = filesContainer.children

    for (let i = 0; i < file_elements.length; i++) {
        
        if (file_elements[i].getAttribute("data-file-name") != globals.focusedFile) {
            file_elements[i].setAttribute("data-selected", "false")
        } else {
            file_elements[i].setAttribute("data-selected", "true")
        }
        
    }

}