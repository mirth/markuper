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
    isValid() {
      const uri = this.source_uri.trim();
      return uri.length !== 0 && isGlob(uri);
    },
  };
}

let newSource = newEmptySource();

let isNewSourceValid = false;
$: isNewSourceValid = newSource.isValid();

// fixme e2e
function addDataSource() {
  if (!isNewSourceValid) {
    return;
  }

  newSource.source_uri = newSource.source_uri.trim();
  dataSources = dataSources.concat([newSource]);
  newSource = newEmptySource();
}

const dirname = (filename) => filename.match(/(.*)[/\\]/)[1] || ''; // fixme fuck

function getDirectory(e) {
  const { files } = e.target;
  if (files.length === 0) {
    return;
  }
  const dir = dirname(files[0].path);
  newSource.source_uri = dir;
}

</script>

<Typography type='subheader' block>Add data source:</Typography>
<Row>
  <Cell>
    <input type='file' on:change={getDirectory} webkitdirectory directory multiple/>
  </Cell>
  <Cell xs={9}>
    <DataSource bind:dataSource={newSource} />
  </Cell>
  <Cell xs={3}>
    <Button
      on:click={addDataSource} iconRight='plus'
      type={isNewSourceValid ? 'filled' : 'flat'}
      >
      Add source
    </Button>
  </Cell>
</Row>
<Row>
<Cell>
  <ul>
    {#each dataSources as dataSource}
    <li>
      <DataSource {dataSource} disabled={true} />
    </li>
    {/each}
  </ul>
</Cell>
</Row>
<Row>
<Cell>
{#if dataSources.length === 0}
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