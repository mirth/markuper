<script>
import Row from 'svelte-atoms/Grids/Row.svelte';

import Button from 'svelte-atoms/Button.svelte';
import api from '../api';
import PageBlank from './PageBlank.svelte';

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
li {
  display: inline;
}

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
<br />
{#await sample}
<p>...waiting</p>
{:then sample}

<Row>
  {#each sample.template.radios as field}
    <Row>
      <ul>
        {#each field.labels as label}
          <li>
            <Button on:click={makeHandleAssess(field, label)} style="display: inline; min-width: 60px;">
              {label.name}
            </Button>
          </li>
        {/each}
      </ul>
    </Row>
  {/each}
</Row>
<div class="image-container">
  <img src="file://{sample.sample.image_uri}" alt="KEK"/>
</div>
{:catch error}
	<p style="color: red">{error}</p>
{/await}
</PageBlank>