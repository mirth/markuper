<script>

import api from './api';

async function fetchNext() {
  const res = await api.get('/next');
  return res;
}

let image = fetchNext();

function makeHandleAssess(label) {
  return function() {
    api.post('/assess', {
        "sample_id": {
          "project_id": "project0",
          "sample_id":0,
        },
        "sample_markup": {
          "markup":{
            "label": label
          }
        }
    })
    image = api.get('/next');
  }
}

const labels = [
  'cat',
  'dog',
  'kek',
]
</script>

{#each labels as label}
  <button on:click={makeHandleAssess(label)}>{label}</button>
{/each}

<br />
{#await image}
<p>...waiting</p>
{:then data}
<img src="file://{data.sample_uri}" alt="KEK"/>
{:catch error}
	<p style="color: red">{error}</p>
{/await}
