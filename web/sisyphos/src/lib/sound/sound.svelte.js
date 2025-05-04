/** @type {AudioContext | undefined} */
let audioContext = undefined;

/** @type {GainNode | undefined} */
let audioGain = undefined;

async function Init() {
  if (!audioContext || audioContext.state === "closed") {
    audioContext = new AudioContext()
    audioGain = audioContext.createGain()
    audioGain.connect(audioContext.destination)
  } else if (audioContext.state === "suspended") {
    await audioContext.resume()
  }
}

/** @param {string} url @returns {Promise<void>} */
export async function Play(url) {
  if (!audioContext || !audioGain || audioContext.state !== "running") {
    Init()
  }
  
  if (!audioContext || !audioGain) {
    console.error("failed to initialize audio context")
    return
  }

  try {
    const response = await fetch(url)
    const buffer = await audioContext.decodeAudioData(await response.arrayBuffer())

    const source = audioContext.createBufferSource()
    source.buffer = buffer

    source.connect(audioGain)

    source.start(0)
  } catch (error) {
    console.error("failed to play sound: ", error)
  }
}

export async function Mute() {
  if (!audioContext || !audioGain || audioContext.state !== "running") {
    Init()
  }

  if (audioContext && audioGain) {
    audioGain.gain.setValueAtTime(0, audioContext.currentTime)
  }
}

export async function Unmute() {
  if (!audioContext || !audioGain || audioContext.state !== "running") {
    Init()
  }

  if (audioContext && audioGain) {
    audioGain.gain.setValueAtTime(1, audioContext.currentTime)
  }
}