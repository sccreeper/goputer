import { db, fileTableName } from "./db.js"
import * as Comlink from "comlink";
import("../static/wasm_exec.js");

/**
 * @typedef {"r00" | "r01" | "r02" | "r03" | "r04" | "r05" | "r06" | "r07" |
 *           "r08" | "r09" | "r10" | "r11" | "r12" | "r13" | "r14" | "r15" |
 *           "vx0" | "vy0" | "vx1" | "vy1" |
 *           "vc" | "vb" | "vt" |
 *           "kc" | "kp" |
 *           "mx" | "my" | "mb" |
 *           "st" | "sv" |
 *           "a0" | "d0" |
 *           "stk" | "stz" |
 *           "io00" | "io01" | "io02" | "io03" | "io04" | "io05" | "io06" | "io07" |
 *           "io08" | "io09" | "io10" | "io11" | "io12" | "io13" | "io14" | "io15" |
 *           "prc" |
 *           "cstk" | "cstz" |
 *           "dl" | "dp" |
 *           "sw"} Register
 * 
 */


self.showError = null;

const goputer = {
    /** @returns {Promise<void>} */
    compileCode() {return compileCode()},
    /** @returns {Promise<void>} */
    initVm() {return initVM()},

    async workerInit() {

        //Init Go WASM before anything else
        const go = new Go();
        await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });

    },

    /**
     * Set register value
     * @param {number} register 
     * @param {number} value 
     * @returns {Promise<void>}
     */
    setRegister(register, value) { return setRegister(register, value) },
    
    /**
     * Get a register value
     * @param {number} register 
     * @returns {Promise<number>}
     */
    getRegister(register) { return getRegister(register) },

    /**
     * Length of dest should be 4 or 128
     * @param {number} register 
     * @param {Uint8Array} dest
     * @returns {Promise<void>}
     */
    getRegisterBytes(register, dest) { return getRegisterBytes(register, dest) },

    /**
     * Gets either video or text buffer
     * @param {"text"|"data"} bufferName 
     * @param {Uint8Array} dest 
     */
    getBuffer(bufferName, dest) { return getBuffer(bufferName, dest) },

    /**
     * Pops interrupt from queue
     * @returns {Promise<number|null>} 
     */
    getInterrupt() {
        let x = getInterrupt()
        return x
    
    },

    /**
     * 
     * @param {Promise<number>} interrupt 
     */
    sendInterrupt(interrupt) { return sendInterrupt(interrupt) },
    
    /**
     * Used for checking wether or not to add interrupt to queue using `sendInterrupt`.
     * @param {number} interrupt 
     * @returns {Promise<boolean>}
     */
    isSubscribed(interrupt) { return isSubscribed(interrupt) },

    /** @type {string} */
    get currentItn() {
        return currentItn().toString()
    },

    cycleVm() {
        cycleVM()
    },

    async run() {
        while (!isFinished()) {
            this.cycleVm()
            await new Promise(resolve => setTimeout(resolve, 0))
        }
    },

    /** @type {Promise<boolean>} */
    isFinished() {
        let x = isFinished()

        return x
    },

    files: {
        /**
         * Overwrite the contents of a specific file. Key doesn't have to be preexisting and as such this method can be used to create new files.
         * @param {string} key typically the name of the file
         * @param {Uint8Array} data 
         * @param {number} size length of the data
         * @param {import("./editor/code_tab.js").FileType} type 1 of 3 specified types
         * @param {boolean} [isNew=false] defaults to false
         * @param {boolean} [writeToDb=true] defaults to true
         */
        update(key, data, size, type, isNew = false, writeToDb = true) {

            if (writeToDb) {
                if (isNew) {
                    db.table(fileTableName).put(
                        {
                            name: key,
                            data: data,
                            type: type,
                        }
                    )
                } else {
                    db.table(fileTableName).update(
                        key,
                        {
                            name: key,
                            data: data,
                            type: type,
                        },
                    )
                }   
            }

            updateFile(key, data, size, type, isNew)
        },

        /**
         * 
         * @param {string} key 
         * @param {string} newKey 
         */
        rename(key, newKey) {
            let fileSize = this.size(key)
            let fileType = this.type(key)
            let fileData = new Uint8Array(fileSize)
            
            this.get(key, fileData)
            this.remove(key)
    
            if (fileType == "image") {
                let imageMapData = imageMap.get(key)
                imageMap.delete(key)
                imageMap.set(newKey, imageMapData)   
            }
            
            /** @type {import("dexie").Table} */
            (db.files).put({
                name: newKey,
                data: fileData,
                type: fileType
            })
       
            updateFile(newKey, fileData, fileSize, fileType, false)
        },
        
        /**
         * Delete file based on key. Panics if no such file exists.
         * @param {string} key 
         */
        remove(key) {
            db.table(fileTableName).delete(key)

            removeFile(key)
        },

        /**
         * Return file data based on key. Panics if no such file exists.
         * @param {string} key
         * @param {Uint8Array} dest
         */
        get(key, dest) {
            getFile(key, dest)
        },

        /**
         * Get the size of a file in bytes. Useful for allocating `Uint8Array`.
         * @param {string} key 
         * @returns {number}
         */
        size(key) {
            return getFileSize(key)
        },

        /**
         * Does this file exist in the map
         * @param {string} key
         * @returns {boolean} 
         */
        exists(key) {
            return doesFileExist(key)
        },

        /**
         * Get the type of a specific file, image, binary, or text.
         * @param {string} key 
         * @returns {import("./editor/code_tab.js").FileType}
         */
        type(key) {
            return getFileType(key)
        },

        /** @type {number} */
        get numFiles() {
            return numFiles()
        },

        /** @type {string[]} */
        get fileNames() {
            return getFiles()
        },
    },

    /**
     * Returns the current program. Can be empty if no program has been set.
     * @param {Uint8Array} dest  
     */
    getProgramBytes(dest) {
        return getProgramBytes(dest)
    },

    /**
     * 
     * @returns {number}
     */
    getProgramLength() {
        return getProgramLength()
    },

    /**
     * 
     * @param {Uint8Array} data 
     * @param {number} size 
     */
    setProgramBytes(data, size) {
        return setProgramBytes(data, size)
    },

    /**
     * Returns a Uint8Array of the framebuffer.
     * 
     * Does some magic with a worker scope variable called `wasmFramebuffer`.
     * @param {Uint8Array<SharedArrayBuffer>} sharedBuffer
     * @returns {Promise<number>}
     */
    updateFramebuffer(sharedBuffer) {
        return updateFramebuffer(sharedBuffer)
    },

    /**
     * Maps a JS event.code to goputer's internal keycodes
     * @param {string} key
     * @returns {Promise<number>} 
     */
    mappedKey(key) {
        return getMappedKey(key)
    },

    /**
     * @param {any[]} bin
     * @returns {any}
     */
    disassembleCode(bin) {
        return JSON.parse(disassembleCode(bin))
    },

    // Constants

    /** @type {Object} */
    get interruptInts() {
        return interruptInts
    },
    /** @type {Object} */
    get instructionInts() {
        return instructionInts
    },
    /** @type {Object} */
    get registerInts() {
        return registerInts
    },

    /** @type {string[]} */
    get instructionArray() {
        return instructionArray
    },

    /** @type {number} */
    get memOffset() {
        return memOffset
    },

    util: {
        /**
         * Converts a uint32 (RGBA8) colour to a rgba(r, g, b, a) string.
         * @param {number} colour
         * @returns {string} 
         */
        convertColour(colour) {
            return convertColour(colour)
        },

        /**
         * Convert a number (typically an address) to a hexadecimal string.
         * @param {number} num 
         * @param {boolean} useOffset 
         */
        convertHex(num, isOffset) {
            return convertHex(num, isOffset)
        }
    },

    /** @type {number} */
    get usableMemorySize() {
        return memSize
    },

    set onerror(val) {
        self.showError = val
    },
    
    get onerror() {
        return self.showError 
    }

}

Comlink.expose(goputer, self);

/**
 * @typedef {typeof goputer} Goputer
 */