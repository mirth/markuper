<script>
import Button from 'svelte-atoms/Button.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import ControlRadio from './ControlRadio.svelte';
import ControlCheckbox from './ControlCheckbox.svelte';
import { assessState } from '../store';
import { getFieldsInOrderFor } from '../template';

export let submitMarkupAndFetchNext;
export let sample;

const [fields, groupsOrder] = getFieldsInOrderFor(sample.template);

let fieldIter = 0;
let keyDown = null;

function focusOnNextField(lastGroupTouched) {
  const lastGroupTouchedIndex = groupsOrder.indexOf(lastGroupTouched);
  if(fieldIter === lastGroupTouchedIndex) {
    fieldIter += 1;
  }
}

function handleKeydown(event) {
  keyDown = event.key;
}

function handleKeyup(event) {
  if (event.key !== keyDown) {
    return;
  }

  if (event.key === 'Enter') {
    submitMarkupAndFetchNext();
  }
}


$: $assessState.focusedGroup = groupsOrder[fieldIter];

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>

{#each fields as field, i }
<div class:selected={$assessState.focusedGroup === field.group} id={`device${i}`}>
  {#if field.type === 'radio'}
    <ControlRadio
      {field}
      handleFieldCompleted={focusOnNextField}
      markup={sample.markup && sample.markup.markup[field.group]}
      />
  {/if}
  {#if field.type === 'checkbox'}
    <ControlCheckbox
      {field}
      handleFieldCompleted={focusOnNextField}
      markup={sample.markup && sample.markup.markup[field.group]}
      />
  {/if}
</div>
{/each}
<Spacer size={16} />
<Button on:click={submitMarkupAndFetchNext}>Submit</Button>

<style>

.selected {
  border: 1px solid red;
}
</style>
