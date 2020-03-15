<script>
import ControlRadio from './ControlRadio.svelte';

export let handleFieldComplete;
export let sample;

const fields = sample.template.radios;

</script>

{#each fields as field}
  <ControlRadio {field} {handleFieldComplete} />
{/each}