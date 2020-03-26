<script>
import Button from 'svelte-atoms/Button.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import ControlRadio from './ControlRadio.svelte';
import ControlCheckbox from './ControlCheckbox.svelte';
import { assessState, activeMarkup } from '../store';
import { getFieldsInOrderFor } from '../template';

export let submitMarkupAndFetchNext;
export let sample;

if (sample.markup) {
  $activeMarkup = sample.markup.markup;
}

const [fields, groupsOrder] = getFieldsInOrderFor(sample.template);
const submitGroup = '[submit]';

groupsOrder.push(submitGroup);

let fieldIter = 0;
let keyDown = null;

function isDeviceFilled() {
  const curField = fields[fieldIter];
  if(curField.type === 'radio' && !Object.hasOwnProperty.call($activeMarkup, curField.group)) {
    return false;
  }

  return true
}

function focusOnNextField() {
  if(!isDeviceFilled()) {
    return;
  }

  fieldIter += 1;
}

function handleKeydown(event) {
  keyDown = event.key;
}

function handleKeyup(event) {
  if (event.key !== keyDown) {
    return;
  }

  // onKeyEnter
  if (event.key === 'Enter') {
    if (fieldIter === groupsOrder.length - 1) {
      submitMarkupAndFetchNext();
    }
    focusOnNextField();
  }
}


$: $assessState.focusedGroup = groupsOrder[fieldIter];

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>

{#each fields as field, i }
<div class:selected={$assessState.focusedGroup === field.group} id={`device${i}`}>
  {#if field.type === 'radio'}
    <ControlRadio {field} />
  {/if}
  {#if field.type === 'checkbox'}
    <ControlCheckbox {field} />
  {/if}
</div>
{/each}
<Spacer size={16} />
<div class:selected={$assessState.focusedGroup === submitGroup} id='device_submit'>
  <Button on:click={submitMarkupAndFetchNext}>Submit</Button>
</div>
<style>

.selected {
  border: 1px solid red;
}
</style>
