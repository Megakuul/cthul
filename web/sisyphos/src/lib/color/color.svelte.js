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
export let Palette = $state(NewSandstormPalette());

export function NewSandstormPalette() {
  return {
    bgPrimary: () => { return "#F0E3CA" },
    bgSecondary: () => { return "#000000" },
    fgPrimary: () => { return "#1B1A17" },
    fgSecondary: () => { return "#A35709" },
    fgRune: () => { return "#5E4A11"},
    fgWave: () => { return "#0D65A4" },
    fgGranit: () => { return "#042F0B"},
    fgProton: () => { return "#57056C"},
    fgSuccess: () => { return "#255F38" },
    fgError: () => { return "#A31D1D" },
  }
}
