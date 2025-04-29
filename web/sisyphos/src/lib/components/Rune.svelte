<script>
  import rune from "$lib/assets/rune.svg";

  let { selected = $bindable() } = $props();
</script>

<button
  aria-label="rune"
  class="w-full bg-slate-100/20 rounded-3xl backdrop-blur-3xl cursor-pointer transition-all duration-700 overflow-hidden"
  class:selected={selected === "rune"}
  onclick={() => {
    if (selected === "rune") {
      selected = undefined;
    } else {
      selected = "rune";
    }
  }}
>
  <img alt="rune" src={rune} />
</button>

<style>
  button {
    border-width: 6px;
    border-color: #5e4a11;
    box-shadow: rgba(255, 255, 255, 0.06) 0px 2px 4px 0px inset;
  }

  button.selected {
    filter: drop-shadow(6px 6px 10px #5e4a11);
  }

  button:before {
    content: "";
    position: absolute;
    top: 50%;
    transform: translateY(-70%);

    clip-path: url(#crystalPath);
    height: 200%;
    width: 20px;
    filter: brightness(200%);
    rotate: 45deg;
    background-color: white;
  }

  button:hover {
    animation: tile linear 3s forwards;
  }

  button:hover:before {
    animation: shine linear 3s forwards;
  }

  @keyframes tile {
    0%,
    50%,
    100% {
      transform: skewX(0);
    }
    25% {
      transform: skewX(5deg);
    }
    75% {
      transform: skewX(-5deg);
    }
  }

  @keyframes shine {
    0%,
    50% {
      rotate: 45deg;
      left: -200%;
    }
    51%,
    100% {
      rotate: -45deg;
      left: 0%;
    }
    25% {
      rotate: 60deg;
      left: 100%;
    }
    75% {
      rotate: -30deg;
      left: 300%;
    }
  }
</style>
