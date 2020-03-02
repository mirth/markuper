<script>
import DataSource from './DataSource.svelte';
import Button from 'svelte-atoms/Button.svelte';

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

<ul>
  {#each dataSources.dataSources as dataSource}
  <li>
    <DataSource {dataSource} />
  </li>
  {/each}

  <li>
    <DataSource dataSource={newSource} />
    <Button on:click={addDataSource} iconRight='plus'>Add source</Button>
  </li>
</ul>
