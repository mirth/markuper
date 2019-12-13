<script>
import axios from 'axios';
import config from './config';

async function get() {
  const url = `${config.BACKEND_URL}/api/v1/next`;
  const resp = await axios.get(url);
  return resp.data;
}

async function fetchNext() {
  const res = await get();
  return res;
}

let image = fetchNext();

function handleClick() {
  image = fetchNext();
}

</script>


<button on:click={handleClick}>Next</button>
<br />
{#await image}
<p>...waiting</p>
{:then data}
<img src="file://{data.sample_uri}" alt="KEK"/>
{:catch error}
	<p style="color: red">{error}</p>
{/await}
