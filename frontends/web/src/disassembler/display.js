// Display disassembled code in UI.

import { instructionsContainer, jumpBlocksContainer, definitionsContainer, interruptTableContainer } from "./main";
import shared from "./shared";

function GenerateInstructionHTML(instrutionObj) {

    let p = document.createElement("p")
    p.classList.add("font-fira")

    let itnElement = document.createElement("span")
    itnElement.classList.add("text-green-500")
    itnElement.textContent = instructionArray[instrutionObj.instruction] + " "

    let argsElement = document.createElement("span")
    argsElement.classList.add("text-cyan-600")

    if (instrutionObj.data != null) {
        instrutionObj.data.forEach(element => {
            argsElement.textContent += element + " ";
        });
    }

    p.appendChild(itnElement)
    p.appendChild(argsElement)

    return p

}

export function DisplayDisassembledCode(codeJSONString) {

    let codeObject = JSON.parse(codeJSONString)
    shared.file_json = JSON.stringify(codeObject, null, 3);

    // Display instructions

    if (codeObject.instructions.length == 0) {
        instructionsContainer.querySelector("p").textContent = "No instructions"
    } else {
        instructionsContainer.replaceChildren();

        codeObject.instructions.forEach(element => {

            instructionsContainer.appendChild(GenerateInstructionHTML(element));

        });
    }

    // Display jump blocks

    if (Object.keys(codeObject.jump_blocks).length == 0) {
        jumpBlocksContainer.querySelector("p").textContent = "No jump blocks"
    } else {
        jumpBlocksContainer.replaceChildren();

        for (const [key, value] of Object.entries(codeObject.jump_blocks)) {

            let jump_block_header = document.createElement("h3")
            jump_block_header.classList.add("font-bold");
            jump_block_header.textContent = `Block ${convertHex(parseInt(key), true)}`

            jumpBlocksContainer.appendChild(jump_block_header)

            value.forEach(element => {
                jumpBlocksContainer.appendChild(GenerateInstructionHTML(element));
            });

        }
    }

    // Display interrupt subscriptions

    interruptTableContainer.replaceChildren();

    var interrupt_table = document.createElement("table")
    interrupt_table.classList.add("w-full");

    for (const [key, value] of Object.entries(codeObject.interrupt_table)) {

        let interrupt_row = document.createElement("tr")

        let interrupt_name = document.createElement("td")
        interrupt_name.textContent = interruptArray[parseInt(key)] + " "
        interrupt_row.appendChild(interrupt_name)

        let interrupt_address = document.createElement("td")
        interrupt_address.classList.add("font-fira")
        interrupt_address.textContent = convertHex(value, false)
        interrupt_row.appendChild(interrupt_address)

        interrupt_table.appendChild(interrupt_row)

    }

    interruptTableContainer.appendChild(interrupt_table);

    // Display definitions

    if (codeObject.program_definitions.length == 0) {

        definitionsContainer.querySelector("p").textContent = "No definitions."

    } else {
        definitionsContainer.replaceChildren();

        codeObject.program_definitions.forEach(element => {

            let definitionContainer = document.createElement("div")
            definitionContainer.classList.add("border-b", "border-gray-400", "p-2")

            let definitionTitle = document.createElement("h3")
            definitionTitle.classList.add("font-bold");
            definitionTitle.textContent = "Text representation"

            let definition = document.createElement("p")
            definition.classList.add("table-code")

            definition.textContent = atob(element)

            let definitionTitle1 = document.createElement("h3")
            definitionTitle1.classList.add("font-bold")
            definitionTitle1.textContent = "Hex representation"

            let definitionHex = document.createElement("p")
            definitionHex.classList.add("table-code")

            // https://stackoverflow.com/a/41106346/7345078
            let int_array = Uint8Array.from(atob(element), c => c.charCodeAt(0));
            
            int_array.forEach(element => {
                
                definitionHex.textContent += element.toString(16).toUpperCase();

            });

            definitionContainer.appendChild(definitionTitle)
            definitionContainer.appendChild(definition)
            definitionContainer.appendChild(definitionTitle1)
            definitionContainer.appendChild(definitionHex)

            definitionsContainer.appendChild(definitionContainer)

        });
    }

}