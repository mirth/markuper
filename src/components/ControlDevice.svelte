<script>
import Spacer from 'svelte-atoms/Spacer.svelte';
import Block from 'svelte-atoms/Block.svelte';
import ControlRadio from './ControlRadio.svelte';
import ControlCheckbox from './ControlCheckbox.svelte';
import ControlSubmit from './ControlSubmit.svelte';
import { assessState, activeMarkup } from '../store';
import { getFieldsInOrderFor } from '../project';
import { submitGroup } from '../control';

export let submitMarkupAndFetchNext;
export let sample;


$: $activeMarkup = (sample.markup && sample.markup.markup) || {};

const [fields, groupsOrder] = getFieldsInOrderFor(sample.project.template);

groupsOrder.push(submitGroup);

let fieldIter = 0;
let keyDown = null;

function isDeviceFilled() {
  const curField = fields[fieldIter];
  if (curField.type === 'radio' && !Object.hasOwnProperty.call($activeMarkup, curField.group)) {
    return false;
  }

  return true;
}

function focusOnNextField() {
  if (!isDeviceFilled()) {
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

  if (event.key === 'Enter') {
    if (fieldIter < groupsOrder.length - 1) {
      focusOnNextField();
    }
  }
}


$: $assessState.focusedGroup = groupsOrder[fieldIter];

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>

{#each fields as field, i }
  <div id={`device${i}`}>
    <Block type={$assessState.focusedGroup === field.group ? 'selected' : 'block1'}>
      {#if field.type === 'radio'}
        <ControlRadio {field} />
      {/if}
      {#if field.type === 'checkbox'}
        <ControlCheckbox {field} />
      {/if}
    </Block>
  </div>
{/each}
<Spacer size={16} />
<div id='device_submit'>
  <ControlSubmit
    {submitMarkupAndFetchNext}
    isSelected={$assessState.focusedGroup === submitGroup}
    />
</div>


<style>
</style>
