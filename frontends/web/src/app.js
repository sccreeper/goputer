import { glContext, canvas, currentInstructionHTML, programCounterHTML, peekRegHTML, peekRegInput } from "./init";
import globals from "./globals.js"
import { ShowError, ErrorTypes } from "./error";
import { drawSceneSimple } from "./gl/index.js";
import { goputer } from "./goputer.js";

var previousMousePos = {
    X: 0,
    Y: 0,
}
var currentMousePos = {
    X: 0,
    Y: 0,
}

var executionStartTime = 0
var cyclesCompleted = 0

export function PeekRegister() {
    if (peekRegInput.value == "" || !globals.vmInited) {
        return;
    } else {
        if (registerInts[peekRegInput.value] != undefined) {
            
            globals.registerPeekValue = peekRegInput.value;
            peekRegHTML.textContent = GetRegisterText(
                registerInts[globals.registerPeekValue], 
                document.getElementById("peek-format-select").value
            ) 
            
            peekRegInput.setAttribute("valid-reg", "true");
        
        } else {
            peekRegInput.setAttribute("valid-reg", "false");
        }
    
    }
}

/**
 * Used in conjunction with PeekRegister
 * @param {number} regInt 
 * @param {string} format 
 * @returns {string}
 */
export function GetRegisterText(regInt, format) {
    

    let bytes = new Uint8Array(
        regInt == registerInts["d0"] || regInt == registerInts["vt"] ? 128 : 4
    );

    if (regInt == registerInts["d0"] || regInt == registerInts["vt"]) {
        
        if (regInt == registerInts["d0"]) {
            goputer.getBuffer("data", bytes)
        } else {
            goputer.getBuffer("text", bytes)
        }

    } else {

        goputer.getRegisterBytes(registerInts[globals.registerPeekValue], bytes)
    
    }

    let result = ""

    switch (format) {
        case "hex":

            bytes.forEach(element => {
                result += element.toString(16).padStart(2, "0")
            });

            return `0x${result}`

        
        case "binary":

            bytes.forEach(element => {
                result += element.toString(2).padStart(8, "0")
                result += " "
            });

            return result
            
        case "text":
            
            for (let i = 0; i < bytes.length; i++) {
                
                if (bytes[i] == 0) {
                    continue
                } else {
                    result += String.fromCharCode(bytes[i])
                }
                
            }

            if (result.length == 0) {
                return "No string found"
            }

            return result

        case "decimal":

            if (regInt == registerInts["d0"] || regInt == registerInts["vt"]) {
                
                bytes.forEach(element => {
                    result += element.toString() + " "
                });

                return result

            } else {

                return getRegister(regInt).toString()

            }

        default:
            break;
    }
    
}

// Main app logic

export function Compile(e) {

    globals.errorDiv.replaceChildren();

    globals.compileFailed = false;
    goputer.compileCode()
    globals.codeHasBeenCompiled = true;

    if(!globals.compileFailed) {

        document.getElementById("run-code-button").disabled = false;
        document.getElementById("download-code-button").disabled = false;

        ShowError(ErrorTypes.Success, "Code compiled successfully!");

    }

}

export function Run(e) { 

    if (!globals.codeHasBeenCompiled) {

        ShowError(ErrorTypes.Error, "No code has been uploaded or compiled!")

    } else {

        goputer.initVm();

        globals.vmIsAlive = true;
        globals.runInterval = setInterval(Cycle, Math.round(1000 / globals.FPS));
        globals.vmInited = true;

        executionStartTime = Date.now()

        canvas.setAttribute("running", "true");
        
    }

}

export function handleMouseMove(e) {

    if (globals.vmIsAlive) {
        if (globals.mouseOverDisplay) {
            currentMousePos.X = Math.round(e.clientX -  canvas.getBoundingClientRect().left);
            currentMousePos.Y = Math.round(e.clientY -  canvas.getBoundingClientRect().top);        
        }
    }

}

/**
 * 
 * @param {KeyboardEvent} e 
 */
export function handleKeyDown(e) {
    
    if (globals.vmIsAlive) {
        e.preventDefault()
        globals.keysDown.push(e.keyCode)
    }

}

/**
 * 
 * @param {KeyboardEvent} e 
 */
export function handleKeyUp(e) {
    
    if (globals.vmIsAlive) {
        e.preventDefault()
        globals.keysUp.push(e.keyCode) // I am aware this is depreceated however, this is the most practical way to get integer keycodes.
    }

}

