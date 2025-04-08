/** 
 * @typedef {Object} Palette
 * @property {() => string} bgPrimary
 * @property {() => string} bgSecondary
 * @property {() => string} fgPrimary
 * @property {() => string} fgRune
 * @property {() => string} fgWave
 * @property {() => string} fgGranit
 * @property {() => string} fgProton
 * @property {() => string} fgSuccess
 * @property {() => string} fgError
 */

/** @type {Palette} */
let palette = $state(NewDefaultPalette());

/** @returns {Palette} */
export function Palette() {
  return palette
}

/** @param {Palette} palette */
export function SetPalette(newPalette) {
  palette = newPalette
}

export function NewDefaultPalette() {
  return {
    bgPrimary: () => { return "#1e2124" },
    bgSecondary: () => { return "#000000" },
    fgPrimary: () => { return "#F5F7F8" },
    fgSecondary: () => { return "#F9F7F7" },
    fgRune: () => { return "#5E4A11"},
    fgWave: () => { return "#0D65A4" },
    fgGranit: () => { return "#042F0B"},
    fgProton: () => { return "#57056C"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}

export function NewRunePalette() {
  return {
    bgPrimary: () => { return "#1e2124" },
    bgSecondary: () => { return "#000000" },
    fgPrimary: () => { return "#5E4A11" },
    fgSecondary: () => { return "#826411" },
    fgRune: () => { return "#5E4A11"},
    fgWave: () => { return "#5E4A11" },
    fgGranit: () => { return "#5E4A11"},
    fgProton: () => { return "#5E4A11"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}

export function NewWavePalette() {
  return {
    bgPrimary: () => { return "#1e2124" },
    bgSecondary: () => { return "#000000" },
    fgPrimary: () => { return "#0D65A4" },
    fgSecondary: () => { return "#0c7ac9" },
    fgRune: () => { return "#0D65A4"},
    fgWave: () => { return "#0D65A4" },
    fgGranit: () => { return "#0D65A4"},
    fgProton: () => { return "#0D65A4"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}

export function NewGranitPalette() {
  return {
    bgPrimary: () => { return "#1e2124" },
    bgSecondary: () => { return "#000000" },
    fgPrimary: () => { return "#042F0B" },
    fgSecondary: () => { return "#045411" },
    fgRune: () => { return "#042F0B"},
    fgWave: () => { return "#042F0B" },
    fgGranit: () => { return "#042F0B"},
    fgProton: () => { return "#042F0B"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}

export function NewProtonPalette() {
  return {
    bgPrimary: () => { return "#1e2124" },
    bgSecondary: () => { return "#000000" },
    fgPrimary: () => { return "#57056C" },
    fgSecondary: () => { return "#6f068a" },
    fgRune: () => { return "#57056C"},
    fgWave: () => { return "#57056C" },
    fgGranit: () => { return "#57056C"},
    fgProton: () => { return "#57056C"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}
