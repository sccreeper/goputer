import * as ioStyles from "url:./io.css";
import globals from "../globals";

class IOSwitch extends HTMLElement {

    /** @type {boolean} */
    #enabled = false;
    /** @type {string} */
    #reg;

    constructor() {
        super()
    }

    get enabled() {
        return this.#enabled
    }

    set enabled(val) {
        this.setAttribute("enabled", val)
        this.setAttribute("aria-checked", val)

        if (globals.vmIsAlive) {
            globals.switchQueue.push({
                register: this.reg,
                enabled: this.enabled,
            })   
        }

        this.#enabled = val
        this.#updateLabel()

        const changeEvent = new CustomEvent("change", {detail: val})
        this.dispatchEvent(changeEvent)
    }

    get reg() {
        return this.#reg
    }

    set reg(val) {
        this.#reg = val
        this.#updateLabel()
    }

    #updateLabel() {
        this.setAttribute("aria-label", `light ${this.reg} is ${this.enabled ? "on" : "off"}`)
    }

    connectedCallback() {

        this.reg = this.getAttribute("reg") ?? "io00"
        this.enabled = this.getAttribute("enabled") ?? false

        const shadow = this.attachShadow({mode: "open"})

        const stylesheet = document.createElement("link")
        stylesheet.rel = "stylesheet"
        stylesheet.href = ioStyles

        const btn = document.createElement("div")
        btn.classList.add("switch")

        this.addEventListener("click", () => {

            if (globals.vmIsAlive) {
                this.enabled = !this.enabled                
            }

        })

        this.addEventListener("keydown", (e) => {
            if (e.key === "Enter") {
                this.click()
            }
        })

        shadow.appendChild(stylesheet)
        shadow.appendChild(btn)

        this.tabIndex = 0
        this.setAttribute("enabled", this.enabled)
        this.setAttribute("reg", this.reg)
        this.setAttribute("role", "switch")
        this.setAttribute("aria-checked", this.enabled)

    }

    /**
     * @param {string} name 
     * @param {string} oldValue 
     * @param {string} newValue 
     */
    attributeChangedCallback(name, oldValue, newValue) {

        if (oldValue != newValue) {
            switch (name) {
                case "enabled":
                    this.enabled = newValue === "true"
                    break;
                case "reg":
                    this.reg = newValue
                    break;
                default:
                    break;
            }
        }

    }

    static observedAttributes = ["enabled", "reg"]

}

class IOLight extends HTMLElement {
    /** @type {boolean} */
    #enabled = false;
    /** @type {string} */
    #reg;

    constructor() {
        super()
    }

    get enabled() {
        return this.#enabled
    }

    set enabled(val) {
        this.setAttribute("enabled", val)
        this.#enabled = val
        this.#updateLabel()
    }

    get reg() {
        return this.#reg
    }

    set reg(val) {
        this.#reg = val
        this.#updateLabel()
    }

    #updateLabel() {
        this.setAttribute("aria-label", `light ${this.reg} is ${this.enabled ? "on" : "off"}`)
    }

    connectedCallback() {

        this.reg = this.getAttribute("reg") ?? "io00"
        this.enabled = this.getAttribute("enabled") ?? false

        const shadow = this.attachShadow({mode: "open"})

        const stylesheet = document.createElement("link")
        stylesheet.rel = "stylesheet"
        stylesheet.href = ioStyles

        const bulb = document.createElement("div")
        bulb.classList.add("bulb")

        shadow.appendChild(stylesheet)
        shadow.appendChild(bulb)

        this.setAttribute("enabled", this.enabled)
        this.setAttribute("reg", this.reg)
        this.#updateLabel()
        
    }

    /**
     * @param {string} name 
     * @param {string} oldValue 
     * @param {string} newValue 
     */
    attributeChangedCallback(name, oldValue, newValue) {

        if (oldValue != newValue) {
            switch (name) {
                case "enabled":
                    this.enabled = newValue === "true"
                    break;
                case "reg":
                    this.reg = newValue
                    break;
                default:
                    break;
            }
        }

    }

    static observedAttributes = ["enabled", "reg"]
}

customElements.define("io-switch", IOSwitch)
customElements.define("io-light", IOLight)

export default {};