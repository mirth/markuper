
<script>
import Spacer from 'svelte-atoms/Spacer.svelte';
import Block from 'svelte-atoms/Block.svelte';
import ControlBoundingBox from './ControlBoundingBox.svelte';
import ControlRadio from './ControlRadio.svelte';
import ControlCheckbox from './ControlCheckbox.svelte';
import { assessState, activeMarkup } from '../../store';
import { getFieldsInOrderFor } from '../../project';
import ControlSubmit from './ControlSubmit.svelte';
import { submitGroup } from '../../control';


export let owner;
export let submitMarkupAndFetchNext;

const ownerGroup = (owner && owner.group) || 'root';

let fieldIter = 0;
let keyDown = null;

const fields = getFieldsInOrderFor(owner)

let groupsOrder = owner.fields_order;
if(ownerGroup === 'root') {
  groupsOrder = [...groupsOrder, submitGroup];
}

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
  $assessState.focusedGroup = groupsOrder[fieldIter];
}

function handleKeydown(event) {
  keyDown = event.key;
}

function handleKeyup(event) {
  if (event.key !== keyDown) {
    return;
  }

  if($assessState.curFocusedOwner !== ownerGroup) {
    return
  }

  if($assessState.curFocusedOwner !== $assessState.lastFocusedOwner) {
    $assessState.focusedGroup = groupsOrder[fieldIter];
    $assessState.lastFocusedOwner = $assessState.curFocusedOwner
    return
  }

  const focusIsOnThisList = $assessState.curFocusedOwner === ownerGroup;
  if(!focusIsOnThisList) {
    return;
  }

  if($assessState.focusedGroup === submitGroup) {
    return;
  }

  if (event.key === 'Enter') {
    if (fieldIter < groupsOrder.length - 1) {
      focusOnNextField();
      $assessState.curFocusedOwner = ownerGroup;
    } else {
      if($assessState.curFocusedOwner !== 'root') {
        $assessState.lastFocusedOwner = $assessState.curFocusedOwner;
        $assessState.curFocusedOwner = ownerGroup;
      }
    }
  }

  keyDown = null;
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>

{#each fields as field, i }
  <div id={`${ownerGroup}/${i}`}>
    <Block type={$assessState.focusedGroup === field.group ? 'selected' : 'block1'}>
      <label>
        <span>{field.group}</span>
        {#if field.type === 'radio'}
          <ControlRadio {field} />
        {/if}
        {#if field.type === 'checkbox'}
          <ControlCheckbox {field} />
        {/if}

        {#if field.type === 'bounding_box'}
          <ControlBoundingBox {field} />
        {/if}
      </label>
    </Block>
  </div>
{/each}

{#if submitMarkupAndFetchNext}
  <Spacer size={16} />
  <div id='device_submit'>
    <ControlSubmit
      {submitMarkupAndFetchNext}
      isSelected={$assessState.focusedGroup === submitGroup}
      />
  </div>
{/if}

<style>

span {
  display: block;
  margin-bottom: 5px;
}

</style>
