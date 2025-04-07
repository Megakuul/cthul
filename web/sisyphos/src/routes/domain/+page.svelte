<!-- just testing if dynamic client side operations work with adapter-static -->
<script>
  let some = $state("");

  let url = $state("");
  async function doSomething() {
    try {
      const result = fetch(url);
      const response = await result;
      some = response.ok ? "OK" : "NOTOK"
    } catch {
      some = "alarm";
    }
  }
</script>

<h1>Domain</h1>
<input bind:value={url} />
<button onclick={doSomething}>Press to request</button>
<div class="bg-slate-400" >{some}</div>
