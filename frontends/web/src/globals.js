import { ShowError } from "./error"
export default {

    codeHasBeenCompiled: false,
    vmIsAlive: false,
    keyboardLocked: false,
    vmInited: false,
    runInterval: null,
    FPS: 512,
    /**
     * @type {AudioContext}
     */
    audioContext: null,
    /**
     * @type {OscillatorNode}
     */
    oscillator: null,
    /**
     * @type {GainNode}
     */
    audioVolume: null,
    /**
     * @type {MediaStreamAudioDestinationNode}
     */
    audioMediaStreamDestination: null,
    soundStarted: false,
    mouseOverDisplay: false,
    keysDown: [],
    keysUp: [],
    videoText: "",
    ioBulbs: {},
    ioBulbNames: ["io00", "io01","io02", "io03", "io04", "io05", "io06", "io07"],
    /** @type {{register: string, enabled: boolean}[]} */
    switchQueue: [],
    compileFailed: false,
    errorDiv: null,
    errorCount: 0,
    registerPeekValue: null,
    currentRegPeekValue: null,
    prevRegPeekValue: null,
    focusedFile: "main.gpasm",
}

window.textureData = new Uint8Array(320*240*3)
window.showError = ShowError