<script>
import _ from 'lodash';
import { onMount } from 'svelte';
import jsFileDownload from 'js-file-download';
import { push } from 'svelte-spa-router';
import Button from 'svelte-atoms/Button.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Block from 'svelte-atoms/Block.svelte';
import api from '../api';
import PageBlank from './PageBlank.svelte';
import { activeProject, fetchProject } from '../store';
import { getFieldsInOrderFor } from '../template';

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
  return Object.keys(markup)
    .sort()
    .map((labelName) => [labelName, markup[labelName]])
    .map(([labelName, labelValue]) => `${labelName}: ${labelValue}`)
    .join(', ');
}

let assessedList = [];
$: assessedList = _.orderBy($activeProject.assessed.list, 'sample_markup.created_at', 'desc');

onMount(async () => {
  await fetchProject(params.project_id);
});

let fields = [];
$: [fields, groupsOrder] = getFieldsInOrderFor($activeProject.template);
</script>



<PageBlank>
<Block>
<Row>
<Cell>
<Row>
  <Cell>
    <Typography type='title' block>{$activeProject.description.name}</Typography>

    {#each $activeProject.data_sources as src}
      <p>Data source: <span>{src.source_uri}</span></p>
    {/each}
    {#each fields as field}
      <p>Labels: <span>{labelsStr(field)}</span></p>
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
        Export CSV
      </Button>
    </div>
  </Cell>
</Row>
</Cell>
</Row>
<Spacer size={32} />
<Row>
<Cell>
{#if assessedList.length === 0}
<Typography type='title' block>No samples have assessed yet</Typography>
{:else}
<Typography type='title' block>Assessed samples:</Typography>
{/if}
<ul>
  {#each assessedList as forSample}
    <li>
      <p>
        <a href={`#/project/${forSample.sample_id.project_id}/samples/${forSample.sample_id.sample_id}`}>
          <small>{forSample.sample_uri}: </small>
        </a>
        <span>{formatMarkup(forSample.sample_markup.markup)}</span>
      </p>
    </li>
  {/each}
</ul>
</Cell>
</Row>
</Block>
</PageBlank>

<style>
li {
  list-style-type: none;
}
ul {
  padding: 0;
}
</style>