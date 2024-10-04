<script>
    import CthulIcon from "$lib/components/CthulIcon.svelte";

</script>


<svelte:head>
	<title>Cthul Project</title>
	<meta name="description" content="Cthul Project - A Next Gen Virtualization Platform." />
	<meta property="og:description" content="Cthul Project - A Next Gen Virtualization Platform." />
	<meta property="og:title" content="Cthul Project">
  <meta property="og:type" content="website">
	<meta property="og:image" content="https://cthul.io/favicon.png" />
	<link rel="canonical" href="https://cthul.io" />
</svelte:head>


<h1 class="text-9xl my-20 py-2 font-bold ">
  <span class="text-transparent bg-clip-text bg-gradient-to-br from-[#b10a79] to-[#83105dce] fade-short">Cthul</span>
  <span class="text-transparent bg-clip-text bg-gradient-to-br from-stone-800 to-stone-950 fade-long">Project</span>
</h1>

<CthulIcon class="w-3/4 lg:w-1/2 cursor-pointer" />

<section class="flex flex-col items-center w-full my-20">
  <div class="flex flex-col w-9/12 lg:w-6/12 p-6 bg-stone-800 text-stone-400 rounded-lg">
    <p class="flex flex-row items-start">
      <span class="mr-2 font-bold">{">"}</span>
      <span>
        <span class="font-bold">cthul --status</span><br>
        analyze request....<br>
        load assets....<br>
        specify target....<br>
        check status....<br>
        -------------------<br>
        <span class="text-indigo-500">Name</span>: Cthul Project<br>
        <span class="text-red-800">Source</span>: <a class="underline" href="https://github.com/megakuul/cthul">github.com/megakuul/cthul</a><br>
        <span class="text-orange-800">Status</span>: Work in progress<br>
      </span>
    </p>
    <p class="flex flex-row items-center">
      <span class="mr-2 font-bold">{">"}</span>
      <span class="font-bold h-4 w-2 bg-stone-200 animate-pulse"></span>
    </p>
  </div>
</section>

<style>
  .fade-short {
    opacity: 0;
    animation: fade ease-out 3s 0s forwards;
  }

  .fade-long {
    opacity: 0;
    animation: fade ease-out 3s .5s forwards;
  }

  @keyframes fade {
    0% {
      opacity: 0;
    }
    100% {
      opacity: 1;
    }
  }
</style>