<script>

import api from './api';
import Button from './components/Button.svelte';
import './global.svelte';

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

const labels = [
  'cat',
  'dog',
  'kek',
];
</script>

{#each labels as label}
  <Button on:click={makeHandleAssess(label)}>{label}</Button>
{/each}

<br />
{#await sample}
<p>...waiting</p>
{:then sample}
<img src="file://{sample.sample_uri}" alt="KEK"/>
{:catch error}
	<p style="color: red">{error}</p>
{/await}
