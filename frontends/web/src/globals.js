import { ShowError } from "./error"

export default {

    codeHasBeenCompiled: false,
    vmIsAlive: false,
    vmInited: false,
    runInterval: null,
    FPS: 240,
    vmWorker: null,
    audioContext: null,
    oscillator: null,
    audioVolume: null,
    soundStarted: false,
    mouseOverDisplay: false,
    keysDown: [],
    keysUp: [],
    videoText: "",
    ioBulbs: {},
    ioBulbNames: ["io00", "io01","io02", "io03", "io04", "io05", "io06", "io07"],
    switchQueue: [],
    compileFailed: false,
    errorDiv: null,
    errorCount: 0,
    registerPeekValue: null,
    currentRegPeekValue: null,
    prevRegPeekValue: null,
    focusedFile: "main.gpasm",
    textureData: new Uint8Array(320*240*3)

}

window.showError = ShowError