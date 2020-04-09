<script>
import ControlRadio from './ControlRadio.svelte';
import ControlCheckbox from './ControlBoundingBox.svelte';

export let field;

</script>


{#if field.type === 'radio'}
  <ControlRadio {field} />
{/if}
{#if field.type === 'checkbox'}
  <ControlCheckbox {field} />
{/if}
