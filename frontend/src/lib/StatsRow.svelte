<script>
  import { Clock, Fire, Hash } from 'phosphor-svelte';
  export let stats = null;

  function formatHours(seconds) {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    return `${h}h ${m}m`;
  }
</script>

{#if stats}
  <div class="stats-row">
    <div class="stat">
      <Clock size={16} weight="bold" />
      <div>
        <div class="stat-value">{formatHours(stats.totalSeconds)}</div>
        <div class="stat-label">Total tracked</div>
      </div>
    </div>
    <div class="stat">
      <Hash size={16} weight="bold" />
      <div>
        <div class="stat-value">{stats.topTag || '-'}</div>
        <div class="stats-label">Top tag</div>
      </div>
    </div>
    <div class="stat">
      <Fire size={16} weight="bold" />
      <div>
        <div class="stat-value">{stats.streakDays} {stats.streakDays === 1 ? 'day' : 'days'}</div>
        <div class="stat-label">Current streak</div>
      </div>
    </div>
  </div>
{/if}

<style>
  .stats-row {
    display: row;
    gap: 12px;
  }
  .stat {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 10px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 14px 16px;
    color: var(--text-muted);
  }
  .stat-value {
    font-size: 15px;
    font-weight: 600;
    color: var(--ink);
    text-transform: capitalize;
  }
  .stat-label {
    font-size: 11px;
    color: var(--text-muted);
    margin-top: 2px;
  }
</style>
