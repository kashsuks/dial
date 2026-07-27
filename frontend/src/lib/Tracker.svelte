<script>
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { Play, Pause, Stop, Tag, FolderSimple } from 'phosphor-svelte';
  import {
    StartSession,
    StopSession,
    PauseSession,
    ResumeSession,
    CurrentSession,
  } from '../../wailsjs/go/gui/App';

  const dispatch = createEventDispatcher();

  let session = null;
  let taskInput = '';
  let projectInput = '';
  let tagsInput = '';
  let elapsedDisplay = '00:00:00';
  let tickHandle;
  let pollHandle;

  function formatElapsed(seconds) {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = Math.floor(seconds % 60);
    const pad = (n) => String(n).padStart(2, '0');
    return `${pad(h)}:${pad(m)}:${pad(s)}`
  }

  function startTick() {
    clearInterval(tickHandle);
    tickHandle = setInterval(() => {
      if (session && !session.isPaused) {
        session.elapsedSeconds += 1;
        elapsedDisplay = formatElapsed(session.elapsedSeconds);
      }
    }, 1000);
  }

  async function refresh() {
    const s = await CurrentSession();
    session = s;
    if (s) elapsedDisplay = formatElapsed(s.elapsedSeconds);
  }

  async function handleStart() {
    if (!taskInput.trim()) return;
    session = await StartSession(taskInput.trim(), projectInput.trim(), tagsInput.trim());
    elapsedDisplay = formatElapsed(session.elapsedSeconds);
    taskInput = '';
    projectInput = '';
    tagsInput = '';
    dispatch('change');
  }

  async function handlePause() {
    session = await PauseSession();
    dispatch('change');
  }

  async function handleResume() {
    session = await ResumeSession();
    dispatch('change');
  }

  async function handleStop() {
    await StopSession();
    session = null;
    elapsedDisplay = '00:00:00';
    dispatch('change');
  }

  onMount(() => {
    refresh();
    startTick();
    pollHandle = setInterval(refresh, 10000);
  });

  onDestroy(() => {
    clearInterval(tickHandle);
    clearInterval(pollHandle);
  });
</script>

<div class="tracker-card">
  {#if session}
    <div class="running-state">
      <div class="running-header">
        <span class="task-name">{session.task}</span>
        {#if session.isPaused}
          <span class="badge badge-yellow">Paused</span>
        {:else}
          <span class="badge badge-green">Running</span>
        {/if}
      </div>

      {#if session.project || session.tags}
        <div class="meta-row">
          {#if session.project}
            <span class="meta-item"><FolderSimple size={13} weight="bold" /> {session.project}</span>
          {/if}
          {#if session.tags}
            <span class="meta-item"><Tag size={13} weight="bold" /> {session.tags}</span>
          {/if}
        </div>
      {/if}

      <div class="elapsed">{elapsedDisplay}</div>

      <div class="controls">
        {#if session.isPaused}
          <button class="btn btn-primary" on:click={handleResume}>
            <Play size={15} weight="fill" /> Resume
          </button>
        {:else}
          <button class="btn btn-secondary" on:click={handlePause}>
            <Pause size={15} weight="fill" /> Pause
          </button>
        {/if}
        <button class="btn btn-stop" on:click={handleStop}>
          <Stop size={15} weight="fill" /> Stop
        </button>
      </div>
    </div>
  {:else}
    <div class="idle-state">
      <input
        class="input input-task"
        type="text"
        placeholder="What are you working on?"
        bind:value={taskInput}
        on:keydown={(e) => e.key === 'Enter' && handleStart()}
      />
      <div class="input-row">
        <input class="input" type="text" placeholder="Project" bind:value={projectInput} />
        <input class="input" type="text" placeholder="Tags" bind:value={tagsInput} />
      </div>
      <button class="btn btn-primary btn-full" on:click={handleStart} disabled={!taskInput.trim()}>
        <Play size={15} weight="fill" /> Start
      </button>
    </div>
  {/if}
</div>

<style>
  .tracker-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: clamp(18px, 2.2vw, 32px);
    width: 100%;
  }

  .running-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 8px;
  }

  .task-name {
    font-size: clamp(15px, 1.4vw, 19px);
    font-weight: 600;
    color: var(--ink);
  }

  .badge {
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    padding: 3px 10px;
    border-radius: 9999px;
  }

  .badge-green { background: var(--pastel-green-bg); color: var(--pastel-green-text); }
  .badge-yellow { background: var(--pastel-yellow-bg); color: var(--pastel-yellow-text); }

  .meta-row {
    display: flex;
    gap: 14px;
    margin-top: 10px;
  }

  .meta-item {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    line-height: 1;
    color: var(--text-muted);
  }
  
  .elapsed {
    font-family: var(--font-mono);
    font-size: clamp(32px, 4vw, 52px);
    font-weight: 500;
    color: var(--ink);
    letter-spacing: -0.01em;
    margin: clamp(16px, 2.4vw, 28px) 0;
  }

  .controls {
    display: flex;
    gap: 8px;
  }

  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    border: none;
    border-radius: 6px;
    padding: 10px 16px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: background 150ms, transform 100ms;
  }

  .btn:active { transform: scale(0.98); }
  .btn:disabled { opacity: 0.4; cursor: not-allowed; }

  .btn-primary {
    background: var(--ink);
    color: #fff;
  }
  .btn-primary:hover:not(:disabled) { background: var(--ink-hover); }

  .btn-secondary {
    background: var(--pastel-yellow-bg);
    color: var(--pastel-yellow-text);
  }

  .btn-stop {
    background: var(--pastel-red-bg);
    color: var(--pastel-red-text);
  }

  .btn-full { width: 100%; margin-top: 14px; }

  .idle-state .input-task {
    font-size: clamp(14px, 1.2vw, 16px);
    padding: clamp(10px, 1vw, 14px) clamp(12px, 1.2vw, 16px);
  }

  .input-row {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }

  .input {
    width: 100%;
    min-width: 0;
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: clamp(9xp, 0.9vw, 12px) clamp(10px, 1vw, 14px);
    font-size: clamp(12px, 1vw, 14px);
    color: var(--text);
    background: var(--canvas);
    outline: none;
    transition: border-color 150ms;
  }
  .input:focus { border-color: #111111; }
  .input::placeholder { color: var(--text-muted); }
</style>
