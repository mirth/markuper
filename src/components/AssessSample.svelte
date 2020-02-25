<script>

import api from '../api';
import Button from './Button.svelte';
import PageBlank from './PageBlank.svelte';

export let params = {};

async function fetchNext(projectID) {
  const res = await api.get(`/project/${projectID}/next`);
  return res;
}

let sample = fetchNext(params.project_id);

function makeHandleAssess(label) {
  return async () => {
    sample = await sample;
    const { sample_id } = sample;
    await api.post(`/project/${sample_id.project_id}/assess`, {
      sample_id: sample_id,
      sample_markup: {
        markup: {
          label,
        },
      },
    });
    sample = fetchNext(sample_id.project_id);
  };
}

</script>

<PageBlank>
<br />
{#await sample}
<p>...waiting</p>
{:then sample}
{#each sample.template.radios as field}
  {#each field.labels as label}
    <Button on:click={makeHandleAssess(label)}>{label.name}</Button>
  {/each}
{/each}
<img src="file://{sample.sample.image_uri}" alt="KEK"/>
{:catch error}
	<p style="color: red">{error}</p>
{/await}
</PageBlank>