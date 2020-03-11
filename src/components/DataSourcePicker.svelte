<script>
import isGlob from 'is-valid-glob';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Typography from 'svelte-atoms/Typography.svelte';
import Button from 'svelte-atoms/Button.svelte';
import DataSource from './DataSource.svelte';

export let dataSources;
function newEmptySource() {
  return {
    type: 'local_directory',
    source_uri: '',
    isValid: function() {
      const uri = this.source_uri.trim()
      return uri.length !== 0 && isGlob(uri)
    },
  };
}

let newSource = newEmptySource();

// fixme e2e
function addDataSource() {
  if(!newSource.isValid()) {
    return
  }

  newSource.source_uri = newSource.source_uri.trim();
  dataSources.dataSources = dataSources.dataSources.concat([newSource]);
  newSource = newEmptySource();
}

</script>

<Typography type='subheader' block>Add data source:</Typography>
<Row>
  <Cell xs={7}>
    <DataSource dataSource={newSource} />

  </Cell>
  <Cell xs={3}>
    <Button on:click={addDataSource} iconRight='plus'>Add source</Button>
  </Cell>
</Row>
<Row>
<Cell>
  <ul>
    {#each dataSources.dataSources as dataSource}
    <li>
      <DataSource {dataSource} disabled={true} />
    </li>
    {/each}
  </ul>
</Cell>
</Row>
<Row>
<Cell>
{#if dataSources.dataSources.length === 0}
<span>There should be at least one data source</span>
{/if}
</Cell>
</Row>

<style>
li {
  list-style-type: none;
  margin-bottom: 15px;
}
ul {
  padding: 0;
}

span {
  color: red;
}
</style>