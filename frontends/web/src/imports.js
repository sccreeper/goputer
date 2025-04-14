import globals from "./globals"
import { files_container, new_file } from "./init"


// Creation of a file from the UI.
export function NewFileUI(e) {
    
    let file_name = prompt("New file name:", `new_${getFiles().length+1}.gpasm`)

    if (file_name == null || file_name == "" || file_name == ".gpasm") {
        file_name = `new_${getFiles().length+1}.gpasm`
    }

    NewFile(file_name)
}

// Initializes a new file, asking the user for it's name.
export function NewFile(file_name) {

    updateFile(file_name, "")

    let tab_div = document.createElement("div")
    tab_div.classList.add("code-name")
    tab_div.setAttribute("data-selected", "true")
    tab_div.setAttribute("data-file-name", file_name)

    tab_div.addEventListener("click", SwitchFocus)

    let file_name_p = document.createElement("p")
    file_name_p.textContent = file_name;

    tab_div.appendChild(file_name_p)

    let delete_file_i = document.createElement("i")
    delete_file_i.classList.add("bi", "bi-x")
    delete_file_i.title = "Delete file"

    delete_file_i.addEventListener("click", DeleteFile)

    tab_div.appendChild(delete_file_i)

    files_container.insertBefore(tab_div, new_file);

    globals.focusedFile = file_name;

    document.getElementById("code-textarea").value = ""

    SwitchFocusedStyle();


}

export function DeleteFile(e) {
    
    let file_name = e.currentTarget.parentElement.getAttribute("data-file-name");
    globals.focusedFile = e.currentTarget.previousSibling.getAttribute("data-file-name");

    e.currentTarget.parentElement.remove()
    removeFile(file_name)

    document.getElementById("code-textarea").value = getFile(globals.focusedFile)
    SwitchFocusedStyle();

}

// Switch focus from one file to another
export function SwitchFocus(e) {
    
    let file_name = e.currentTarget.getAttribute("data-file-name");
    document.getElementById("code-textarea").value = getFile(file_name)

    globals.focusedFile = file_name;
    SwitchFocusedStyle();

}

export function SwitchFocusedStyle() {
    
    let file_elements = files_container.children

    for (let i = 0; i < file_elements.length; i++) {
        
        if (file_elements[i].getAttribute("data-file-name") != globals.focusedFile) {
            file_elements[i].setAttribute("data-selected", "false")
        } else {
            file_elements[i].setAttribute("data-selected", "true")
        }
        
    }

}