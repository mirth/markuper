
<script>
import _ from 'lodash';
import { push } from 'svelte-spa-router';
import Button from 'svelte-atoms/Button.svelte';
import Typography from "svelte-atoms/Typography.svelte";
import Block from "svelte-atoms/Block.svelte";
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import jsFileDownload from 'js-file-download';

import PageBlank from './PageBlank.svelte';
import api from '../api';

export let params = {};

$: project = api.get(`/project/${params.project_id}`);
$: assessedList = api.get(`/project/${params.project_id}/assessed`);

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

</script>

<style>
button {
  display: inline;
}
</style>


<PageBlank>
<Block>
<Row>
<Cell>
{#await project then p}
  <Row>
    <Cell>
      <Typography type='title' block>{p.description.name}</Typography>

      <p>Template: <b>{p.template.task}</b></p>
      <p>Data source: {p.data_source.source_uri}</p>
    </Cell>
  </Row>
  <Row>
    <Cell>
      <div style='display: flex; justify-content: space-between; flex-direction: row;'>
        <Button on:click={() => push(`/project/${p.project_id}/assess_sample`)} iconRight='chevron-right'>
          Begin assess
        </Button>
        <Button on:click={exportProject(p)} iconLeft='download'>
          Export
        </Button>
      </div>
    </Cell>
  </Row>
  {#each p.template.radios as radio}
    <h3>{labelsStr(radio)}</h3>
  {/each}
{/await}
</Cell>
</Row>
</Block>


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




