
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


export let templateLike;
export let submitMarkupAndFetchNext;
export let prefix;

let fieldIter = 0;
let keyDown = null;

// console.log('templateLike: ', templateLike)
const [fields, groupsOrder] = getFieldsInOrderFor(templateLike);
if(prefix === 'root') {
  groupsOrder.push(submitGroup);
}

// console.log('fields: ', fields)
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
  console.log('$assessState.controlPrefix prefix: ', $assessState.controlPrefix, prefix)
  if($assessState.controlPrefix !== prefix) {
    return;
  }

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

  keyDown = null;
}

$: if($assessState.controlPrefix === prefix) {
  $assessState.focusedGroup = groupsOrder[fieldIter];
}
</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>

{#each fields as field, i }
  <div id={`device/${prefix}/${i}`}>
    <!-- {console.log('field: ', prefix, i, field)}
    {console.log('fieldsfields: ', prefix, fields)} -->
    <Block type={$assessState.focusedGroup === field.group ? 'selected' : 'block1'}>
      <label>
        <span>{field.group}</span>
        <!-- <ControlClassification {field} /> -->
        {#if field.type === 'radio'}
          <ControlRadio {field} />
        {/if}
        {#if field.type === 'checkbox'}
          <ControlCheckbox {field} />
        {/if}

        {#if field.type === 'bounding_box'}
          <ControlBoundingBox {field} prefix={`${prefix}/${i}`}/>
        {/if}
      </label>
    </Block>
  </div>
{/each}

{#if prefix === 'root'}
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
