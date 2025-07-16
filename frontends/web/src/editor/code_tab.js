import * as tabStyles from "url:./tab.css";
import { goputer } from "../goputer";
import { NewFile, SwitchFocus } from "../imports";

/**
 * @typedef {"text"|"image"|"bin"} FileType
 */

export class CodeTabElement extends HTMLElement {

    /** @type {boolean} */
    #deletable = true;
    /** @type {string} */
    #filename;
    /** @type {FileType} */
    #type = "text";
    /** @type {boolean} */
    #selected;
    /** @type {boolean} */
    #renamable = true; 

    /**
     * 
     * @param {string} filename 
     * @param {boolean} deletable defaults to true
     * @param {FileType} fileType defaults to text
     */
    constructor() {
        super()
    }

    get selected() {
        return this.#selected;
    }

    set selected(val) {
        this.#selected = val;
        this.setAttribute("selected", this.#selected)
    }

    get type() {
        return this.#type;
    }

    set type(val) {
        this.#type = val;
        this.setAttribute("type", val)
    }

    get filename() {
        return this.#filename;
    }

    set filename(val) {
        this.#filename = val
        this.setAttribute("filename", val)
    }

    get deletable() {
        return this.#deletable
    }
    
    set deletable(val) {
        this.#deletable = val;
        this.setAttribute("deletable", val)
    }

    get renamable() {
        return this.#renamable;
    }

    set renamable(val) {
        this.#renamable = val;
        this.setAttribute("renamable", val)
    }

    connectedCallback() {

        const shadow = this.attachShadow({mode: "open"})

        // Instantiate element tree

        let stylesheet = document.createElement("link")
        stylesheet.rel = "stylesheet"
        stylesheet.href = tabStyles

        let bsIcons = document.createElement("link")
        bsIcons.rel = "stylesheet"
        bsIcons.href = new URL("npm:bootstrap-icons/font/bootstrap-icons.css", import.meta.url)

        let mainParent = document.createElement("div")
        mainParent.classList.add("code-name")

        let fileNameEl = document.createElement("p")
        fileNameEl.textContent = this.filename

        fileNameEl.addEventListener("click", (e) => this.#focusFile(e))

        if (this.renamable) {
            fileNameEl.addEventListener("dblclick", (e) => this.#renameFile(e))   
        }

        mainParent.appendChild(fileNameEl)

        if (this.deletable) {
            let deleteFileEl = document.createElement("i")
            deleteFileEl.classList.add("bi", "bi-x", "delete-file-button")
            deleteFileEl.title = "Delete file"
            deleteFileEl.addEventListener("click", (e) => this.#deleteFile(e))
            mainParent.appendChild(deleteFileEl)   
        }

        shadow.appendChild(stylesheet)
        shadow.appendChild(bsIcons)
        shadow.appendChild(mainParent)

        // Set attributes

        this.setAttribute("selected", this.selected)
        this.setAttribute("filename", this.filename)
        this.setAttribute("type", this.type)

    }

    #deleteFile(e) {

        if (confirm(`Are you sure you want to delete ${this.filename}?`)) {
            
            const deleteEvent = new CustomEvent("filedelete", {detail: this.filename, composed: true})
            this.dispatchEvent(deleteEvent)

            goputer.files.remove(this.filename)

            this.remove()

        }

    }

    #focusFile(e) {

        const focusEvent = new CustomEvent("fileselect", {detail: this.filename, composed: true})
        this.dispatchEvent(focusEvent)

        this.selected = true;

    }

    #renameFile(e) {

        let newFilename = prompt("Choose a new filename:", this.filename) ?? ""

        if (newFilename.trim().length == 0) {
            return
        }

        if (goputer.files.exists(newFilename)) {
            return
        }

        let fileSize = goputer.files.size(this.filename)
        let fileData = new Uint8Array(fileSize)
        goputer.files.get(this.filename, fileData)
        goputer.files.remove(this.filename)

        NewFile(newFilename)
        goputer.files.update(newFilename, fileData, fileData.length, )
        SwitchFocus(newFilename)

        const renameEvent = new CustomEvent("filerename", {detail: {oldName: this.filename, newName: newFilename }})
        this.dispatchEvent(renameEvent)

        this.remove()

    }

    /**
     * 
     * @param {string} name 
     * @param {string} oldValue 
     * @param {string} newValue 
     */
    attributeChangedCallback(name, oldValue, newValue) {

        if (name === "selected" && oldValue !== newValue) {
            this.selected = newValue === "true"
        } else if (name == "filename" && oldValue !== newValue) {
            this.filename = newValue
        } else if (name == "type" && oldValue !== newValue) {
            this.type = newValue
        }

    }

    static observedAttributes = ["selected", "filename", "type"];

    disconnectedCallback() {}
    
}

customElements.define("code-tab", CodeTabElement)