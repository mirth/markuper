<script>
import Button from 'svelte-atoms/Button.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import api from '../api';
import PageBlank from './PageBlank.svelte';
import ControlDevice from './controls/ControlDevice.svelte';
import SampleView from './SampleView.svelte';
import { sampleMarkup, assessState } from '../store';
import { goToProject } from '../project';

export let params = {};

async function fetchNext(projectID) {
  const res = await api.get(`/project/${projectID}/next`);
  return res;
}

async function getProjStats(projectID) {
  const res = await api.get(`/project/${projectID}/stats`);
  return res;
}

let sample = (async () => {
  if (Object.prototype.hasOwnProperty.call(params, 'sample_id')) {
    return api.get(`/project/${params.project_id}/samples/${params.sample_id}`);
  }

  return fetchNext(params.project_id);
})();
let projStats = getProjStats(params.project_id);

async function submitMarkupAndFetchNext() {
  sample = await sample;
  const { sample_id: sampleId } = sample;

  await api.post(`/project/${sampleId.project_id}/assess`, {
    sample_id: sampleId,
    sample_markup: {
      markup: $sampleMarkup,
    },
  });
  sample = fetchNext(sampleId.project_id);
  $assessState.focusedGroup = null;

  projStats = getProjStats(sampleId.project_id);
}

</script>

<PageBlank>
<Row>
{#await sample then sample}
  <Cell xs={8}>
    {#if sample.sample === null}
      <Typography type="title" block>No samples left</Typography>
    {:else}
      <SampleView {sample} />
    {/if}
  </Cell>
  <Cell xs={4}>
    <ControlDevice {sample} submitMarkupAndFetchNext={
      sample.sample === null ? () => {} : submitMarkupAndFetchNext
    } />
    <hr />
    <span>
      Project:
      <Button
        type='empty'
        on:click={goToProject(sample.project.project_id)}
        style='padding: 0; display: inline;'
      >
        {sample.project.description.name}
      </Button>
      {#await projStats then projStats}
        ({projStats.assessed_number_of_samples} / {projStats.total_number_of_samples})
      {/await}
    </span>
  </Cell>
{/await}
</Row>
</PageBlank>


<style>


hr {
  border: none;
  background-color: lightgray;
  height: 1px;
}

</style>
