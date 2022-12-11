import { renderContext, canvas } from "./init";
import global from "./globals.js";

// Main app logic

export function Compile() {  

    compileCode(document.getElementById("code-textarea").value)
    codeHasBeenCompiled = true;

}

export function Run() { 

    if (!global.codeHasBeenCompiled) {

        console.error("No code has been compiled!");

    } else {

        runInterval = setInterval(Cycle, 1000 / FPS);

    }

}

//Performs one cycle of the VM
export function Cycle() {  

    if (!global.vmIsAlive) {
        
        console.error("VM isn't alive therefore can't run code.");

    } else {

    }

}