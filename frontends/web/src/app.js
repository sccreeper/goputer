import { glContext, canvas, currentInstructionHTML, programCounterHTML, peekRegHTML, peekRegInput } from "./init";
import globals from "./globals.js"
import { ShowError, ErrorTypes } from "./error";
import { drawSceneSimple } from "./gl/index.js";

var previousMousePos = {
    X: 0,
    Y: 0,
}
var currentMousePos = {
    X: 0,
    Y: 0,
}

//Other app logic

/**
 * 
 * @param {MouseEvent} e 
 * @returns {null}
 */
export function IOToggle(e) {

    if (!globals.vmIsAlive) {
        return
    }

    if (e.target.getAttribute("on") == "false") {
        e.target.setAttribute("on", "true")
    } else {
        e.target.setAttribute("on", "false")
    }

    globals.switchQueue.push(
        {
            register: e.target.getAttribute("reg"),
            enabled: (e.target.getAttribute("on") == "true") ? true : false,
        }
        
    )
}

export function PeekRegister() {
    if (peekRegInput.value == "" || !globals.vmInited) {
        return;
    } else {
        if (registerInts[peekRegInput.value] != undefined) {
            globals.registerPeekValue = peekRegInput.value;
            peekRegHTML.textContent = GetRegisterText(registerInts[globals.registerPeekValue]) 
            peekRegInput.setAttribute("valid-reg", "true");
        } else {
            peekRegInput.setAttribute("valid-reg", "false");
        }
    
    }
}

// Used in conjunction with PeekRegister
export function GetRegisterText(reg_int) {

    if (reg_int == registerInts["d0"]) {
        
        let data = getBuffer("data");
        let data_string = "";

        let last_val_zero = true;

        data.forEach(element => {

            if (element != 0) {
                last_val_zero = false;
                data_string += element.toString() + " ";
            } else if (element == 0 && !last_val_zero) {
                last_val_zero = true;
                data_string += ".. "
            }

        });

        return data_string;


    } else if (reg_int == registerInts["vt"]) {
        
        let t_buff = getBuffer("text");
        
        var t_codes = []

        t_buff.forEach(element => {
            
            if (element != 0) {
                t_codes.push(element)
            }

        });

        // Convert from array of ints to chars.
        return String.fromCharCode(...t_codes)

    } else {
        let hex_string = getRegister(registerInts[globals.registerPeekValue]).toString(16)
        hex_string = hex_string.split("")

        hex_string = hex_string.reverse()
        hex_string = hex_string.join("")

        return `0x${hex_string.toUpperCase().padStart(8, "0")} (${getRegister(registerInts[globals.registerPeekValue])})`;
    }
    
}

// Main app logic

export function Compile(e) {

    globals.errorDiv.replaceChildren();

    globals.compileFailed = false;
    compileCode(document.getElementById("code-textarea").value)
    globals.codeHasBeenCompiled = true;

    if(!globals.compileFailed) {

        document.getElementById("run-code-button").disabled = false;
        document.getElementById("download-code-button").disabled = false;

        ShowError(ErrorTypes.Success, "Code compiled successfully!");

    }

}

export function Run(e) { 

    if (!globals.codeHasBeenCompiled) {

        ShowError(ErrorTypes.Error, "No code has been compiled!")

    } else {

        initVM();

        globals.vmIsAlive = true;
        globals.runInterval = setInterval(Cycle, Math.round(1000 / globals.FPS));
        globals.vmInited = true;

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

export function handleKeyDown(e) {
    
    if (globals.vmIsAlive) {
        globals.keysDown.push(e.keyCode)
    }

}

export function handleKeyUp(e) {
    
    if (globals.vmIsAlive) {
        globals.keysUp.push(e.keyCode)
    }

}

//Performs one cycle of the VM & Updates UI
export function Cycle() {
    
    if (isFinished()) {
        
        clearInterval(globals.runInterval);
        canvas.setAttribute("running", "false");
        return;

    }

    if (!globals.vmIsAlive) {
        
        console.error("VM isn't alive therefore can't run code.");

    } else {

        //Handle called interrupts.

        var x = getInterrupt()

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
                        "on",
                        (getRegister(registerInts[globals.ioBulbNames[i]]) > 0) ? "true" : "false"
                    )

                }

            default:
                break;
        }

        // Video
        updateFramebuffer();
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
            
            setRegister(registerInts["mx"], previousMousePos.X);
            setRegister(registerInts["my"], previousMousePos.Y);

            previousMousePos.X = currentMousePos.X;
            previousMousePos.Y = currentMousePos.Y
        
            if (isSubscribed(interruptInts["mm"])) {

                sendInterrupt(interruptInts["mm"]);
            }

        }

        //Keyboard

        if (globals.keysDown.length > 0 ) {
            
            setRegister(registerInts["kc"], globals.keysDown.pop())

            if (isSubscribed(interruptInts["kd"])) {
                sendInterrupt(interruptInts["kd"])
            }

        }

        if (globals.keysUp.length > 0) {
            
            setRegister(registerInts["kp"], globals.keysUp.pop())

            if (isSubscribed(interruptInts["ku"])) {
                sendInterrupt(interruptInts["ku"])
            }

        }

        //IO Switches

        globals.switchQueue.forEach(element => {
        
            setRegister(registerInts[element.register], (element.enabled) ? 1 : 0)

            if (isSubscribed(interruptInts[element.register])) {
                sendInterrupt(interruptInts[element.register])
            }

        });

        globals.switchQueue = [];

        //Update hardware info

        currentInstructionHTML.innerHTML = String(currentItn());
        programCounterHTML.innerHTML = getRegister(registerInts["prc"])

        if (globals.registerPeekValue != null && GetRegisterText(registerInts[globals.registerPeekValue]) != globals.prevRegPeekValue) {

            globals.currentRegPeekValue = GetRegisterText(registerInts[globals.registerPeekValue])
            peekRegHTML.textContent = globals.currentRegPeekValue
            globals.prevRegPeekValue = globals.currentRegPeekValue

        }

        //Finally cycle VM & update graphics.

        cycleVM();

    }

}