export default {

    codeHasBeenCompiled: false,
    vmIsAlive: false,
    runInterval: null,
    FPS: 60,
    vmWorker: null,
    audio_context: null,
    oscillator: null,
    audio_volume: null,
    sound_started: false,
    mouse_over_display: false,
    keys_down: [],
    keys_up: [],
    video_text: "",
    io_bulbs: {},
    io_bulb_names: ["io00", "io01","io02", "io03", "io04", "io05", "io06", "io07"],
    switch_queue: [],

}