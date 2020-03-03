<script>
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
  };
}

let newSource = newEmptySource();

function addDataSource() {
  dataSources.dataSources = dataSources.dataSources.concat([newSource]);
  newSource = newEmptySource();
}

</script>

<Typography type="subheader" block>Add data source:</Typography>
<Row>
  <Cell xs={7}>
    <DataSource dataSource={newSource} />
  </Cell>
  <Cell xs={3}>
    <Button on:click={addDataSource} iconRight='plus'>Add source</Button>
  </Cell>
</Row>
<ul>
  {#each dataSources.dataSources as dataSource}
  <li>
    <DataSource {dataSource} />
  </li>
  {/each}
</ul>