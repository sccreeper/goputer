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

/**
 * 
 * @param {HTMLElement} elm 
 * @returns {boolean}
 */
export function checkVisible(elm) {
  var rect = elm.getBoundingClientRect();
  var viewHeight = Math.max(document.documentElement.clientHeight, window.innerHeight);
  return !(rect.bottom < 0 || rect.top - viewHeight >= 0);
}