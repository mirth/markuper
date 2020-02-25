
<script>
import _ from 'lodash';
import { link } from 'svelte-spa-router';
import api from '../api';
import PageBlank from './PageBlank.svelte';
import Button from 'svelte-atoms/Button.svelte';
import jsFileDownload from 'js-file-download';

import axios from 'axios';
import config from '../config';

export let params = {};

$: project = api.get(`/project/${params.project_id}`);
$: assessedList = api.get(`/project/${params.project_id}/assessed`);

function labelsStr(radio) {
  return _.map(radio.labels, 'value').join(', ');
}

function exportProject(p) {
  return async () => {
    const res = api.downloadFile(`/project/${p.project_id}/export`)
    const [data, filename] = await res;
    jsFileDownload(data, filename)
  }
}
</script>



<PageBlank>
{#await project then p}
<h3>{p.description.name}</h3>
<h3>{p.template.task}</h3>
<h3>{p.data_source.source_uri}</h3>
<h3>
  <a href={`/project/${p.project_id}/assess_sample`} use:link>Begin assess</a>
</h3>
<Button on:click={exportProject(p)}>Export</Button>
{#each p.template.radios as radio}
<h3>{labelsStr(radio)}</h3>
{/each}
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




