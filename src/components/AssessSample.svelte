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

let sample = (async () => {
  if (Object.prototype.hasOwnProperty.call(params, 'sample_id')) {
    return api.get(`/project/${params.project_id}/samples/${params.sample_id}`);
  }

  return fetchNext(params.project_id);
})();

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
  $assessState.markup = {};
  $assessState.focusedGroup = null;
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
    <ControlDevice {sample} {submitMarkupAndFetchNext} />
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
