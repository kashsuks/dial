<script>
  export let data = []; // [{ date, seconds }]

  $: max = Math.max(1, ...data.map((d) => d.seconds));

  function formatDay(dateStr) {
    const d= new Date(dateStr + 'T00:00:00');
    return d.toLocaleDateString(undefined, { weekday: 'short' }).slice(0, 2);
  }

  function formatHours(seconds) {
    return (seconds / 3600).toFixed(1);
  }
</script>

<div class="bar-chart"
  {#if data.length === 0}
    <div class="empty">No time logged yet</div>
  {:else}
    {#each data as d}
      <div class="bar-col">
        <div class="bar-track">
          <div class="bar-fill" style="height: {(d.seconds / max) * 100}%" title="{formatHours(d.seconds)}h"></div>
        </div>
        <span class="bar-label">{formatDay(d.date)}</span>
      </div>
    {/each}
  {/if}
</div>

<style>
  .bar-chart {
    display: flex;
    align-items: flex-end;
    gap: 10px;
    height: 120px;
    padding-top: 8px;
  }
  .empty {
    color: var(--text-muted);
    font-size: 13px;
    margin: auto;
  }
  .bar-col {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    flex: 1;
    height: 100%;
  }
  .bar-track {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: flex-end;
  }
  .bar-fill {
    width: 100%;
    background: var(--pastel-blue-bg);
    border-radius: 3px 3px 0 0;
    min-height: 2px;
    transition: height 300ms ease;
  }
  .bar-label {
    font-size: 10px;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
</style>
