<script>
import Button from 'svelte-atoms/Button.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import api from '../../api';
import { fetchProjectStats, activeSample, fetchNextSampleAndResetState } from '../../store';

export let projectId;

$: projStats = fetchProjectStats(projectId);

async function skipSampleAndFetchNext(stats) {
  if (stats.assessed_number_of_samples === stats.total_number_of_samples) {
    return;
  }

  const sample = await $activeSample;
  const { sample_id: sampleId } = sample;

  await api.post(`/project/${projectId}/assess`, {
    sample_id: sampleId,
    sample_markup: {
      markup: null,
    },
  });

  $activeSample = fetchNextSampleAndResetState(projectId);
}


</script>

{#await projStats then projStats}
  <Row>
    <Button
      on:click={() => skipSampleAndFetchNext(projStats)}
      type='flat'
      iconRight='chevron-right'
      style='display: inline'
      disabled={projStats.assessed_number_of_samples === projStats.total_number_of_samples}>
      Skip
    </Button>
  </Row>
{/await}