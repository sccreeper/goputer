/**
 * Clamp a number
 * @param {number} num 
 * @param {number} min 
 * @param {number} max 
 * @returns {number}
 */
export function clamp(num, min, max) {

    if (num > max) {
        return max
    } else if (num < min) {
        return min
    } else {
        return num
    }

}