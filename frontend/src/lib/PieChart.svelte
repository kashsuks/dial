<script>
  export let data = []; // [{ tag, seconds }]

  const colors = ['#1F6C9F', '#346538', '#956400', '#9F2F2D', '$6B5B95', '#4A7C7C'];

  $: total = data.reduce((sum, d) => sum + d.seconds, 0);

  $: slices = (() => {
    let angle = 0;
    return data.map((d, i) => {
      const fraction = total > 0 ? d.seconds / total : 0;
      const startAngle = angle;
      const endAngle = angle + fraction * 360;
      angle = endAngle;
      return { ...d, fraction, startAngle, endAngle, color: colors[i % colors.length] };
    });
  })();

  function arcPath(startAngle, endAngle, r = 80, cx = 100, cy = 100) {
    const toRad = (deg) => ((deg - 90) * Math.PI) / 180;
    const x1 = cx + r * Math.cos(toRad(startAngle));
    const y1 = cy + r * Math.sin(toRad(startAngle));
    const x2 = cx + r * Math.cos(toRad(endAngle));
    const y2 = cy + r * Math.sin(toRad(endAngle));
    const largeArc = endAngle - startAngle > 180 ? 1 : 0;
    return `M ${cx} ${cy} L ${x1} ${y1} A ${r} A ${r} ${r} 0 ${largeArc} 1 ${x2} ${y2} Z`;
  }

  function formatHours(seconds) {
    const h = seconds / 3600;
    return h < 1 ? `${Math.round(seconds / 60)}m` : `${h.toFixed(1)}h`;
  }
</script>

<div class="pie-wrap">
  {#if total === 0}
    <div class="empty">No time logged yet</div>
  {:else}
    <svg viewBox="0 0 200 200" class="pie-svg">
      {#each slices as slice}
        <path d={arcPath(slice.startAngle, slice.endAngle)} fill={slice.color} opacity="0.85" />
      {/each}
      <circle cx="100" cy="100" r="46" fill="var(--surface)" />
    </svg>
    <div class="legend">
      {#each slices as slice}
        <div class="legend-row">
          <span class="dot" style="background:{slice.color}"></span>
          <span class="legend-tag">{slice.tag}</span>
          <span class="legend-time">{formatHours(slice.seconds)}</span>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .pie-wrap {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }

  .pie-svg { width: 160px; height: 160px; }
  .empty {
    color: var(--text-muted);
    font-size: 13px;
    padding: 40px 0;
  }
  .legend {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .legend-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
  }
  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }
  .legend-tag {
    color: var(--text);
    flex: 1;
    text-transform: capitalize;
  }
  .legend-time {
    color: var(--text-muted);
    font-family: var(--font-mono);
    font-size: 11px;
  }
</style>
