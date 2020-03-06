<script>
import api from '../api';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import PageBlank from './PageBlank.svelte';
import ControlDevice from './ControlDevice.svelte';

export let params = {};

async function fetchNext(projectID) {
  const res = await api.get(`/project/${projectID}/next`);
  return res;
}

let sample = fetchNext(params.project_id);

function makeHandleAssess(field, label) {
  return async () => {
    sample = await sample;
    const { sample_id: sampleId } = sample;
    const markup = { [field.name.value]: label.value };

    await api.post(`/project/${sampleId.project_id}/assess`, {
      sample_id: sampleId,
      sample_markup: {
        markup,
      },
    });
    sample = fetchNext(sampleId.project_id);
  };
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
<ControlDevice field={sample.template.radios[0]} {makeHandleAssess}/>
<Spacer size={24} />
<div class='image-container'>
  <img src='file://{sample.sample.image_uri}' alt='KEK'/>
</div>
{/await}
</PageBlank>