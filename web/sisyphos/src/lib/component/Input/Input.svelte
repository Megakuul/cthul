<script>
  import { cn } from "$lib/utils";

  let {
    title,
    type,
    class: classNames,
    value = $bindable(),
  } = $props()
</script>

<p class={cn("relative", classNames)}>
  <span class="absolute -top-2 left-1 text-xs font-bold">{title}</span>
  <input placeholder={title} type={type} bind:value={value}
    class="text-xl w-full p-1 rounded-md bg-slate-50/10 focus:bg-slate-50/20 focus:outline-0 transition-all overflow-hidden" />
</p>