<script>
import Button from './Button.svelte';

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
    <input bind:value={newLabel} />
    <Button on:click={addLabel}>+</Button>
  </li>
  {#each field.labels as label}
    <li>{label.name}</li>
  {/each}
</ul>


