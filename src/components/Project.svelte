
<script>
import { link } from 'svelte-spa-router';
import api from '../api';
import PageBlank from './PageBlank.svelte';

export let params = {};

$: project = api.get(`/project/${params.project_id}`);

</script>


{#await project then p}
<PageBlank>
<h3>{p.description.name}</h3>
<h3>{p.template.task}</h3>
<h3>{p.data_source.source_uri}</h3>
<h3>
  <a href={`/project/${project.project_id}/assess_sample`} use:link>Begin assess</a>
</h3>
</PageBlank>
{/await}
