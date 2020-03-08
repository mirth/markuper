<script>
import _ from 'lodash';
import { onMount } from 'svelte';
import jsFileDownload from 'js-file-download';
import { push } from 'svelte-spa-router';
import api from '../api';
import Button from 'svelte-atoms/Button.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Block from 'svelte-atoms/Block.svelte';
import PageBlank from './PageBlank.svelte';
import { activeProject, fetchProject } from '../store';

export let params;

function labelsStr(radio) {
  return _.map(radio.labels, 'value').join(', ');
}

function exportProject(p) {
  return async () => {
    const res = api.downloadFile(`/project/${p.project_id}/export`);
    const [data, filename] = await res;
    jsFileDownload(data, filename);
  };
}

function formatMarkup(markup) {
  return _(markup).toPairs().map(([labelName, labelValue]) => `${labelName}:${labelValue}`).join('\n');
}

let assessedList = [];
$: assessedList = _.orderBy($activeProject.assessed.list, 'sample_markup.created_at', 'desc')

onMount(async () => {
  await fetchProject(params.project_id);
});
</script>



<PageBlank>
<Block>
<!-- fixme -->
{#if $activeProject.template}
<Row>
<Cell>
<Row>
  <Cell>
    <Typography type='title' block>{$activeProject.description.name}</Typography>

    <p>Template: <b>{$activeProject.template.task}</b></p>
    {#each $activeProject.data_sources as src}
      <p>Data source: <span>{src.source_uri}</span></p>
    {/each}
    {#each $activeProject.template.radios as radio}
      <p>Labels: <span>{labelsStr(radio)}</span></p>
    {/each}
  </Cell>
</Row>
<Spacer size={24} />
<Row>
  <Cell>
    <div style='display: flex; justify-content: space-between; flex-direction: row;'>
      <Button on:click={() => push(`/project/${$activeProject.project_id}/assess_sample`)} iconRight='chevron-right'>
        Begin assess
      </Button>
      <Button on:click={exportProject($activeProject)} iconLeft='download'>
        Export
      </Button>
    </div>
  </Cell>
</Row>
</Cell>
</Row>
<Spacer size={32} />
<Row>
<Cell>
<!-- fixme sort by date -->
<ul>
  {#each assessedList as forSample}
    <li>
      <p>Sample ID: {forSample.sample_id.sample_id}|Value: {formatMarkup(forSample.sample_markup.markup)}</p>
    </li>
  {/each}
</ul>
</Cell>
</Row>
{/if}
</Block>
</PageBlank>
