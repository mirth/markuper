<script>
import Button from 'svelte-atoms/Button.svelte';
import Input from 'svelte-atoms/Input.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Chip from 'svelte-atoms/Chip.svelte';

export let field;
let newLabel = '';

function addLabel() {
  field.labels = field.labels.concat({
    value: newLabel,
    name: newLabel,
  });
  newLabel = '';
}
function removeLabel(index) {
  return () => {
    field.labels = field.labels.filter((_unused, iter) => iter !== index);
  };
}
</script>

<Row>
<Cell>
  <Input bind:value={newLabel} placeholder='Label goes here...' size="small" />
</Cell>
<Cell>
  <Button on:click={addLabel} iconRight='plus'>Add label</Button>
</Cell>
</Row>

<Row>
<Cell>
<ul>
  {#each field.labels as label, i}
    <li>
      <Chip text={label.name} onClose={removeLabel(i)} />
    </li>
  {/each}
</ul>
{#if field.labels.length === 0}
  <span>There should be at least one label</span>
{/if}
</Cell>
</Row>

<style>
ul {
  padding: 0;
}

li {
  list-style-type: none;
  border: 1px solid black;
  border-radius: 15px;
  display: inline;
  margin-right: 10px;
}

span {
  color: red;
}
</style>