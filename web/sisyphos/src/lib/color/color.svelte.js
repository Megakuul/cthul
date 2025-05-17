/** 
 * @typedef {Object} Palette
 * @property {() => string} bgPrimary
 * @property {() => string} fgPrimary
 * @property {() => string} bgCthul
 * @property {() => string} fgRune
 * @property {() => string} fgWave
 * @property {() => string} fgGranit
 * @property {() => string} fgProton
 * @property {() => string} fgSuccess
 * @property {() => string} fgError
 */

/** @type {Palette} */
let palette = $state(NewSandstormPalette());

/** @returns {Palette} */
export function Palette() {
  return palette
}

/** @param {Palette} newPalette */
export function SetPalette(newPalette) {
  palette = newPalette
}

export function NewSandstormPalette() {
  return {
    bgPrimary: () => { return "#E4D4C8" },
    fgPrimary: () => { return "#58391C" },
    bgCthul: () => { return "#000000" },
    fgRune: () => { return "#5E4A11"},
    fgWave: () => { return "#0D65A4" },
    fgGranit: () => { return "#042F0B"},
    fgProton: () => { return "#57056C"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}

export function NewSlatePalette() {
  return {
    bgPrimary: () => { return "#252A34" },
    fgPrimary: () => { return "#B7C4CF" },
    bgCthul: () => { return "#000000" },
    fgRune: () => { return "#5E4A11"},
    fgWave: () => { return "#0D65A4" },
    fgGranit: () => { return "#042F0B"},
    fgProton: () => { return "#57056C"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}