//Performs one cycle of the VM & Updates UI
export function Cycle() {
    
    if (goputer.isFinished) {

        console.log(`Time elapsed: ${Date.now()-executionStartTime}ms`)
        console.log(`Average time per cycle: ${(Date.now()-executionStartTime)/cyclesCompleted}ms`)

        executionStartTime = 0
        cyclesCompleted = 0
        
        clearInterval(globals.runInterval);
        canvas.setAttribute("running", "false");
        return;

    }

    if (!globals.vmIsAlive) {
        
        console.error("VM isn't alive therefore can't run code.");

    } else {

        //Handle called interrupts.

        var x = goputer.getInterrupt()

        switch (x) {
            case interruptInts["ss"]:
                globals.oscillator.frequency.value = 0;
                globals.audioVolume.gain.value = 0;
                break;
            case interruptInts["sf"]:
                globals.oscillator.type = (getRegister(registerInts["sw"]) == 0) ? "square" : "sine";
                globals.oscillator.frequency.value = getRegister(registerInts["st"])
                globals.audioVolume.gain.value = getRegister(registerInts["sv"]) / 255;
                if (!globals.soundStarted) {
                    globals.oscillator.start()
                    globals.soundStarted = true;
                }

                break;

            case interruptInts["iof"]:
                //Set IO states for IO bulbs.

                for (let i = 0; i < globals.ioBulbNames.length; i++) {
                    
                    globals.ioBulbs[globals.ioBulbNames[i]].setAttribute(
                        "enabled",
                        (getRegister(registerInts[globals.ioBulbNames[i]]) > 0) ? "true" : "false"
                    )

                }

            default:
                break;
        }

        // Video
        goputer.updateFramebuffer();
        drawSceneSimple(glContext)

        // Video brightness

        /**
         * @type {number}
         */
        let col = 0.0;

        // Avoid divide by zero error.
        if (getRegister(registerInts["vb"]) == 0) {
            col = 1.0;
        } else {
            col = 1 - Math.pow((Math.pow(getRegister(registerInts["vb"]), -1)) * 255, -1);
        }

        glContext.clearColor(0.0, 0.0, 0.0, col);

        // Handle subscribed interrupts

        //Mouse

        if ((previousMousePos.X != currentMousePos.X) || (previousMousePos.Y != currentMousePos.Y)) {
            
            goputer.setRegister(registerInts["mx"], Math.floor(previousMousePos.X / 2));
            goputer.setRegister(registerInts["my"], Math.floor(previousMousePos.Y / 2));

            previousMousePos.X = currentMousePos.X;
            previousMousePos.Y = currentMousePos.Y
        
            if (goputer.isSubscribed(interruptInts["mm"])) {

                goputer.sendInterrupt(interruptInts["mm"]);
            }

        }

        //Keyboard

        if (globals.keysDown.length > 0 ) {
            
            goputer.setRegister(registerInts["kc"], globals.keysDown.pop())

            if (goputer.isSubscribed(interruptInts["kd"])) {
                goputer.sendInterrupt(interruptInts["kd"])
            }

        }

        if (globals.keysUp.length > 0) {
            
            goputer.setRegister(registerInts["kp"], globals.keysUp.pop())

            if (goputer.isSubscribed(interruptInts["ku"])) {
                goputer.sendInterrupt(interruptInts["ku"])
            }

        }

        //IO Switches

        globals.switchQueue.forEach(element => {

            console.log(element)
        
            goputer.setRegister(registerInts[element.register], (element.enabled) ? 1 : 0)

            if (goputer.isSubscribed(interruptInts[element.register])) {
                goputer.sendInterrupt(interruptInts[element.register])
            }

        });

        globals.switchQueue = [];

        //Update hardware info

        currentInstructionHTML.innerHTML = String(goputer.currentItn);
        programCounterHTML.innerHTML = getRegister(registerInts["prc"])

        if (globals.registerPeekValue != null && GetRegisterText(registerInts[globals.registerPeekValue], document.getElementById("peek-format-select").value) != globals.prevRegPeekValue) {

            console.log("updating register...")

            globals.currentRegPeekValue = GetRegisterText(registerInts[globals.registerPeekValue], document.getElementById("peek-format-select").value)
            peekRegHTML.textContent = globals.currentRegPeekValue
            globals.prevRegPeekValue = globals.currentRegPeekValue

        }

        //Finally cycle VM & update graphics.

        goputer.cycleVm();
        cyclesCompleted++

    }

}

/**
 * Save current canvas state as a PNG
 * @param {MouseEvent} e 
 */
export function SaveVideo(e) {

    let downloadLink = document.createElement("a")
    downloadLink.href = canvas.toDataURL("image/png") 
    downloadLink.download = "video.png"
    downloadLink.click()

}