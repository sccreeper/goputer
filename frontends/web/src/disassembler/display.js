// Display disassembled code in UI.

import { goputer } from "../goputer";
import { instructionsContainer, definitionsContainer, interruptTableContainer } from "./main";
import shared from "./shared";

/**
 * 
 * @param {any} instructionObj 
 * @param {number} offset 
 * @returns 
 */
function GenerateInstructionHTML(instructionObj, offset) {

    let p = document.createElement("p")
    p.classList.add("font-fira")

    let addrElement = document.createElement("span")
    addrElement.classList.add("instruction-address")
    addrElement.textContent = `0x${offset.toString(16).padStart(8, "0").toUpperCase()} `
    addrElement.id = "addr_" + addrElement.textContent.trim()

    let itnElement = document.createElement("span")
    itnElement.classList.add("text-green-500")
    itnElement.textContent = instructionArray[instructionObj.instruction] + " "

    let argsElement = null

    if (instructionArray[instructionObj.instruction] == "jmp" || instructionArray[instructionObj.instruction] == "cndjmp" || instructionArray[instructionObj.instruction] == "call" || instructionArray[instructionObj.instruction] == "cndcall") {
        argsElement = document.createElement("a")
        argsElement.style.textDecoration = "underline"
        argsElement.href = `#addr_${instructionObj.data[0]}`
    } else if (instructionArray[instructionObj.instruction] == "lda" || instructionArray[instructionObj.instruction] == "sta") {
        argsElement = document.createElement("a")
        argsElement.style.textDecoration = "underline"
        argsElement.href = `#def_${instructionObj.data[0]}`
    } else {
        argsElement = document.createElement("span")
    }

    argsElement.classList.add("text-cyan-600")

    if (instructionObj.data != null) {
        instructionObj.data.forEach(element => {
            argsElement.textContent += element + " ";
        });
    }

    p.appendChild(addrElement)
    p.appendChild(itnElement)
    p.appendChild(argsElement)

    return p

}

export function DisplayDisassembledCode(codeObject) {

    shared.file_json = JSON.stringify(codeObject, null, 3);

    // Display instructions

    if (codeObject.instructions.length == 0) {
        instructionsContainer.querySelector("p").textContent = "No instructions"
    } else {
        instructionsContainer.replaceChildren();

        for (let i = 0; i < codeObject.instructions.length; i++) {
            const element = codeObject.instructions[i];
            
            instructionsContainer.appendChild(
                GenerateInstructionHTML(
                    element,
                    (i*5) + memOffset + codeObject.start_indexes[2], 
                )
            );
        }

    }

    // Display interrupt subscriptions

    interruptTableContainer.replaceChildren();

    var interruptTable = document.createElement("table")
    interruptTable.classList.add("w-full");

    for (const [key, value] of Object.entries(codeObject.interrupt_table)) {

        let interruptRow = document.createElement("tr")

        let interruptName = document.createElement("td")
        interruptName.textContent = interruptArray[parseInt(key)] + " "
        interruptRow.appendChild(interruptName)

        let interruptAddress = document.createElement("td")
        interruptAddress.classList.add("font-fira")
        interruptAddress.textContent = goputer.util.convertHex(value, false)
        interruptRow.appendChild(interruptAddress)

        interruptTable.appendChild(interruptRow)

    }

    interruptTableContainer.appendChild(interruptTable);

    // Display definitions

    if (codeObject.program_definitions.length == 0) {

        definitionsContainer.querySelector("p").textContent = "No definitions."

    } else {
        definitionsContainer.replaceChildren();

        let totalDefinitionLength = 0;

        for (let i = 0; i < codeObject.program_definitions.length; i++) {
            const element = codeObject.program_definitions[i];

            let definitionAddress = "0x" + (memOffset + codeObject.start_indexes[1] + totalDefinitionLength).toString(16).padStart(8, "0").toUpperCase();
            totalDefinitionLength += atob(element).length + 4;

            let definitionContainer = document.createElement("div")
            definitionContainer.classList.add("p-2")

            let definitionInfo = document.createElement("p");
            definitionInfo.textContent = `${atob(element).length} bytes at ${definitionAddress}`
            definitionInfo.id = `def_${definitionAddress}`
            definitionInfo.classList.add("definition-info")

            let definitionText = document.createElement("h3")
            definitionText.classList.add("font-bold");
            definitionText.textContent = "Text representation"

            let definition = document.createElement("p")
            definition.classList.add("table-code")

            definition.textContent = atob(element)

            let definitionTitle1 = document.createElement("h3")
            definitionTitle1.classList.add("font-bold")
            definitionTitle1.textContent = "Hex representation"

            let definitionHex = document.createElement("p")
            definitionHex.classList.add("table-code")

            // https://stackoverflow.com/a/41106346/7345078
            let intArray = Uint8Array.from(atob(element), c => c.charCodeAt(0));
            
            intArray.forEach(element => {
                
                definitionHex.textContent += element.toString(16).toUpperCase().padStart(2, "0");

            });

            definitionContainer.appendChild(definitionInfo)
            definitionContainer.appendChild(definitionText)
            definitionContainer.appendChild(definition)
            definitionContainer.appendChild(definitionTitle1)
            definitionContainer.appendChild(definitionHex)

            definitionsContainer.appendChild(definitionContainer)
            
        }

    }

}