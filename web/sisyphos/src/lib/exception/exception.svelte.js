/** 
 * @typedef {Object} Exception
 * @property {string} title
 * @property {string} desc
 */

/** @type {Exception | undefined} */
let exception = $state(undefined);

/** @type {number | undefined} */
let exceptionTimeout = undefined;

/** @returns {Exception | undefined} */
export function Exception() {
  return exception
}

/** @param {Exception | undefined} newException */
export function SetException(newException) {
  exception = newException
}