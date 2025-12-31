import { canvas, currentInstructionHTML, programCounterHTML, peekRegHTML, peekRegInput } from "./init";
import globals from "./globals.js"
import { ShowError, ErrorTypes } from "./error";
import { goputer, registerInts, interruptInts } from "./goputer.js";
import { checkVisible } from "./util.js";
import * as Comlink from "comlink";

var previousMousePos = {
    X: 0,
    Y: 0,
}
var currentMousePos = {
    X: 0,
    Y: 0,
}

var executionStartTime = 0
var uiUpdatesCompleted = 0

export async function PeekRegister() {
    if (peekRegInput.value == "" || !globals.vmInited) {
        return;
    } else {
        if (registerInts[peekRegInput.value] != undefined) {
            
            globals.registerPeekValue = peekRegInput.value;
            peekRegHTML.textContent = await GetRegisterText(
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
 * @param {"hex"|"binary"|"text"|"decimal"} format 
 * @returns {string}
 */
export async function GetRegisterText(regInt, format) {
    

    let bytes = new Uint8Array(new SharedArrayBuffer(
        regInt == registerInts["d0"] || regInt == registerInts["vt"] ? 128 : 4
    ));

    if (regInt == registerInts["d0"] || regInt == registerInts["vt"]) {
        
        if (regInt == registerInts["d0"]) {
            await goputer.getBuffer("data", Comlink.transfer(bytes, bytes.data))
        } else {
            await goputer.getBuffer("text", Comlink.transfer(bytes, bytes.data))
        }

    } else {

        await goputer.getRegisterBytes(regInt, Comlink.transfer(bytes, bytes.data))
    
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

                return (await goputer.getRegister(regInt)).toString()

            }

        default:
            break;
    }
    
}

// Main app logic

export async function Compile(e) {

    globals.errorDiv.replaceChildren();

    globals.compileFailed = false;
    await goputer.compileCode()
    globals.codeHasBeenCompiled = true;

    if(!globals.compileFailed) {

        document.getElementById("run-code-button").disabled = false;
        document.getElementById("download-code-button").disabled = false;

        ShowError(ErrorTypes.Success, "Code compiled successfully!");

    } else {
        globals.codeHasBeenCompiled = false;
    }

}

let lastUpdateTime = 0;
const targetFrameTime = Math.round(1000 / globals.UPS);

function scheduleNextUpdate() {
    if (!globals.vmIsAlive) return;

    requestAnimationFrame(async (timestamp) => {
        const elapsed = timestamp - lastUpdateTime;

        if (elapsed >= targetFrameTime) {
            lastUpdateTime = timestamp;
            await UiUpdate();
        }

        scheduleNextUpdate();
    })
}

export async function Run(e) { 

    if (!globals.codeHasBeenCompiled) {

        ShowError(ErrorTypes.Error, "No code has been uploaded or compiled!")

    } else {

        await goputer.initVm();

        globals.vmIsAlive = true;
        
        SetKeyboardLocking(true)
    
        globals.vmInited = true;

        goputer.run();

        executionStartTime = Date.now()

        canvas.setAttribute("running", "true");

        lastUpdateTime = performance.now()
        scheduleNextUpdate()
        
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
export async function handleKeyDown(e) {

    console.log("help")
    
    if (globals.keyboardLocked) {
        e.preventDefault()
        globals.keysDown.push(await goputer.mappedKey(e.code))
    }

}

/**
 * 
 * @param {KeyboardEvent} e 
 */
export async function handleKeyUp(e) {
    
    if (globals.keyboardLocked) {
        e.preventDefault()
        globals.keysUp.push(await goputer.mappedKey(e.code)) // I am aware this is depreceated however, this is the most practical way to get integer keycodes.
    }

}

/**
 * 
 * @param {boolean} locked 
 */
export function SetKeyboardLocking(locked) {

    if (locked) {
        globals.keyboardLocked = true

        document.getElementById("kbd-locked-message").querySelector("span").innerText = " Keyboard locked"
        document.getElementById("kbd-locked-message").querySelector("i").classList.remove("bi-unlock")
        document.getElementById("kbd-locked-message").querySelector("i").classList.add("bi-lock")
    } else {
        globals.keyboardLocked = false

        document.getElementById("kbd-locked-message").querySelector("span").innerText = " Keyboard unlocked"
        document.getElementById("kbd-locked-message").querySelector("i").classList.add("bi-unlock")
        document.getElementById("kbd-locked-message").querySelector("i").classList.remove("bi-lock")
    }
}

export async function UiUpdate() {

    if (!globals.vmIsAlive) {
        
        console.error("VM isn't alive therefore can't run code.");

    } else {

        //Handle called interrupts.

        var x = await goputer.getInterrupt()

        switch (x) {
            case interruptInts["ss"]:
                globals.oscillator.frequency.value = 0;
                globals.audioVolume.gain.value = 0;
                break;
            case interruptInts["sf"]:
                globals.oscillator.type = (await goputer.getRegister(registerInts["sw"]) == 0) ? "square" : "sine";
                globals.oscillator.frequency.value = await goputer.getRegister(registerInts["st"])
                globals.audioVolume.gain.value = await goputer.getRegister(registerInts["sv"]) / 255;
                if (!globals.soundStarted) {
                    globals.oscillator.start()
                    globals.soundStarted = true;
                }

                break;

            case interruptInts["iof"]:
                //Set IO states for IO bulbs.

                for (let i = 0; i < globals.ioBulbNames.length; i++) {

                    console.log(`${i}: ${(await goputer.getRegister(registerInts[globals.ioBulbNames[i]]) > 0)}`)
                    
                    globals.ioBulbs[globals.ioBulbNames[i]].setAttribute(
                        "enabled",
                        (await goputer.getRegister(registerInts[globals.ioBulbNames[i]]) > 0) ? "true" : "false"
                    )

                }
                break;
            default:
                break;
        }

        await goputer.drawing.drawScene();

        // Handle subscribed interrupts

        //Mouse

        if ((previousMousePos.X != currentMousePos.X) || (previousMousePos.Y != currentMousePos.Y)) {
            
            await goputer.setRegister(registerInts["mx"], Math.floor(previousMousePos.X / 2));
            await goputer.setRegister(registerInts["my"], Math.floor(previousMousePos.Y / 2));

            previousMousePos.X = currentMousePos.X;
            previousMousePos.Y = currentMousePos.Y
        
            if (await goputer.isSubscribed(interruptInts["mm"])) {

                await goputer.sendInterrupt(interruptInts["mm"]);
            }

        }

        //Keyboard

        while (globals.keysDown.length > 0) {
            
            await goputer.setRegister(registerInts["kc"], globals.keysDown.pop())
            await goputer.setRegister(registerInts["kp"], 1)

            if (await goputer.isSubscribed(interruptInts["kd"])) {
                await goputer.sendInterrupt(interruptInts["kd"])
            }

        }

        while (globals.keysUp.length > 0) {
            
            await goputer.setRegister(registerInts["kc"], globals.keysUp.pop())
            await goputer.setRegister(registerInts["kp"], 0)

            if (await goputer.isSubscribed(interruptInts["ku"])) {
                await goputer.sendInterrupt(interruptInts["ku"])
            }

        }

        //IO Switches

        const switchQueueCopy = [...globals.switchQueue]

        for (const element of switchQueueCopy) {
        
            await goputer.setRegister(registerInts[element.register], (element.enabled) ? 1 : 0)
            console.log(element.register)

            if (await goputer.isSubscribed(interruptInts[element.register])) {
                console.log("subbed")
                await goputer.sendInterrupt(interruptInts[element.register])
            }

        };

        globals.switchQueue = [];

        //Update hardware info

        if (checkVisible(currentInstructionHTML)) {
            currentInstructionHTML.textContent = String(await goputer.currentItn);   
        }

        if (checkVisible(programCounterHTML)) {
            programCounterHTML.textContent = await GetRegisterText(registerInts["prc"], "hex")   
        }

        if (globals.registerPeekValue != null && await GetRegisterText(registerInts[globals.registerPeekValue], document.getElementById("peek-format-select").value) != globals.prevRegPeekValue) {

            globals.currentRegPeekValue = await GetRegisterText(registerInts[globals.registerPeekValue], document.getElementById("peek-format-select").value)
            peekRegHTML.textContent = globals.currentRegPeekValue
            globals.prevRegPeekValue = globals.currentRegPeekValue

        }

        //Finally cycle VM & update graphics.
        uiUpdatesCompleted++

    }

    if ((await goputer.isFinished()) && uiUpdatesCompleted != 0) {

        console.log(`Time elapsed: ${Date.now()-executionStartTime}ms`)
        console.log(`Average time per UI update: ${(Date.now()-executionStartTime)/uiUpdatesCompleted}ms`)

        executionStartTime = 0
        uiUpdatesCompleted = 0
        
        globals.vmIsAlive = false;
        canvas.setAttribute("running", "false");
        SetKeyboardLocking(false);
        return;

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