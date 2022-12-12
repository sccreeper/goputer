import { renderContext, canvas } from "./init";
import global from "./globals.js";
import { clearCanvas, drawLine, drawRect, drawText, setPixel } from "./canvas_util";
import globals from "./globals.js";

var video_text = ""

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
            case interruptInts["vl"]:
                drawLine(
                    renderContext,
                    convertColour(getRegister(registerInts["vc"])),
                    getRegister(registerInts["vx0"]),
                    getRegister(registerInts["vy0"]),
                    getRegister(registerInts["vx1"]),
                    getRegister(registerInts["vy1"])
                )
            case interruptInts["vt"]:
                
                var t = getBuffer("text")

                if (t[0] == 0) {
                    video_text = "";
                } else {
                    
                    var t1 = []

                    t.forEach(element => {
                        
                        if (element != 0) {
                            t1.push(element)
                        }

                    });

                    // Convert from array of ints to chars.
                    video_text += String.fromCharCode(...t1)

                }


                drawText(
                    renderContext,
                    convertColour(getRegister(registerInts["vc"])),
                    getRegister(registerInts["vx0"]),
                    getRegister(registerInts["vy0"]),
                    video_text
                )
            case interruptInts["vp"]:
                setPixel(
                    renderContext,
                    convertColour(getRegister(registerInts["vc"])),
                    getRegister(registerInts["vx0"]),
                    getRegister(registerInts["vy0"])
                )
            case interruptInts["ss"]:
                globals.oscillator.frequency.value = 0;
                globals.audioVolume.gain.value = 0;
                break;
            case interruptInts["sf"]:
                globals.oscillator.type = (getRegister(registerInts["sw"]) == 0) ? "square" : "sine";
                globals.oscillator.frequency.value = getRegister(registerInts["st"])
                globals.audioVolume.gain.value = getRegister(registerInts["sv"]) / 255;
                if (!globals.sound_started) {
                    globals.oscillator.start()
                    globals.sound_started = true;
                }

                break;

            default:
                break;
        }

        stepVM();

    }

}