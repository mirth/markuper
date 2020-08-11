<script>
import Button from 'svelte-atoms/Button.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import api from '../api';
import PageBlank from './PageBlank.svelte';
import ControlDevice from './controls/ControlDevice.svelte';
import SampleView from './SampleView.svelte';
import {
  sampleMarkup, assessState, activeSample, fetchNextSampleAndResetState, fetchSampleByIdAndResetState,
} from '../store';
import { goToProject, getProjectIDFromSampleID } from '../project';
import Skip from './sample_view/Skip.svelte';
import Progress from './sample_view/Progress.svelte';

export let params = {};



async function submitMarkupAndFetchNext() {
  const sample = await $activeSample;
  const { sample_id: sampleId } = sample;

  const projID = getProjectIDFromSampleID(sampleId);
  await api.post(`/project/${projID}/assess`, {
    sample_id: sampleId,
    sample_markup: {
      markup: $sampleMarkup,
    },
  });


  $activeSample = fetchNextSampleAndResetState(projID);
}


if (Object.prototype.hasOwnProperty.call(params, 'sample_id')) {
  $activeSample = fetchSampleByIdAndResetState(params.project_id, params.sample_id);
} else {
  $activeSample = fetchNextSampleAndResetState(params.project_id);
}

</script>

<PageBlank>
<Row>
  {#if $activeSample != null}
    {#await $activeSample then sample}
      <Cell xs={8}>
        {#if sample.sample == null}
          <Typography type="title" block>No samples left</Typography>
        {:else}
          <SampleView sample={sample} />
        {/if}
      </Cell>
      <Cell xs={4}>
        {#if sample.sample != null}
          <ControlDevice
            sample={sample}
            submitMarkupAndFetchNext={sample.sample == null ? () => {} : submitMarkupAndFetchNext} />
          <hr />
        {/if}
        <p>
          Project:
          <Button
            type='empty'
            on:click={goToProject(sample.project.project_id)}
            style='padding: 0; display: inline;'
          >
            {sample.project.description.name}
          </Button>
          <Progress projectId={params.project_id} />
        </p>
        {#if sample.sample != null}
          <p>Sample:<small id='sample_uri'>{sample.sample.media_uri}</small></p>
        {/if}
        <Skip projectId={params.project_id} />
      </Cell>
    {/await}
  {/if}
</Row>
</PageBlank>


<style>


hr {
  border: none;
  background-color: lightgray;
  height: 1px;
}

small {
  word-wrap: break-word;
  margin-left: 10px;
}

</style>
