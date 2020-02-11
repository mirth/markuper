
<script>
import { onMount } from 'svelte';
import api from '../api.js';
import PageBlank from './PageBlank.svelte';

export let params = {};

$: project = api.get(`/project/${params.project_id}`)

</script>


{#await project then p}
<PageBlank>
<h3>{p.description.name}</h3>
<h3>{p.template.type}</h3>
</PageBlank>
{/await}
