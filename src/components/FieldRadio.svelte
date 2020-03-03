<script>
import Button from 'svelte-atoms/Button.svelte';
import Input from 'svelte-atoms/Input.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';

export let field;
let newLabel = '';

function addLabel() {
  field.labels = field.labels.concat({
    value: newLabel,
    name: newLabel,
  });
  newLabel = '';
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

<ul>
  {#each field.labels as label}
    <li>{label.name}</li>
  {/each}
</ul>


