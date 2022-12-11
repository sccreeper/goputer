import { renderContext, canvas } from "./init";
import global from "./globals.js";
import { clearCanvas, drawRect } from "./canvas_util";
import globals from "./globals.js";

// Main app logic

export function Compile() {  

    compileCode(document.getElementById("code-textarea").value)
    global.codeHasBeenCompiled = true;

}

export function Run() { 

    if (!global.codeHasBeenCompiled) {

        console.error("No code has been compiled!");

    } else {

        initVM();

        global.vmIsAlive = true;
        global.runInterval = setInterval(Cycle, Math.round(1000 / global.FPS));
        
    }

}

//Performs one cycle of the VM & Updates UI
export function Cycle() {
    
    if (isFinished()) {
        
        clearInterval(global.runInterval);
        return;

    }

    if (!global.vmIsAlive) {
        
        console.error("VM isn't alive therefore can't run code.");

    } else {

        var x = getInterrupt()

        switch (x) {
            case interruptInts["va"]:
                drawRect(
                    renderContext,
                    convertColour(getRegister(registerInts["vc"])),
                    getRegister(registerInts["vx0"]),
                    getRegister(registerInts["vy0"]),
                    getRegister(registerInts["vx1"]),
                    getRegister(registerInts["vy1"]),
                    )
                break;

            case interruptInts["vc"]:
                clearCanvas(renderContext, convertColour(getRegister(registerInts["vc"])))
                break;

            default:
                break;
        }

        stepVM();

    }

}