
<script>
import Spacer from 'svelte-atoms/Spacer.svelte';
import Block from 'svelte-atoms/Block.svelte';
import ControlBoundingBox from './ControlBoundingBox.svelte';
import ControlRadio from './ControlRadio.svelte';
import ControlCheckbox from './ControlCheckbox.svelte';
import { assessState } from '../../store';
import { getFieldsInOrderFor } from '../../project';
import ControlSubmit from './ControlSubmit.svelte';
import { submitGroup } from '../../control';


export let owner;
export let submitMarkupAndFetchNext;

const ownerGroup = (owner && owner.group) || 'root';
const fields = getFieldsInOrderFor(owner);

</script>

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
