/** @type {boolean} */
let gimmick = $state(true);

/** @returns {boolean} */
export function Gimmick() {
  return gimmick
}

/** @param {boolean} enabled */
export function SetGimmick(enabled) {
  gimmick = enabled
}