<script>
import Button from 'svelte-atoms/Button.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import { fetchProjectStats, activeSample, fetchNextSample } from '../../store';

export let projectId;

$: projStats = fetchProjectStats(projectId);

async function skipSample(stats) {
  if (stats.assessed_number_of_samples === stats.total_number_of_samples) {
    return;
  }

  $activeSample = fetchNextSample(projectId);
}

</script>

{#await projStats then projStats}
  <Row>
    <Button
      on:click={() => skipSample(projStats)}
      type='flat'
      iconRight='chevron-right'
      style='display: inline'
      disabled={projStats.assessed_number_of_samples === projStats.total_number_of_samples}>
      Skip
    </Button>
  </Row>
{/await}