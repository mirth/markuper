<script>
import ControlRadio from './ControlRadio.svelte';

export let handleFieldCompleted;
export let sample;

const fields = sample.template.radios;

</script>

{#each fields as field}
  <ControlRadio
    {field}
    {handleFieldCompleted}
    markup={sample.markup && sample.markup.markup[field.name.value]}
    />
{/each}
