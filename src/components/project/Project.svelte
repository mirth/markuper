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
import api from '../../api';
import PageBlank from '../PageBlank.svelte';
import { activeProject, fetchProject } from '../../store';
import ControlList from '../controls/ControlList.svelte';
import ProjectAssessedSamplesList from './ProjectAssessedSamplesList.svelte';


export let params;

function exportProject(p) {
  return async () => {
    const res = api.downloadFile(`/project/${p.project_id}/export`);
    const [data, filename] = await res;
    jsFileDownload(data, filename);
  };
}

let assessedList = [];
$: assessedList = _.orderBy($activeProject.assessed.list, 'sample_markup.created_at', 'desc');

onMount(async () => {
  await fetchProject(params.project_id);
});

</script>



<PageBlank>
<Block>
<Row>
  <Cell xs={8}>
    <Row>
      <Cell>
        <Typography type='title' block><b>{$activeProject.description.name}</b></Typography>
        <Typography type='subheader' block>Data sources:</Typography>
        <ul>
        {#each $activeProject.data_sources as src}
          <li><span>{src.source_uri}</span></li>
        {/each}
        </ul>
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
  <Cell xs={4}>
    <Typography type='subheader' block>Sample UI</Typography>
    <Block type="block3">
      <ControlList
        submitMarkupAndFetchNext={() => {}}
        owner={$activeProject.template}
        />
    </Block>
  </Cell>
</Row>
<Spacer size={32} />
<Row>
<Cell>
{#if assessedList.length === 0}
<Typography type='title' block>No samples have been assessed yet</Typography>
{:else}
<Typography type='title' block>Assessed samples:</Typography>
{/if}

<ProjectAssessedSamplesList {assessedList} />

</Cell>
</Row>
</Block>
</PageBlank>

<style>
li {
  /* list-style-type: none; */
  list-style: circle inside none;
}
ul {
  padding: 0;
}
</style>