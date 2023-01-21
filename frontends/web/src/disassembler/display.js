// Display disassembled code in UI.

import { ContainerInstructions, ContainerJumpBlocks, DefinitionsContainer, InterruptTableContainer } from "./main";
import shared from "./shared";

function GenerateInstructionHTML(instruction_obj) {

    let p = document.createElement("p")
    p.classList.add("font-fira")

    let itn_el = document.createElement("span")
    itn_el.classList.add("text-green-500")
    itn_el.textContent = instructionArray[instruction_obj.instruction] + " "

    let args_el = document.createElement("span")
    args_el.classList.add("text-cyan-600")

    instruction_obj.data.forEach(element => {
        args_el.textContent += element + " ";
    });

    p.appendChild(itn_el)
    p.appendChild(args_el)

    return p

}

export function DisplayDisassembledCode(code_string) {

    let code_object = JSON.parse(code_string)
    shared.file_json = JSON.stringify(code_object, null, 3);

    // Display instructions

    if (code_object.instructions.length == 0) {
        ContainerInstructions.querySelector("p").textContent = "No instructions"
    } else {
        ContainerInstructions.replaceChildren();

        code_object.instructions.forEach(element => {

            ContainerInstructions.appendChild(GenerateInstructionHTML(element));

        });
    }

    // Display jump blocks

    if (Object.keys(code_object.jump_blocks).length == 0) {
        ContainerJumpBlocks.querySelector("p").textContent = "No jump blocks"
    } else {
        ContainerJumpBlocks.replaceChildren();

        for (const [key, value] of Object.entries(code_object.jump_blocks)) {

            let jump_block_header = document.createElement("h3")
            jump_block_header.classList.add("font-bold");
            jump_block_header.textContent = `Block ${convertHex(parseInt(key), true)}`

            ContainerJumpBlocks.appendChild(jump_block_header)

            value.forEach(element => {
                ContainerJumpBlocks.appendChild(GenerateInstructionHTML(element));
            });

        }
    }

    // Display interrupt subscriptions

    InterruptTableContainer.replaceChildren();

    var interrupt_table = document.createElement("table")
    interrupt_table.classList.add("w-full");

    for (const [key, value] of Object.entries(code_object.interrupt_table)) {

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

    InterruptTableContainer.appendChild(interrupt_table);

    // Display definitions

    if (code_object.program_definitions.length == 0) {

        DefinitionsContainer.querySelector("p").textContent = "No definitions."

    } else {
        DefinitionsContainer.replaceChildren();

        code_object.program_definitions.forEach(element => {

            let definition_container = document.createElement("div")
            definition_container.classList.add("border-b", "border-gray-400", "p-2")

            let definition_title = document.createElement("h3")
            definition_title.classList.add("font-bold");
            definition_title.textContent = "Text representation"

            let definition = document.createElement("p")
            definition.classList.add("table-code")

            definition.textContent = atob(element)

            let definition_title_1 = document.createElement("h3")
            definition_title_1.classList.add("font-bold")
            definition_title_1.textContent = "Hex representation"

            let definition_hex = document.createElement("p")
            definition_hex.classList.add("table-code")

            // https://stackoverflow.com/a/41106346/7345078
            let int_array = Uint8Array.from(atob(element), c => c.charCodeAt(0));
            
            int_array.forEach(element => {
                
                definition_hex.textContent += element.toString(16).toUpperCase();

            });

            definition_container.appendChild(definition_title)
            definition_container.appendChild(definition)
            definition_container.appendChild(definition_title_1)
            definition_container.appendChild(definition_hex)

            DefinitionsContainer.appendChild(definition_container)

        });
    }

}