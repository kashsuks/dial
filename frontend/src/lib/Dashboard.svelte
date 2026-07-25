<script>
  import { onMount } from 'svelte';
  import { CaretDown } from 'phosphor-svelte';
  import Tracker from './Tracker.svelte';
  import PieChart from './PieChart.svelte';
  import BarChart from './BarChart.svelte';
  import StatsRow from './StatsRow.svelte';
  import { GetTagBreakdown, GetDailyTotals, GetStats } from '../../wailsjs/go/gui/App';

  const ranges = [
    { value: 'today', label: 'Today' },
    { value: 'week', label: 'This Week' },
    { value: 'month', label: 'This Month' },
    { value: 'year', label: 'This Year'},
  ];

  let selectedRange = 'week';
  let tagData = [];
  let dailyData = [];
  let stats = null;

  async function loadReports() {
    [tagData, dailyData, stats] = await Promise.all([
      GetTagBreakdown(selectedRange),
      GetDailyTotals(selectedRange),
      GetStats(selectedRange),
    ]);
  }

  $: selectedRange, loadReports();

  onMount(loadReports);

  // Let Tracker.sveltes start/stop/pause trigger a refresh too
  function handleSessionChange() {
    loadReports();
  }
</script>

<div class="dashboard">
  <div class="top-row">
    <div class="panel tracker-panel">
      <Tracker on:change={handleSessionChange} />
    </div>

    <div class="panel pie-panel">
      <div class="panel-header">
        <span class="panel-title">By Tag</span>
        <div class="select-wrap">
          <select bind:value={selectedRange}>
            {#each ranges as r}
              <option value={r.value}>{r.label}</option>
            {/each}
          </select>
          <CaretDown size={12} weight="bold" />
        </div>
      </div>
      <PieChart data={tagData} />
    </div>
  </div>

  <div class="panel">
    <div class="panel-header">
      <span class="panel-title">Time Logged</span>
    </div>
    <BarChart data={dailyData} />
  </div>

  <StatsRow {stats} />
</div>

<style>
  .dashboard {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .top-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }
  .panel {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 24px;
  }
  .tracker-panel, .pie-panel {
    padding: 0;
  }
  .tracker-panel :global(.tracker-card),
  .pie-panel {
    padding: 24px;
    border: none;
  }
  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 14px;
  }
  .panel-title {
    font-size: 13px;
    font-weight: 600;
    color: var(--ink);
  }
  .select-wrap {
    position: relative;
    display: flex;
    align-items: center;
  }
  select {
    appearance: none;
    background: var(--canvas);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 5px 24px 5px 10px;
    font-size: 11px;
    font-family: var(--font-sans);
    color: var(--text);
    cursor: pointer;
  }
  .select-wrap :global(svg) {
    position: relative;
    right: 8px;
    pointer-events: none;
    color: var(--text-muted);
  }
</style>
