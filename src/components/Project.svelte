
<script>
import _ from 'lodash';
import { push } from 'svelte-spa-router';
import Button from 'svelte-atoms/Button.svelte';
import Typography from "svelte-atoms/Typography.svelte";
import Block from "svelte-atoms/Block.svelte";
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte'
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


<PageBlank>
<Block>
<Row>
<Cell>
{#await project then p}
  <Row>
    <Cell>
      <Typography type='title' block>{p.description.name}</Typography>

      <p>Template: <b>{p.template.task}</b></p>
      <p>Data source: <span>{p.data_source.source_uri}</span></p>
      {#each p.template.radios as radio}
        <p>Labels: <span>{labelsStr(radio)}</span></p>
      {/each}
    </Cell>
  </Row>
  <Spacer size={24} />
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
{/await}
</Cell>
</Row>
<Spacer size={32} />
<Row>
<Cell>
{#await assessedList then list}
  <ul>
    {#each list.list as forSample}
      <li>
        <p>Sample ID: {forSample.sample_id.sample_id}|Value: {forSample.sample_markup.markup.label.value}</p>
      </li>
    {/each}
  </ul>
{/await}
</Cell>
</Row>
</Block>
</PageBlank>




