<script>
  import Dashboard from './lib/Dashboard.svelte';
</script>

<main>
  <h1>Dial</h1>
  <Dashboard />
</main>

<style>
  main {
    width: 100%;
    max-width: 960px;
    margin: 0 auto;
    padding: clamp(20px, 4vw, 48px) clamp(20px, 5vw, 56px);
  }

  h1 {
    font-family: 'Newsreader', 'Instrument Serif', serif;
    font-size: clamp(20px, 2.4vw, 28px);
    font-weight: 500;
    letter-spacing: -0.02em;
    color: var(--ink);
    margin: 0 0 clamp(16px, 2vw, 24px) 0;
  }
</style>
