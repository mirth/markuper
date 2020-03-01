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
  dataSources.push(newSource);
  newSource = newEmptySource();
}

</script>

<ul>
  <li>
    <DataSource dataSource={newSource} />
    <Button on:click={addDataSource} iconRight='plus'>Add Source</Button> -->
  </li>

  {#each dataSources as dataSource}
  <li>
    <DataSource {dataSource} />
  </li>
  {/each}
</ul>
