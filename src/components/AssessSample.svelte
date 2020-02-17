<script>

import api from '../api';
import Button from './Button.svelte';
import PageBlank from './PageBlank.svelte';

async function fetchNext() {
  const res = await api.get('/next');
  return res;
}

let sample = fetchNext();

function makeHandleAssess(label) {
  return async () => {
    sample = await sample;
    await api.post('/assess', {
      sample_id: sample.sample_id,
      sample_markup: {
        markup: {
          label,
        },
      },
    });
    sample = fetchNext();
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