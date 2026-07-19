<script>
  import Tracker from './lib/Tracker.svelte';
</script>

<main>
  <h1>Dial</h1>
  <Tracker />
</main>

<style>
  main {
    max-width: 420px;
    margin: 0 auto;
    padding: 32px 24px;
  }

  h1 {
    font-family: 'Newsreader', 'Instrument Serif', serif;
    font-size: 22px;
    font-weight: 500;
    letter-spacing: -0.02em;
    color: var(--ink);
    margin: 0 0 20px 0;
  }
</style>
