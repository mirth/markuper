
<script>
import { link } from 'svelte-spa-router';
import api from '../api';
import PageBlank from './PageBlank.svelte';

export let params = {};

$: project = api.get(`/project/${params.project_id}`);
$: assessedList = api.get(`/project/${params.project_id}/assessed`);

</script>



<PageBlank>
{#await project then p}
<h3>{p.description.name}</h3>
<h3>{p.template.task}</h3>
<h3>{p.data_source.source_uri}</h3>
<h3>
  <a href={`/project/${p.project_id}/assess_sample`} use:link>Begin assess</a>
</h3>
{/await}

{#await assessedList then list}
<ul>
  {#each list.list as forSample}
    <li>
      <p>Sample ID: {forSample.sample_id.sample_id}|Value: {forSample.sample_markup.markup.label.value}</p>
    </li>
  {/each}
</ul>
{/await}
</PageBlank>




