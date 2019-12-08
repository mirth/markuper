<script>
import axios from 'axios';
import config from './config';

async function get() {
  const url = `${config.BACKEND_URL}/api/v1/images`;
  const resp = await axios.get(url);
  return resp.data;
}

const name = 'new world';
const image = get();

</script>

<h1>Hello {name}!</h1>

{#await image}
<p>...waiting</p>
{:then data}
<img src="file://{data.urls[0]}" alt="KEK"/>
{:catch error}
	<p style="color: red">{error}</p>
{/await}
