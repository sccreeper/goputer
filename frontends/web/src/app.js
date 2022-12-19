import { renderContext, canvas, currentInstructionHTML, programCounterHTML } from "./init";
import globals from "./globals.js"
import { clearCanvas, drawLine, drawRect, drawText, setPixel } from "./canvas_util";
import { ShowError, ErrorTypes } from "./error";

var previous_mouse_pos = {
    X: 0,
    Y: 0,
}
var current_mouse_pos = {
    X: 0,
    Y: 0,
}

//Other app logic

export function IOToggle(e) {

    if (!globals.vmIsAlive) {
        return
    }

    if (e.target.getAttribute("on") == "false") {
        e.target.setAttribute("on", "true")
    } else {
        e.target.setAttribute("on", "false")
    }

    globals.switch_queue.push(
        {
            register: e.target.getAttribute("reg"),
            enabled: (e.target.getAttribute("on") == "true") ? true : false,
        }
        
    )
}

// Main app logic

export function Compile() {

    globals.error_div.replaceChildren();

    globals.compile_failed = false;
    compileCode(document.getElementById("code-textarea").value)
    globals.codeHasBeenCompiled = true;

    if(!globals.compile_failed) {

        ShowError(ErrorTypes.Success, "Code compiled successfully!");

    }

}

export function Run() { 

    if (!globals.codeHasBeenCompiled) {

        ShowError(ErrorTypes.Error, "No code has been compiled!")

    } else {

        initVM();

        globals.vmIsAlive = true;
        globals.runInterval = setInterval(Cycle, Math.round(1000 / globals.FPS));

        canvas.setAttribute("running", "true");
        
    }

}

export function handleMouseMove(e) {

    if (globals.vmIsAlive) {
        if (globals.mouse_over_display) {
            current_mouse_pos.X = Math.round(e.clientX -  canvas.getBoundingClientRect().left);
            current_mouse_pos.Y = Math.round(e.clientY -  canvas.getBoundingClientRect().top);        
        }
    }

}

export function handleKeyDown(e) {
    
    if (globals.vmIsAlive) {
        globals.keys_down.push(e.keyCode)
    }

}

export function handleKeyUp(e) {
    
    if (globals.vmIsAlive) {
        globals.keys_up.push(e.keyCode)
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
                    globals.video_text = "";
                } else {
                    
                    var t1 = []

                    t.forEach(element => {
                        
                        if (element != 0) {
                            t1.push(element)
                        }

                    });

                    // Convert from array of ints to chars.
                    globals.video_text += String.fromCharCode(...t1)

                }


                drawText(
                    renderContext,
                    convertColour(getRegister(registerInts["vc"])),
                    getRegister(registerInts["vx0"]),
                    getRegister(registerInts["vy0"]),
                    globals.video_text
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
                globals.audio_volume.gain.value = 0;
                break;
            case interruptInts["sf"]:
                globals.oscillator.type = (getRegister(registerInts["sw"]) == 0) ? "square" : "sine";
                globals.oscillator.frequency.value = getRegister(registerInts["st"])
                globals.audio_volume.gain.value = getRegister(registerInts["sv"]) / 255;
                if (!globals.sound_started) {
                    globals.oscillator.start()
                    globals.sound_started = true;
                }

                break;

            default:
                break;
        }

        // Handle subscribed interrupts

        //Mouse

        if ((previous_mouse_pos.X != current_mouse_pos.X) || (previous_mouse_pos.Y != current_mouse_pos.Y)) {
            
            setRegister(registerInts["mx"], previous_mouse_pos.X);
            setRegister(registerInts["my"], previous_mouse_pos.Y);

            previous_mouse_pos.X = current_mouse_pos.X;
            previous_mouse_pos.Y = current_mouse_pos.Y
        
            if (isSubscribed(interruptInts["mm"])) {

                sendInterrupt(interruptInts["mm"]);
            }

        }

        //Keyboard

        if (globals.keys_down.length > 0 ) {
            
            setRegister(registerInts["kc"], globals.keys_down.pop())

            if (isSubscribed(interruptInts["kd"])) {
                sendInterrupt(interruptInts["kd"])
            }

        }

        if (globals.keys_up.length > 0) {
            
            setRegister(registerInts["kp"], globals.keys_up.pop())

            if (isSubscribed(interruptInts["ku"])) {
                sendInterrupt(interruptInts["ku"])
            }

        }

        //IO Switches

        globals.switch_queue.forEach(element => {
        
            setRegister(registerInts[element.register], (element.enabled) ? 1 : 0)

            if (isSubscribed(interruptInts[element.register])) {
                sendInterrupt(interruptInts[element.register])
            }

        });

        globals.switch_queue = [];

        //IO bulbs

        for (let i = 0; i < globals.io_bulb_names.length; i++) {
            
            globals.io_bulbs[globals.io_bulb_names[i]].setAttribute(
                "on",
                (getRegister(registerInts[globals.io_bulb_names[i]]) > 0) ? "true" : "false"
            )

        }

        //Update hardware info

        currentInstructionHTML.innerHTML = String(currentItn());
        programCounterHTML.innerHTML = getRegister(registerInts["prc"])

        //Finally cycle VM.

        stepVM();

    }

}