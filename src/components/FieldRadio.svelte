<script>
import Button from 'svelte-atoms/Button.svelte';

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

<ul>
  <li>
    <input bind:value={newLabel} placeholder='Label goes here...' />
    <Button on:click={addLabel}>+</Button>
  </li>
  {#each field.labels as label}
    <li>{label.name}</li>
  {/each}
</ul>


