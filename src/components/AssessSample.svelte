<script>
import Spacer from 'svelte-atoms/Spacer.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import api from '../api';
import PageBlank from './PageBlank.svelte';
import ControlDevice from './ControlDevice.svelte';
import { activeMarkup } from '../store';

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
      markup: $activeMarkup,
    },
  });
  sample = fetchNext(sampleId.project_id);
}

</script>

<style>

.image-container {
  padding: 0 45px 45px 0;
  width: 100%;
}

img {
  max-width: 100%;
  border: 1px solid black;
}

</style>

<PageBlank>
{#await sample then sample}
{#if sample.sample === null}
<Typography type="title" block>No samples left</Typography>
{:else}
<ControlDevice {sample} {submitMarkupAndFetchNext} />
<Spacer size={24} />
<div class='image-container'>
  <img src='file://{sample.sample.image_uri}' alt='KEK'/>
</div>
{/if}
{/await}
</PageBlank>