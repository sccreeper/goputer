import { renderContext, canvas, currentInstructionHTML, programCounterHTML, peekRegHTML, peekRegInput } from "./init";
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

export function PeekRegister() {
    if (peekRegInput.value == "" || !globals.vmInited) {
        return;
    } else {
        if (registerInts[peekRegInput.value] != undefined) {
            globals.register_peek_value = peekRegInput.value;
            peekRegHTML.textContent = GetRegisterText(registerInts[globals.register_peek_value]) 
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
        let hex_string = getRegister(registerInts[globals.register_peek_value]).toString(16)
        hex_string = hex_string.split("")

        hex_string = hex_string.reverse()
        hex_string = hex_string.join("")

        return `0x${hex_string.toUpperCase().padStart(8, "0")} (${getRegister(registerInts[globals.register_peek_value])})`;
    }
    
}

// Main app logic

export function Compile(e) {

    globals.error_div.replaceChildren();

    globals.compile_failed = false;
    compileCode(document.getElementById("code-textarea").value)
    globals.codeHasBeenCompiled = true;

    if(!globals.compile_failed) {

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

            case interruptInts["iof"]:
                //Set IO states for IO bulbs.

                for (let i = 0; i < globals.io_bulb_names.length; i++) {
                    
                    globals.io_bulbs[globals.io_bulb_names[i]].setAttribute(
                        "on",
                        (getRegister(registerInts[globals.io_bulb_names[i]]) > 0) ? "true" : "false"
                    )

                }

            default:
                break;
        }

        // Video brightness

        let col = "";

        // Avoid divide by zero error.
        if (getRegister(registerInts["vb"]) == 0) {
            col = "rgba(0, 0, 0, 1)";
        } else {
            col = `rgba(0, 0, 0, ${1 - Math.pow((Math.pow(getRegister(registerInts["vb"]), -1)) * 255, -1)})`;
        }

        console.log(1 - Math.pow((Math.pow(getRegister(registerInts["vb"]), -1)) * 255, -1))

        drawRect(renderContext, col, 0, 0, 640, 480);

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

        //Update hardware info

        currentInstructionHTML.innerHTML = String(currentItn());
        programCounterHTML.innerHTML = getRegister(registerInts["prc"])

        if (globals.register_peek_value != null && GetRegisterText(registerInts[globals.register_peek_value]) != globals.prev_reg_peek_value) {

            globals.current_reg_peek_value = GetRegisterText(registerInts[globals.register_peek_value])
            peekRegHTML.textContent = globals.current_reg_peek_value
            globals.prev_reg_peek_value = globals.current_reg_peek_value

        }

        //Finally cycle VM.

        stepVM();

    }

}